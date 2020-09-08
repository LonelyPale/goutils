package session

import (
	"sync"

	"github.com/LonelyPale/goutils/random/uuid"
)

type Session interface {
	ID() string
	Get(key interface{}) interface{}
	Set(key, value interface{})
	Delete(key interface{})
	Save() error
}

// Simple Session
type session struct {
	store  Store
	id     string
	data   map[interface{}]interface{}
	dataMu sync.RWMutex
}

func NewSession(store Store, ids ...string) Session {
	var id string
	if len(ids) > 0 && len(ids[0]) > 0 {
		id = ids[0]
	} else {
		id = uuid.New().String()
	}

	return &session{
		store: store,
		id:    id,
		data:  make(map[interface{}]interface{}),
	}
}

func (s *session) ID() string {
	return s.id
}

func (s *session) Get(key interface{}) interface{} {
	s.dataMu.RLock()
	defer s.dataMu.RUnlock()

	val, ok := s.data[key]
	if !ok {
		//return nil, errors.New("key does not exist")
		return nil
	}

	return val
}

func (s *session) Set(key, value interface{}) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	s.data[key] = value
}

func (s *session) Delete(key interface{}) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
	}
}

func (s *session) Save() error {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	return s.store.Save(s)
}
