package session

import (
	"sync"

	"github.com/LonelyPale/goutils/cache"
	"github.com/LonelyPale/goutils/uuid"
)

var store *cache.Cache

func init() {
	var err error
	store, err = cache.New()
	if err != nil {
		panic(err)
	}
}

type Session interface {
	ID() string
	Get(key interface{}) interface{}
	Set(key, value interface{})
	Delete(ey interface{})
	Save() error
}

type SimpleSession struct {
	mutex sync.RWMutex
	Id    string
	Data  map[interface{}]interface{}
}

func NewSession(ids ...string) (Session, error) {
	s := &SimpleSession{Data: make(map[interface{}]interface{})}

	if len(ids) > 0 && len(ids[0]) > 0 {
		s.Id = ids[0]
		if err := store.Get(s.Id, s); err != nil {
			return nil, err
		}
	} else {
		s.Id = uuid.New().String()
	}

	return s, nil
}

func (s *SimpleSession) ID() string {
	return s.Id
}

func (s *SimpleSession) Get(key interface{}) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	val, ok := s.Data[key]
	if !ok {
		//return nil, errors.New("key does not exist")
		return nil
	}
	return val
}

func (s *SimpleSession) Set(key, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Data[key] = value
}

func (s *SimpleSession) Delete(key interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.Data[key]; ok {
		delete(s.Data, key)
	}
}

func (s *SimpleSession) Save() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return store.Set(s.Id, s)
}
