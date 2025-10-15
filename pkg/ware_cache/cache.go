package ware_cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache missing")

type CacheItem struct {
	Key        string
	Expiration time.Duration
	Data       any
}

func (c *CacheItem) Serialize() ([]byte, error) {
	return json.Marshal(c.Data)
}

func (c *CacheItem) Deserialize(data []byte) error {
	return json.Unmarshal(data, c.Data)
}

type Cache interface {
	Add(ctx context.Context, item *CacheItem) error
	Replace(ctx context.Context, item *CacheItem) error
	Get(ctx context.Context, key string, data any) error
	GetRaw(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
