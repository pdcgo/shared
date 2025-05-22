package streampipe

import (
	"log/slog"
	"time"
)

type windowAggregateImpl struct {
	dur time.Duration
}

type AggItem struct {
	Key string
	// Data map[]
}

func WindowAggregate[T any, Ag any](
	dur time.Duration,
	input <-chan T,
	getKey func(item T) any,
	handler func(ag Ag, item T) (Ag, error),
) <-chan Ag {
	retc := make(chan Ag, 1)

	go func() {
		var err error
		tick := time.NewTicker(dur)
		defer tick.Stop()
		defer close(retc)

		bulk := make(map[any]Ag)

		for {
			select {
			case <-tick.C:
				if len(bulk) == 0 {
					continue
				}

				for _, d := range bulk {
					retc <- d
				}

				bulk = make(map[any]Ag)
				tick.Reset(dur)

			case item, ok := <-input:
				if !ok {
					for _, d := range bulk {
						retc <- d
					}
					return
				}

				key := getKey(item)
				var ag Ag
				ag, err = handler(bulk[key], item)
				if err != nil {
					slog.Error(err.Error(), slog.String("streampipe", "window aggregate"))
					continue
				}

				bulk[key] = ag
			}
		}
	}()

	return retc
}
