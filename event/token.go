package event

import (
	"sync"

	"github.com/LonelyPale/goutils/errors"
	"github.com/LonelyPale/goutils/types"
)

type Token struct {
	id           types.ObjectID
	delTokenChan chan<- *Token
	eventChan    chan *Event
	callback     Handler
	quit         chan struct{}
	complete     chan struct{}
	err          error
	errMu        sync.RWMutex
}

func newToken(callback Handler, delTokenChan chan<- *Token) *Token {
	token := &Token{
		id:           types.NewObjectID(),
		delTokenChan: delTokenChan,
		eventChan:    make(chan *Event),
		callback:     callback,
		quit:         make(chan struct{}),
		complete:     make(chan struct{}),
	}

	go token.process()
	return token
}

func (t *Token) IsClose() bool {
	select {
	case <-t.quit:
		return true
	default:
		return false
	}
}

func (t *Token) Close() {
	select {
	case <-t.quit:
	default:
		close(t.quit)
		close(t.eventChan)
	}
}

func (t *Token) Done() {
	t.Close()
	<-t.complete
}

func (t *Token) flowComplete() {
	select {
	case <-t.complete:
	default:
		close(t.complete)
		t.delTokenChan <- t
	}
}

func (t *Token) Error() error {
	t.errMu.RLock()
	defer t.errMu.RUnlock()
	return t.err
}

func (t *Token) setError(e error) {
	t.errMu.Lock()
	defer t.errMu.Unlock()
	t.err = e
	t.Close()
}

//处理事件
func (t *Token) process() {
	defer t.flowComplete()
	defer func() {
		if r := recover(); r != nil {
			t.setError(errors.Error(r))
			//debug.PrintStack()
		}
	}()

	for event := range t.eventChan {
		t.callback(*event)
	}
}

//发布事件
func (t *Token) publish(event *Event) {
	select {
	case <-t.quit:
	default:
		t.eventChan <- event
	}
}
