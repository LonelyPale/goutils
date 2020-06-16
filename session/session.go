package session

import (
	"sync"

	"github.com/LonelyPale/goutils/uuid"
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
	mutex sync.RWMutex
	store Store
	Id    string
	Data  map[interface{}]interface{}
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
		Id:    id,
		Data:  make(map[interface{}]interface{}),
	}
}

func (s *session) ID() string {
	return s.Id
}

func (s *session) Get(key interface{}) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	val, ok := s.Data[key]
	if !ok {
		//return nil, errors.New("key does not exist")
		return nil
	}
	return val
}

func (s *session) Set(key, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Data[key] = value
}

func (s *session) Delete(key interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.Data[key]; ok {
		delete(s.Data, key)
	}
}

func (s *session) Save() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.store.Save(s)
}
