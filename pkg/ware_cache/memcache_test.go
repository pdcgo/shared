package ware_cache_test

import (
	"testing"
)

type MockCacheData struct {
	Data string
}

func TestMemcache(t *testing.T) {

	// cache := ware_cache.NewMemcache()
	// err := cache.Add(context.Background(), &ware_cache.CacheItem{
	// 	Key:        "test",
	// 	Expiration: time.Hour,
	// 	Data: &MockCacheData{
	// 		Data: "asdasdasdasdasdasdasdasdasd",
	// 	},
	// })

	// assert.Nil(t, err)
}
