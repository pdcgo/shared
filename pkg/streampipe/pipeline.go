package streampipe

import (
	"sync"
	"time"
)

func Merge[T any](inputs ...<-chan T) <-chan T {
	retc := make(chan T, 9)

	var wg sync.WaitGroup

	for _, dd := range inputs {
		input := dd
		wg.Add(1)

		go func() {
			defer wg.Done()
			for item := range input {
				retc <- item
			}
		}()
	}

	go func() {
		wg.Wait()
		close(retc)
	}()

	return retc
}

func Split[T any](input <-chan T, splitSize int) []<-chan T {
	retc := make([]chan T, splitSize)
	i := 0
	for {
		retc[i] = make(chan T, 5)

		i += 1
		if i == splitSize {
			break
		}
	}

	go func() {
		defer func() {
			for _, ret := range retc {
				close(ret)
			}
		}()

		for item := range input {
			for _, ret := range retc {
				ret <- item
			}
		}
	}()

	hasil := make([]<-chan T, len(retc))
	for i, cc := range retc {
		hasil[i] = cc
	}
	return hasil
}

func Release[T any](input <-chan T) {
	for range input {
	}
}

func Sink[T any](input <-chan T, handle func(item T)) <-chan T {
	retc := make(chan T, 3)
	go func() {
		defer close(retc)
		for item := range input {
			handle(item)
			retc <- item
		}
	}()

	return retc
}

func Filter[T any](input <-chan T, handle func(item T) bool) <-chan T {
	retc := make(chan T, 3)
	go func() {
		defer close(retc)
		for item := range input {
			if handle(item) {
				retc <- item
			}

		}
	}()

	return retc
}
func Map[T any, R any](input <-chan T, handle func(item T) R) <-chan R {
	retc := make(chan R, 3)
	go func() {
		defer close(retc)
		for dd := range input {
			item := dd
			retc <- handle(item)
		}
	}()

	return retc
}

func UnSlice[T any](input <-chan []T) <-chan T {
	retc := make(chan T, 3)
	go func() {
		defer close(retc)
		for items := range input {
			for _, dd := range items {
				item := dd
				retc <- item
			}
		}
	}()

	return retc
}

func Unique[T any](input <-chan []T, handle func(item T) string) <-chan []T {
	retc := make(chan []T, 3)

	go func() {
		defer close(retc)
		for items := range input {
			bulk := make(map[string]T, 0)
			for _, dd := range items {
				item := dd
				key := handle(item)
				if key == "" {
					continue
				}
				bulk[key] = item
			}

			datas := make([]T, len(bulk))
			i := 0
			for _, val := range bulk {
				datas[i] = val
				i += 1
			}

			retc <- datas
			bulk = make(map[string]T, 0)

		}
	}()

	return retc
}

func TimeWindow[T any](dur time.Duration, input <-chan T) <-chan []T {
	retc := make(chan []T, 3)

	go func() {
		tick := time.NewTicker(dur)
		defer tick.Stop()
		defer close(retc)

		bulk := make([]T, 0)

		for item := range input {
			select {
			case <-tick.C:
				if len(bulk) == 0 {
					continue
				}

				retc <- bulk
				bulk = make([]T, 0)
				tick.Reset(dur)
			default:
				bulk = append(bulk, item)
			}
		}

		retc <- bulk
	}()

	return retc
}
