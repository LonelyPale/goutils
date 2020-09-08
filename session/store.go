package session

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/LonelyPale/goutils/cache"
	"github.com/LonelyPale/goutils/errors"
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

func (m *MemoryStore) Session(request interface{}) (Session, error) {
	switch req := request.(type) {
	case http.Request:
		authorization := req.Header.Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer"))
		return m.Get(token)
	case context.Context:
		md, ok := metadata.FromIncomingContext(req)
		if !ok {
			return nil, errors.New("无Token认证信息")
		}

		var token string
		val, ok := md["token"]
		if !ok {
			return nil, errors.New("无Token认证信息")
		}

		token = val[0]
		return m.Get(token)
	default:
		return nil, errors.New("bad request auth")
	}
}
