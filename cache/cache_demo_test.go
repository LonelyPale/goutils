package cache

import (
	"testing"
	"time"

	"github.com/allegro/bigcache"
)

var config = bigcache.Config{
	// number of shards (must be a power of 2)
	Shards: 1024,

	// time after which entry can be evicted
	LifeWindow: 5 * time.Second,

	// Interval between removing expired entries (clean up).
	// If set to <= 0 then no action is performed.
	// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
	CleanWindow: 5 * time.Second,

	// rps * lifeWindow, used only in initial memory allocation
	MaxEntriesInWindow: 1000 * 10 * 60,

	// max entry size in bytes, used only in initial memory allocation
	MaxEntrySize: 10240,

	// prints information about additional memory allocation
	Verbose: true,

	// cache will not allocate more memory than this limit, value in MB
	// if value is reached then the oldest entries can be overridden for the new ones
	// 0 value means no size limit
	HardMaxCacheSize: 8192,

	// callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	OnRemove: nil,

	// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
	// for the new entry, or because delete was called. A constant representing the reason will be passed through.
	// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
	// Ignored if OnRemove is specified.
	OnRemoveWithReason: nil,
}

func Test(t *testing.T) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	if err := cache.Set("my-unique-key", []byte("value")); err != nil {
		t.Fatal(err)
	}

	entry, err := cache.Get("my-unique-key")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(entry))
}

func Test1(t *testing.T) {
	cache, initErr := bigcache.NewBigCache(config)
	if initErr != nil {
		t.Fatal(initErr)
	}

	if err := cache.Set("my-unique-key", []byte("value")); err != nil {
		t.Fatal(err)
	}

	if entry, err := cache.Get("my-unique-key"); err == nil {
		t.Log(string(entry))
	} else {
		t.Fatal(err)
	}
}

func Test2(t *testing.T) {
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		t.Fatal(err)
	}

	if err := cache.Set("my-unique-key", []byte("value")); err != nil {
		t.Fatal(err)
	} else {
		t.Log(time.Now())
	}

	time.Sleep(10 * time.Second)
	if err := cache.Set("my-unique-key", []byte("value")); err != nil {
		t.Fatal(err)
	} else {
		t.Log(time.Now())
	}

	time.Sleep(10 * time.Second)
	entry, err := cache.Get("my-unique-key")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(time.Now(), string(entry))
	}

	time.Sleep(2 * time.Second)
	entry, err = cache.Get("my-unique-key")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(time.Now(), string(entry))
	}
}
