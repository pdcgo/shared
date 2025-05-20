package streampipe

import (
	"errors"
	"log"
	"log/slog"
)

func SplitMap[T any, K ~string](input <-chan T, keys []K, handle func(item T) ([]K, error)) map[K]chan T {
	// for key := range cha
	mappers := make(map[K]chan T)
	for _, key := range keys {
		chanitem := make(chan T, 1)
		mappers[key] = chanitem
	}

	go func() {
		var err error

		for _, key := range keys {
			// chanitem := make(chan T, 1)
			// mappers[key] = chanitem
			defer close(mappers[key])
		}

		for item := range input {

			var okeys []K
			okeys, err = handle(item)
			if err != nil {
				if errors.Is(err, ErrDropFromStream) {
					continue
				}
				slog.Error(err.Error(), slog.String("pipe", "split_map"))
				continue
			}

			for _, key := range okeys {
				mappers[key] <- item
			}

		}

		log.Println("closing pipe")

	}()

	return mappers
}
