package ware_cache

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/cache_iface/v1"
	"github.com/pdcgo/schema/services/cache_iface/v1/cache_ifaceconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type customCacheImpl struct {
	client cache_ifaceconnect.CacheServiceClient
}

// Add implements Cache.
func (c *customCacheImpl) Add(ctx context.Context, item *CacheItem) error {
	val, err := json.Marshal(item.Data)
	if err != nil {
		return err
	}
	_, err = c.client.Add(ctx, &connect.Request[cache_iface.AddRequest]{
		Msg: &cache_iface.AddRequest{
			Key:      item.Key,
			ExpireAt: timestamppb.New(time.Now().Add(item.Expiration)),
			Value:    val,
		},
	})
	return err
}

// Delete implements Cache.
func (c *customCacheImpl) Delete(ctx context.Context, key string) error {
	_, err := c.client.Delete(ctx, &connect.Request[cache_iface.DeleteRequest]{
		Msg: &cache_iface.DeleteRequest{
			Key: key,
		},
	})
	return err
}

// Flush implements Cache.
func (c *customCacheImpl) Flush(ctx context.Context) error {
	_, err := c.client.Flush(ctx, &connect.Request[cache_iface.FlushRequest]{
		Msg: &cache_iface.FlushRequest{},
	})
	return err
}

// Get implements Cache.
func (c *customCacheImpl) Get(ctx context.Context, key string, data any) error {
	res, err := c.client.Get(ctx, &connect.Request[cache_iface.GetRequest]{
		Msg: &cache_iface.GetRequest{
			Key: key,
		},
	})

	if err != nil {
		return err
	}
	if res.Msg.Missed {
		return ErrCacheMiss
	}

	err = json.Unmarshal(res.Msg.Value, data)
	return err
}

// GetRaw implements Cache.
func (c *customCacheImpl) GetRaw(ctx context.Context, key string) ([]byte, error) {
	res, err := c.client.Get(ctx, &connect.Request[cache_iface.GetRequest]{
		Msg: &cache_iface.GetRequest{
			Key: key,
		},
	})

	if err != nil {
		return nil, err
	}
	if res.Msg.Missed {
		return []byte{}, ErrCacheMiss
	}

	return res.Msg.Value, nil
}

// Replace implements Cache.
func (c *customCacheImpl) Replace(ctx context.Context, item *CacheItem) error {
	val, err := item.Serialize()
	if err != nil {
		return err
	}
	_, err = c.client.Replace(ctx, &connect.Request[cache_iface.ReplaceRequest]{
		Msg: &cache_iface.ReplaceRequest{
			Key:      item.Key,
			ExpireAt: timestamppb.New(time.Now().Add(item.Expiration)),
			Value:    val,
		},
	})

	return err
}

func NewCustomCache(endpoint string) Cache {
	return &customCacheImpl{
		client: cache_ifaceconnect.NewCacheServiceClient(http.DefaultClient, endpoint),
	}
}
