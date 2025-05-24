package streampipe

import "time"

type ID interface {
	GetID() string
}

type TsGroup[T any] struct {
	Ts   time.Time
	Data T
}

func FirstEventInTime[T any](dur time.Duration, groupKey func(item T) (string, time.Time), input <-chan T) <-chan T {
	retc := make(chan T, 3)

	go func() {
		tick := time.NewTicker(dur)
		defer tick.Stop()
		defer close(retc)

		bulk := make(map[string]*TsGroup[T], 0)

		for {
			select {
			case <-tick.C:
				if len(bulk) == 0 {
					continue
				}
				for _, dd := range bulk {
					retc <- dd.Data
				}
				bulk = make(map[string]*TsGroup[T], 0)
				tick.Reset(dur)
			case item, ok := <-input:
				if !ok {
					for _, dd := range bulk {
						retc <- dd.Data
					}
					return
				}

				key, ts := groupKey(item)
				if bulk[key] == nil {
					bulk[key] = &TsGroup[T]{
						Ts:   ts,
						Data: item,
					}

					continue
				}

				if !bulk[key].Ts.Before(ts) {
					bulk[key] = &TsGroup[T]{
						Ts:   ts,
						Data: item,
					}
				}

			}
		}
	}()

	return retc
}

func LastEventInTime[T ID](dur time.Duration, input <-chan T) <-chan map[string]T {
	retc := make(chan map[string]T, 3)

	go func() {
		tick := time.NewTicker(dur)
		defer tick.Stop()
		defer close(retc)

		bulk := make(map[string]T, 0)

		for item := range input {
			select {
			case <-tick.C:
				if len(bulk) == 0 {
					continue
				}
				retc <- bulk
				bulk = make(map[string]T, 0)
				tick.Reset(dur)
			default:
				bulk[item.GetID()] = item
			}
		}

		retc <- bulk
	}()

	return retc
}
