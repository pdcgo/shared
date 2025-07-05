package ware_cache

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type localCacheImpl struct {
	data       map[string][]byte
	created    map[string]time.Time
	expiration map[string]time.Duration
}

// Flush implements Cache.
func (l *localCacheImpl) Flush(ctx context.Context) error {
	l.created = map[string]time.Time{}
	l.data = map[string][]byte{}
	l.expiration = map[string]time.Duration{}
	return nil
}

// Delete implements Cache.
func (l *localCacheImpl) Delete(ctx context.Context, key string) error {
	l.data[key] = nil
	return nil
}

// Add implements Cache.
func (l *localCacheImpl) Add(ctx context.Context, item *CacheItem) error {

	raw, err := item.Serialize()
	if err != nil {
		return err
	}
	l.created[item.Key] = time.Now()
	l.expiration[item.Key] = item.Expiration
	l.data[item.Key] = raw
	log.Println("add cache", l.data[item.Key], item.Key)

	return nil
}

// Get implements Cache.
func (l *localCacheImpl) Get(ctx context.Context, key string, data any) error {
	// log.Println("get cache", l.data, key)

	item := l.data[key]
	if item == nil {
		return ErrCacheMiss
	}

	log.Println("expiration", l.created[key], l.expiration[key])

	if l.expiration[key] != 0 {
		dur := time.Since(l.created[key])
		if dur > l.expiration[key] {
			return ErrCacheMiss
		}
	}

	if len(l.data[key]) == 0 {
		return ErrCacheMiss
	}

	err := json.Unmarshal(l.data[key], data)
	if err != nil {
		return err
	}
	return nil
}

// Replace implements Cache.
func (l *localCacheImpl) Replace(ctx context.Context, item *CacheItem) error {
	return l.Add(ctx, item)
}

func NewLocalCache() Cache {
	return &localCacheImpl{
		data:       map[string][]byte{},
		created:    map[string]time.Time{},
		expiration: map[string]time.Duration{},
	}
}
