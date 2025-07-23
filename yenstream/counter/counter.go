package counter

import (
	"encoding/json"
	"log"

	"github.com/pdcgo/shared/yenstream"
)

type Counter[T any] interface {
	yenstream.HaveMeta
	Merge(data Counter[T])
	GetValue() T
}

var COUNTER_KEY = "ckey"

func NewLogStateCounter[T any](ctx *yenstream.RunnerContext, pipe yenstream.Pipeline) yenstream.Pipeline {
	state := map[string]Counter[T]{}
	return pipe.
		Via("simple_counter", yenstream.NewMap(ctx, func(data Counter[T]) (Counter[T], error) {

			ckey := data.GetMeta(COUNTER_KEY).(string)
			if state[ckey] == nil {
				state[ckey] = data
			}

			state[ckey].Merge(data)
			raw, _ := json.Marshal(state[ckey])
			log.Println(string(raw))

			return data, nil
		}))
}
