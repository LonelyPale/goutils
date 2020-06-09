package cache

import (
	"time"

	"github.com/allegro/bigcache"
)

const (
	DefaultLifeWindow  = 1 * time.Hour
	DefaultCleanWindow = 10 * time.Minute
)

func DefaultConfig(ts ...time.Duration) bigcache.Config {
	life := DefaultLifeWindow
	clean := DefaultCleanWindow

	switch len(ts) {
	case 1:
		life = ts[0]
	case 2:
		life = ts[0]
		clean = ts[1]
	}

	return bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: life,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: clean,

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

		Logger: bigcache.DefaultLogger(),
	}
}

func New(configs ...bigcache.Config) (*bigcache.BigCache, error) {
	var config bigcache.Config
	if len(configs) > 0 {
		config = configs[0]
	} else {
		config = DefaultConfig()
	}

	return bigcache.NewBigCache(config)
}
