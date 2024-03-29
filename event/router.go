package event

import (
	"container/list"
	"sync"
)

type router struct {
	sync.RWMutex
	routes       *list.List
	filter       Filter
	delRouteChan chan *route
	quit         chan struct{}
	complete     chan struct{}
}

func newRouter(filters ...Filter) *router {
	var filter Filter
	if len(filters) > 0 && filters[0] != nil {
		filter = filters[0]
	}

	r := &router{
		routes:       list.New(),
		filter:       filter,
		delRouteChan: make(chan *route),
		quit:         make(chan struct{}),
		complete:     make(chan struct{}),
	}

	go r.process()
	return r
}

func (r *router) close() {
	select {
	case <-r.quit:
	default:
		close(r.quit)
		r.RLock()
		defer r.RUnlock()
		for e := r.routes.Front(); e != nil; e = e.Next() {
			e.Value.(*route).close()
		}
	}
}

func (r *router) done() {
	r.close()
	<-r.complete
}

func (r *router) flowQuit() {
	select {
	case <-r.quit:
	default:
		close(r.quit)
	}
}

func (r *router) flowComplete() {
	select {
	case <-r.complete:
	default:
		close(r.complete)
		close(r.delRouteChan)
	}
}

func (r *router) addRoute(topic string, callback Handler) *Token {
	select {
	case <-r.quit:
		return nil
	default:
		r.Lock()
		defer r.Unlock()
		for e := r.routes.Front(); e != nil; e = e.Next() {
			if e.Value.(*route).topic == topic {
				return e.Value.(*route).addToken(callback)
			}
		}

		rou, token := newRoute(topic, callback, r.delRouteChan)
		r.routes.PushBack(rou)
		return token
	}
}

func (r *router) delRoute(rou *route) {
	r.Lock()
	defer r.Unlock()
	defer func() {
		if r.routes.Len() == 0 {
			r.flowQuit()
			r.flowComplete()
		}
	}()
	for e := r.routes.Front(); e != nil; e = e.Next() {
		if e.Value.(*route).topic == rou.topic {
			r.routes.Remove(e)
			return
		}
	}
}

func (r *router) process() {
	for rou := range r.delRouteChan {
		r.delRoute(rou)
	}
}

func (r *router) publish(event *Event) {
	select {
	case <-r.quit:
	default:
		r.RLock()
		defer r.RUnlock()
		if r.filter != nil {
			for e := r.routes.Front(); e != nil; e = e.Next() {
				if r.filter(event.Type, e.Value.(*route).topic) {
					e.Value.(*route).publish(event)
				}
			}
		} else {
			for e := r.routes.Front(); e != nil; e = e.Next() {
				if e.Value.(*route).topic == event.Type {
					e.Value.(*route).publish(event)
				}
			}
		}
	}
}

func (r *router) unsubscribe(topics ...string) {
	select {
	case <-r.quit:
	default:
		r.RLock()
		defer r.RUnlock()
		for _, topic := range topics {
			for e := r.routes.Front(); e != nil; e = e.Next() {
				if e.Value.(*route).topic == topic {
					e.Value.(*route).close()
					break
				}
			}
		}
	}
}
