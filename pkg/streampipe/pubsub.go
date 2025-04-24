package streampipe

import "log/slog"

type Event interface {
	EventPath() string
}

type PublishProvider interface {
	Send(topic string, event Event) error
}

func PublishStream[T Event](
	provider PublishProvider,
	topics []string,
	input <-chan T,
	size int,
) <-chan T {

	hasil := make(chan T, size)

	go func() {
		var err error
		defer close(hasil)

		for event := range input {
			for _, topic := range topics {
				err = provider.Send(topic, event)
				if err != nil {
					slog.Error(err.Error())
				}
			}
		}
	}()

	return hasil
}

func ListenStream() {}
