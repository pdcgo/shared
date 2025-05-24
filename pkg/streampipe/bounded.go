package streampipe

import (
	"errors"
	"log/slog"
)

func WithPrevious[T any](input <-chan T, handle func(prev T, item T) error) <-chan T {
	retc := make(chan T, 3)
	go func() {
		var err error
		defer close(retc)
		var lastitem T
		for item := range input {
			err = handle(lastitem, item)
			if err != nil {
				if errors.Is(err, ErrDropFromStream) {
					continue
				}
				slog.Error(err.Error(), slog.String("lib", "streampipe"), slog.String("method", "sink"))
				continue
			}
			retc <- item
			lastitem = item
		}
	}()

	return retc
}
