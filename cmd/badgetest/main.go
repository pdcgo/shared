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
			source := yenstream.NewSliceSource(ctx, []uint{
				1,
				2,
				3,
				4,
				5,
			})

			return source.
				Via("sum", yenstream.NewCombiner(ctx, &sumCombiner{}, nil)).
				Via("debug", yenstream.NewMap(ctx, func(data any) (any, error) {
					raw, err := json.Marshal(data)
					log.Println(string(raw))
					return data, err
				}))
		})
}

type UintValue struct {
	yenstream.Metadata
	Value uint `json:"value"`
}

type sumCombiner struct{}

func (*sumCombiner) AddInput(data uint, acc *UintValue) *UintValue {
	acc.Value += data
	return acc
}

func (*sumCombiner) CreateAccumulator() *UintValue {
	return &UintValue{
		Value: 0,
	}
}
