package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger/v4"
)

var BADGER_DB = "badger_db"

func RegisterBadgerStore(pctx context.Context, streamID string) (context.Context, *badger.DB, error) {
	dirname := fmt.Sprintf("/tmp/%s", streamID)
	opts := badger.DefaultOptions(dirname).WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(opts)
	if err != nil {
		return pctx, nil, err
	}

	ctx := context.WithValue(pctx, BADGER_DB, db)
	return ctx, db, nil
}

func NewBadgerStore[R any](db *badger.DB, id string, combiner func() R) StateStore {
	return &badgerImpl[R]{
		id:       id,
		db:       db,
		combiner: combiner,
	}
}

type badgerImpl[R any] struct {
	id       string
	db       *badger.DB
	combiner func() R
}

func (b *badgerImpl[R]) idWithNamespace(key string) []byte {
	res := b.id + "/" + key
	return []byte(res)
}

func (b *badgerImpl[R]) serializeKey(key any) ([]byte, error) {
	switch tkey := key.(type) {
	case string:
		return b.idWithNamespace(tkey), nil
	case int:
		key := strconv.FormatInt(int64(tkey), 10)
		return b.idWithNamespace(key), nil
	case uint:
		key := strconv.FormatUint(uint64(tkey), 10)
		return b.idWithNamespace(key), nil
	case int64:
		key := strconv.FormatInt(tkey, 10)
		return b.idWithNamespace(key), nil
	default:
		tipe := reflect.TypeOf(key)
		if tipe.Kind() == reflect.Ptr {
			tipe = tipe.Elem()
		}
		name := tipe.Name()
		err := fmt.Errorf("key with type %s not supported", name)
		return []byte{}, err
	}
}

// Get implements yenstream.StateStore.
func (b *badgerImpl[R]) Get(key any) any {
	bkey, err := b.serializeKey(key)
	if err != nil {
		panic(err)
	}
	var res any = b.combiner()
	err = b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(bkey)
		if err != nil {
			if !errors.Is(err, badger.ErrKeyNotFound) {
				return err
			}

			return nil
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(val, res)
		return err
	})

	if err != nil {
		slog.Error(err.Error(), slog.String("store_id", b.id))
	}
	return res
}

// GetAll implements yenstream.StateStore.
func (b *badgerImpl[R]) GetAll(emitter func(key any, data any)) {
	prefix := []byte(b.id + "/")

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.Prefix = prefix // ✅ only keys starting with "user_"

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var res any = b.combiner()
			err = json.Unmarshal(v, res)
			if err != nil {
				return err
			}
			key, _ := strings.CutPrefix(string(k), string(prefix))
			// log.Println(string(k), key, res)
			emitter(key, res)
		}
		return nil
	})

	if err != nil {
		slog.Error(err.Error(), slog.String("store_id", b.id))
	}
}

// Set implements yenstream.StateStore.
func (b *badgerImpl[R]) Set(key any, value any) {
	bkey, err := b.serializeKey(key)
	if err != nil {
		panic(err)
	}
	err = b.db.Update(func(txn *badger.Txn) error {
		raw, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return txn.Set(bkey, raw)
	})

	if err != nil {
		slog.Error(err.Error(), slog.String("store_id", b.id))
	}
}
