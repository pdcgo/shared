package streampipe

import "time"

type ID interface {
	GetID() string
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
