package ware_cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/badger/v4"
)

type BadgerCache struct {
	db *badger.DB
}

// GetRaw implements Cache.
func (b *BadgerCache) GetRaw(ctx context.Context, key string) ([]byte, error) {
	var err error
	var data []byte
	err = b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrCacheMiss
			}
			return err
		}
		return item.Value(func(val []byte) error {
			data = val
			return nil
		})
	})
	return data, err
}

// NewBadgerCache opens a BadgerDB instance at the given path.
func NewBadgerCache(path string) (*BadgerCache, error) {
	opts := badger.DefaultOptions(path).
		WithLogger(nil) // silence internal logs, optional

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerCache{db: db}, nil
}

func (b *BadgerCache) Close() error {
	return b.db.Close()
}

// Add inserts a new item. If the key already exists, it returns an error.
func (b *BadgerCache) Add(ctx context.Context, item *CacheItem) error {
	data, err := item.Serialize()
	if err != nil {
		return err
	}

	return b.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(item.Key))
		if err == nil {
			return nil
		}
		e := badger.NewEntry([]byte(item.Key), data)
		if item.Expiration > 0 {
			e = e.WithTTL(item.Expiration)
		}
		return txn.SetEntry(e)
	})
}

// Flush deletes all keys (dangerous, use carefully).
func (b *BadgerCache) Flush(ctx context.Context) error {
	return b.db.DropAll()
}

func (b *BadgerCache) Get(ctx context.Context, key string, data any) error {
	return b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrCacheMiss
			}
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, data)
		})
	})
}

func (b *BadgerCache) Replace(ctx context.Context, item *CacheItem) error {
	data, err := item.Serialize()
	if err != nil {
		return err
	}

	return b.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(item.Key))
		if err != nil {
			if !errors.Is(err, badger.ErrKeyNotFound) {
				return err
			}

		}
		e := badger.NewEntry([]byte(item.Key), data)
		if item.Expiration > 0 {
			e = e.WithTTL(item.Expiration)
		}
		return txn.SetEntry(e)
	})
}

func (b *BadgerCache) Delete(ctx context.Context, key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		return nil
	})
}
