package session

import (
	"github.com/LonelyPale/goutils/cache"
)

type Store interface {
	New() Session
	Get(id string) (Session, error)
	Delete(id string) error
	Save(s Session) error
}

type MemoryStore struct {
	store *cache.Cache
}

func NewMemoryStore() (Store, error) {
	store, err := cache.New()
	if err != nil {
		return nil, err
	}

	return &MemoryStore{store}, nil
}

func (m *MemoryStore) New() Session {
	return NewSession(m)
}

func (m *MemoryStore) Get(id string) (Session, error) {
	s := NewSession(m, id)
	if err := m.store.Get(id, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *MemoryStore) Delete(id string) error {
	return m.store.Delete(id)
}

func (m *MemoryStore) Save(s Session) error {
	return m.store.Set(s.ID(), s)
}
