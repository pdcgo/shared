package ware_cache

import (
	"context"
	"encoding/json"
	"errors"

	"google.golang.org/appengine/memcache"
)

type memcacheImpl struct{}

// Flush implements Cache.
func (m *memcacheImpl) Flush(ctx context.Context) error {
	return memcache.Flush(ctx)
}

// Delete implements Cache.
func (m *memcacheImpl) Delete(ctx context.Context, key string) error {
	err := memcache.Delete(ctx, key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil
		}
		return err
	}
	return nil
}

// Add implements Cache.
func (m *memcacheImpl) Add(ctx context.Context, item *CacheItem) error {
	value, err := item.Serialize()
	if err != nil {
		return err
	}
	err = memcache.Add(ctx, &memcache.Item{
		Key:        item.Key,
		Value:      value,
		Expiration: item.Expiration,
	})

	if err != nil {
		if errors.Is(err, memcache.ErrNotStored) {
			return nil
		}
		return err
	}

	return nil
}

// Get implements Cache.
func (m *memcacheImpl) Get(ctx context.Context, key string, data any) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return ErrCacheMiss
		}

		return err
	}

	err = json.Unmarshal(item.Value, data)
	if err != nil {
		return err
	}

	return nil
}

func NewMemcache() Cache {
	return &memcacheImpl{}
}
