package main

import (
	"context"
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

		window := truesource.
			Via("flatmapping", yenstream.NewFlatMap(ctx, func(data uint) ([]*yenstream.TimestampedValue, error) {
				datas := make([]uint, data)
				var c uint = 0
				result := []*yenstream.TimestampedValue{}
				for range datas {
					c += 1
					result = append(result, &yenstream.TimestampedValue{
						Key:  time.Now().AddDate(0, 0, int(data)),
						Data: c,
					})
				}

				return result, nil
			})).
			Via("windowing", yenstream.NewWindowInto(ctx, yenstream.DailyWindow(func(rctx *yenstream.RunnerContext, window yenstream.Window, source yenstream.Source) yenstream.Pipeline {
				return source.
					Via("map get c", yenstream.NewMap(rctx, func(data *yenstream.TimestampedValue) (uint, error) {
						c := (data.Data).(uint)
						return c, nil
					})).
					Via("combine all", yenstream.NewCombiner(rctx, &sumCombiner{})).
					Via("log combine", yenstream.NewMap(rctx, func(data uint) (uint, error) {

						log.Println(window.Start(), data)
						return data, nil
					}))

			}))).
			Via("after windowing", yenstream.NewMap(ctx, func(data uint) (uint, error) {
				log.Println("after windowing", data)
				return data, nil
			}))

		// return summing
		flatall := yenstream.NewFlatten(ctx, window)
		return flatall

	})

}
