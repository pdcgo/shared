package yenstream

import "context"

var _ Source = (*ChannelSource[any])(nil)

type ChannelSource[T any] struct {
	out chan any
}

func (c *ChannelSource[T]) Emit(data T) T {
	go func() {
		c.out <- data
	}()
	return data
}

// Out implements Source.
func (c *ChannelSource[T]) Out() <-chan any {
	return c.out
}

// Via implements Source.
func (c *ChannelSource[T]) Via(label string, pipe Pipeline) Pipeline {
	DoStream(label, c, pipe)
	return pipe
}

func NewChannelSource[T any](ctx context.Context) *ChannelSource[T] {
	source := ChannelSource[T]{
		out: make(chan any, 1),
	}

	go func() {
		<-ctx.Done()
		close(source.out)
	}()

	return &source
}
