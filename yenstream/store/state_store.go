package store

import (
	"context"
	"log/slog"

	"github.com/dgraph-io/badger/v4"
)

type StateStore interface {
	Get(key any) any
	Set(key any, value any)
	GetAll(emitter func(key any, data any))
}

func CreateStoreFromCtx[R any](ctx context.Context, key string, combiner func() R) StateStore {
	var storeimpl StateStore
	db, ok := ctx.Value(BADGER_DB).(*badger.DB)
	if ok {
		storeimpl = NewBadgerStore(db, key, combiner)
	} else {
		slog.Warn("yenstream using memory keymap primitive")
		storeimpl = &keyMapStoreImpl{
			state: map[any]any{},
		}
	}

	return storeimpl
}
