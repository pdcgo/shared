package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pdcgo/shared/yenstream"
	"github.com/pdcgo/shared/yenstream/store"
)

func main() {
	pctx := context.Background()
	ctx, db, err := store.RegisterBadgerStore(pctx, "badgertest")
	if err != nil {
		panic(err)
	}

	defer func() {
		db.DropAll()
		db.Close()
	}()

	yenstream.
		NewRunnerContext(ctx).
		CreatePipeline(func(ctx *yenstream.RunnerContext) yenstream.Pipeline {
			source := yenstream.NewSliceSource(ctx, [][]uint{
				{1, 2},
				{1, 2},
				{1, 2},
				{1, 2},
				{2, 2},
				{3, 6},
				{3, 6},
				{3, 6},
				{3, 6},
			})

			return source.
				Via("debug", yenstream.NewMap(ctx, debugFunc))
		})
}

func debugFunc(data any) (any, error) {
	raw, err := json.Marshal(data)
	log.Println(string(raw))
	return data, err
}
