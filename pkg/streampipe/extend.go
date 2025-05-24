package streampipe

import (
	"errors"
	"log/slog"
)

func MapExtend[T any, R any](input <-chan T, handle func(item T) ([]R, error)) <-chan R {
	retc := make(chan R, 1)

	go func() {
		defer close(retc)

		for dd := range input {
			items, err := handle(dd)
			if err != nil {
				if errors.Is(err, ErrDropFromStream) {
					continue
				}
				slog.Error(err.Error(), slog.String("streampipe", "MapExtend"))
				continue
			}
			for _, dat := range items {
				datc := dat
				retc <- datc
			}

		}
	}()

	return retc
}

func MapFilterddd[T any, R any](input <-chan T, handle func(item T) (R, error)) <-chan R {
	retc := make(chan R, 1)
	go func() {
		defer close(retc)
		for dd := range input {

			item, err := handle(dd)
			if err != nil {
				if errors.Is(err, ErrDropFromStream) {
					continue
				}
				slog.Error(err.Error(), slog.String("streampipe", "MapFilter"))
				continue
			}
			retc <- item
		}
	}()

	return retc
}
