package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/LonelyPale/goutils/errors"
)

// Nil reply Redis returns when key does not exist.
const Nil = redis.Nil

type DB struct {
	client *redis.Client
}

func NewRedisDB(cfg *Config) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "ping redis")
	}

	return &DB{
		client: client,
	}, nil
}

type CacheAble interface {
	CacheKey() string
}

type Cache interface {
	Set(string, interface{}, time.Duration) error
	Get(string, interface{}) error
	Del(string) *redis.IntCmd
	Scan(uint64, string, int64) *redis.ScanCmd
	Reset() *redis.StatusCmd
	Close() error
}

func (r *DB) Set(key string, obj interface{}, expiration time.Duration) error {
	bytes, err := msgpack.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "msgpack marshal")
	}

	return r.client.Set(key, bytes, expiration).Err()
}

func (r *DB) Get(key string, obj interface{}) error {
	result, err := r.client.Get(key).Bytes()
	if err != nil {
		return Nil
	}

	return msgpack.Unmarshal(result, obj)
}

func (r *DB) Del(key string) *redis.IntCmd {
	return r.client.Del(key)
}

func (r *DB) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.Scan(cursor, match, count)
}

func (r *DB) Reset() *redis.StatusCmd {
	return r.client.FlushDB()
}

func (r *DB) Close() error {
	return r.client.Close()
}
