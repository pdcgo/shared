package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/pdcgo/shared/yenstream"
)

func NewBadgerStore[R any](db *badger.DB, id string, combiner func() R) yenstream.StateStore {
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

func (b *badgerImpl[R]) serializeKey(key any) (string, error) {
	switch tkey := key.(type) {
	case string:
		return tkey, nil
	case int:
		return strconv.FormatInt(int64(tkey), 10), nil
	case uint:
		return strconv.FormatUint(uint64(tkey), 10), nil
	default:
		tipe := reflect.TypeOf(key)
		if tipe.Kind() == reflect.Ptr {
			tipe = tipe.Elem()
		}
		name := tipe.Name()
		err := fmt.Errorf("key with type %s not supported", name)
		return "", err
	}
}

// Get implements yenstream.StateStore.
func (b *badgerImpl[R]) Get(key any) any {
	b.db.View(func(txn *badger.Txn) error {
		panic("unimplemented")
	})
	return nil
}

// GetAll implements yenstream.StateStore.
func (b *badgerImpl[R]) GetAll(emitter func(key any, data any)) {
	panic("unimplemented")
}

// Set implements yenstream.StateStore.
func (b *badgerImpl[R]) Set(key any, value any) {
	panic("unimplemented")
}

func main() {
	opts := badger.DefaultOptions("/tmp/streamid")
	// WithLoggingLevel(badger.INFO)
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("name"), []byte("pdcoke"))
		return txn.SetEntry(e)
	})

	defer db.Close()
}
