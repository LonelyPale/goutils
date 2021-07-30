package validator

import (
	"sync"
	"sync/atomic"
)

type validateCache struct {
	cache atomic.Value // map[string]*validate
	lock  sync.Mutex
}

func newValidateCache() *validateCache {
	s := new(validateCache)
	s.cache.Store(make(map[string]*validate))
	return s
}

func (s *validateCache) Get(key string) (v *validate, found bool) {
	v, found = s.cache.Load().(map[string]*validate)[key]
	return
}

func (s *validateCache) Set(key string, value *validate) {
	cache := s.cache.Load().(map[string]*validate)
	newCache := make(map[string]*validate, len(cache)+1)
	for k, v := range cache {
		newCache[k] = v
	}
	newCache[key] = value
	s.cache.Store(newCache)
}
