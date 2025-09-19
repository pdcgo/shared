package ware_cache_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pdcgo/shared/pkg/ware_cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T) *ware_cache.BadgerCache {
	dir := filepath.Join(os.TempDir(), "badger_test")
	_ = os.RemoveAll(dir) // cleanup before test

	c, err := ware_cache.NewBadgerCache(dir)
	require.NoError(t, err)

	t.Cleanup(func() {
		c.Close()
		os.RemoveAll(dir)
	})

	return c
}

func TestAddAndGet(t *testing.T) {
	c := newTestCache(t)
	ctx := context.Background()

	item := &ware_cache.CacheItem{
		Key:        "foo",
		Expiration: time.Minute,
		Data:       map[string]string{"msg": "hello"},
	}
	err := c.Add(ctx, item)
	require.NoError(t, err)

	var result map[string]string
	err = c.Get(ctx, "foo", &result)
	require.NoError(t, err)
	require.Equal(t, "hello", result["msg"])

	// Adding duplicate should fail
	err = c.Add(ctx, item)
	assert.Nil(t, err)
}

func TestReplace(t *testing.T) {
	c := newTestCache(t)
	ctx := context.Background()

	// Replace on missing key should fail
	item := &ware_cache.CacheItem{
		Key:  "bar",
		Data: "first",
	}
	err := c.Replace(ctx, item)
	assert.Nil(t, err)

	// Add first
	require.NoError(t, c.Add(ctx, item))

	// Replace should work
	item.Data = "second"
	require.NoError(t, c.Replace(ctx, item))

	var result string
	require.NoError(t, c.Get(ctx, "bar", &result))
	require.Equal(t, "second", result)
}

func TestDelete(t *testing.T) {
	c := newTestCache(t)
	ctx := context.Background()

	item := &ware_cache.CacheItem{
		Key:  "baz",
		Data: 123,
	}
	require.NoError(t, c.Add(ctx, item))

	// Delete it
	require.NoError(t, c.Delete(ctx, "baz"))

	// Should not be found
	var out int
	err := c.Get(ctx, "baz", &out)
	require.ErrorIs(t, err, ware_cache.ErrCacheMiss)
}

func TestFlush(t *testing.T) {
	c := newTestCache(t)
	ctx := context.Background()

	require.NoError(t, c.Add(ctx, &ware_cache.CacheItem{
		Key:  "flushme",
		Data: true,
	}))

	require.NoError(t, c.Flush(ctx))

	var out bool
	err := c.Get(ctx, "flushme", &out)
	require.ErrorIs(t, err, ware_cache.ErrCacheMiss)
}
