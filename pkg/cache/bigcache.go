package cache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/chenyahui/gin-cache/persist"
	"go.uber.org/zap"
	"time"
)

// BigCacheStore local memory cache store
type BigCacheStore struct {
	Cache *bigcache.BigCache
}

// NewBigCacheStore allocate a local memory store with default expiration
func NewBigCacheStore(defaultExpiration time.Duration, logger *zap.Logger) *BigCacheStore {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         defaultExpiration,
		CleanWindow:        3 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1500,
		StatsEnabled:       false,
		Verbose:            true,
		HardMaxCacheSize:   256, // memory limit, value in MB
		Logger:             NewZapLogger(logger),
	}
	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return &BigCacheStore{
		Cache: cache,
	}
}

// Set put key value pair to memory store, and expire after expireDuration
func (c *BigCacheStore) Set(key string, value interface{}, expireDuration time.Duration) error {
	_ = expireDuration
	payload, err := persist.Serialize(value)
	if err != nil {
		return err
	}
	return c.Cache.Set(key, payload)
}

// Delete remove key in memory store, do nothing if key doesn't exist
func (c *BigCacheStore) Delete(key string) error {
	return c.Cache.Delete(key)
}

// Get get key in memory store, if key doesn't exist, return ErrCacheMiss
func (c *BigCacheStore) Get(key string, value interface{}) error {
	payload, err := c.Cache.Get(key)

	if err == bigcache.ErrEntryNotFound {
		return persist.ErrCacheMiss
	}

	if err != nil {
		return err
	}
	return persist.Deserialize(payload, value)
}
