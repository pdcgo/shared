package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pdcgo/shared/yenstream"
)

type sumCombiner struct{}

// AddInput implements yenstream.Accumulator.
func (s *sumCombiner) AddInput(item uint, acc uint) uint {
	return item + acc
}

// CreateAccumulator implements yenstream.Accumulator.
func (s *sumCombiner) CreateAccumulator() uint {
	return 0
}

type listCombiner struct{}

// AddInput implements yenstream.Accumulator.
func (l *listCombiner) AddInput(item uint, acc []uint) []uint {
	acc = append(acc, item)
	return acc
}

// CreateAccumulator implements yenstream.Accumulator.
func (l *listCombiner) CreateAccumulator() []uint {
	return []uint{}
}

func main() {

	pctx := context.Background()
	// pctx = context.WithValue(pctx, yenstream.DEBUG_NODE, true)
	runCtx := yenstream.NewRunnerContext(pctx)

	runCtx.CreatePipeline(func(ctx *yenstream.RunnerContext) yenstream.Pipeline {
		truesource := yenstream.NewSliceSource(ctx, []uint{
			1,
			2,
			3,
			4,
			5,
			6,
			7,
			8,
			9,
			10,
			11,
		})

		source := truesource.
			Via("flatmapping", yenstream.NewFlatMap(ctx, func(data uint) ([]uint, error) {
				datas := make([]uint, data)
				var c uint = 0
				result := []uint{}
				for range datas {
					c += 1
					result = append(result, c)
				}

				return result, nil
			}))

		keyed := truesource.
			Via("flatmappdata", yenstream.NewFlatMap(ctx, func(data uint) ([]yenstream.KeyedItem[uint], error) {
				datas := make([]yenstream.KeyedItem[uint], data)
				var c uint = 0
				result := []yenstream.KeyedItem[uint]{}
				for range datas {
					c += 1
					dakey := yenstream.NewKeyedItem(fmt.Sprintf("%d", data), c)
					result = append(result, dakey)
				}

				return result, nil
			})).
			Via("combine key", yenstream.NewKeyCombiner(ctx, &listCombiner{})).
			Via("log", yenstream.NewMap(ctx, func(data yenstream.KeyedItem[[]uint]) (yenstream.KeyedItem[[]uint], error) {
				log.Println("combine per key", data.Key(), data.Data())
				return data, nil
			}))

		count := source.
			Via("getting count", yenstream.NewMap(ctx, func(data uint) (uint, error) {
				return 1, nil
			})).
			Via("sum global1", yenstream.NewCombiner(ctx, &sumCombiner{}))
			// Via("log", yenstream.NewMap(ctx, func(data uint) (uint, error) {
			// 	log.Println("datacount", data)
			// 	return data, nil
			// }))

		batch := source.
			Via("log", yenstream.NewMap(ctx, func(data uint) (uint, error) {
				// log.Println(data, "other")
				return data, nil
			})).
			Via("gather with size", yenstream.NewBatch[uint](ctx, 10, time.Minute))
			// Via("log2", yenstream.NewMap(ctx, func(data []uint) ([]uint, error) {
			// 	log.Println(len(data), "len")
			// 	return data, nil
			// }))

		summing := source.
			Via("sum global", yenstream.NewCombiner(ctx, &sumCombiner{}))
			// Via("log", yenstream.NewMap(ctx, func(data uint) (uint, error) {
			// 	log.Println("after summing", data)
			// 	return data, nil
			// }))

		// return summing
		flatall := yenstream.NewFlatten(ctx, batch, summing, count, keyed)
		return flatall

	})

	// sourceErr := yenstream.
	// 	NewChannelSource[error](context.Background())

	// sourceLog := sourceErr.
	// 	Via("log_error", yenstream.NewMap(func(err error) (error, error) {
	// 		log.Println(err)
	// 		return err, nil
	// 	}))

	// go yenstream.Drain(sourceLog)

	// source := yenstream.
	// 	NewSliceSource([]uint{
	// 		1,
	// 		2,
	// 		3,
	// 		4,
	// 		5,
	// 		6,
	// 		7,
	// 	}).
	// 	Via("flatmapping", yenstream.NewFlatMap(func(data uint) ([]uint, error) {
	// 		datas := make([]uint, data)
	// 		var c uint = 0
	// 		result := []uint{}
	// 		for range datas {
	// 			c += 1
	// 			result = append(result, c)

	// 			// if c == 4 {
	// 			// 	sourceErr.Emit(fmt.Errorf("mock 4 error"))
	// 			// }
	// 		}

	// 		return result, nil
	// 	})).
	// 	Via("Kali 2", yenstream.NewMap(func(data uint) (uint, error) {
	// 		// log.Println(data)
	// 		return data * 2, nil
	// 	})).
	// 	Via("Mapping To String", yenstream.NewMap(func(data uint) (string, error) {
	// 		// log.Println(data)
	// 		return fmt.Sprintf("asdasd-%d", data), nil
	// 	})).
	// 	Via("Mapping To String", yenstream.NewMap(func(data string) (string, error) {
	// 		log.Println(data)
	// 		return data, nil
	// 	}))

	// yenstream.Drain(source)

}
