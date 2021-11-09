package event

import (
	"container/list"
	"sync"
)

type route struct {
	sync.RWMutex
	topic        string
	tokens       *list.List
	delTokenChan chan *Token
	delRouteChan chan<- *route
	quit         chan struct{}
	complete     chan struct{}
}

func newRoute(topic string, callback Handler, delRouteChan chan<- *route) (*route, *Token) {
	delTokenChan := make(chan *Token)
	token := newToken(callback, delTokenChan)
	tokens := list.New()
	tokens.PushBack(token)
	r := &route{
		topic:        topic,
		tokens:       tokens,
		delTokenChan: delTokenChan,
		delRouteChan: delRouteChan,
		quit:         make(chan struct{}),
		complete:     make(chan struct{}),
	}

	go r.process()
	return r, token
}

func (r *route) isClose() bool {
	select {
	case <-r.quit:
		return true
	default:
		return false
	}
}

func (r *route) close() {
	select {
	case <-r.quit:
	default:
		close(r.quit)
		r.RLock()
		defer r.RUnlock()
		for e := r.tokens.Front(); e != nil; e = e.Next() {
			e.Value.(*Token).Close()
		}
	}
}

func (r *route) done() {
	r.close()
	<-r.complete
}

func (r *route) flowQuit() {
	select {
	case <-r.quit:
	default:
		close(r.quit)
	}
}

func (r *route) flowComplete() {
	select {
	case <-r.complete:
	default:
		close(r.complete)
		close(r.delTokenChan)
		r.delRouteChan <- r
	}
}

func (r *route) addToken(callback Handler) *Token {
	select {
	case <-r.quit:
		return nil
	default:
		r.Lock()
		defer r.Unlock()
		token := newToken(callback, r.delTokenChan)
		r.tokens.PushBack(token)
		return token
	}
}

func (r *route) delToken(token *Token) {
	r.Lock()
	defer r.Unlock()
	defer func() {
		if r.tokens.Len() == 0 {
			r.flowQuit()
			r.flowComplete()
		}
	}()
	for e := r.tokens.Front(); e != nil; e = e.Next() {
		if e.Value.(*Token).id == token.id {
			r.tokens.Remove(e)
			return
		}
	}
}

func (r *route) process() {
	for token := range r.delTokenChan {
		r.delToken(token)
	}
}

func (r *route) publish(event *Event) {
	select {
	case <-r.quit:
	default:
		r.RLock()
		defer r.RUnlock()
		for e := r.tokens.Front(); e != nil; e = e.Next() {
			e.Value.(*Token).publish(event)
		}
	}
}
