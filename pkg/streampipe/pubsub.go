package streampipe

import (
	"log/slog"
)

type Event interface {
	EventPath() string
}

type PullEvent interface {
	EventPath() string
	Ack()
	Decode(v any) error
}

type PublishProvider interface {
	Send(topic string, event Event) error
	Close() error
}

type PullProvider interface {
	Receive(handler func(event PullEvent)) error
	Close() error
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

func PullStream(topic string, size int) <-chan PullEvent {
	hasil := make(chan PullEvent, size)

	return hasil

}
