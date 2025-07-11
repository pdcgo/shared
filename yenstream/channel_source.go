package yenstream

var _ Source = (*ChannelSource[any])(nil)

type ChannelSource[T any] struct {
	ctx *RunnerContext
	out NodeOut
}

// Process implements Outlet.
func (c *ChannelSource[T]) Process() {}

// Out implements Source.
func (c *ChannelSource[T]) Out() NodeOut {
	return c.out
}

func (c *ChannelSource[T]) Emit(data T) T {
	go func() {
		c.out.C() <- data
	}()
	return data
}

// Via implements Source.
func (c *ChannelSource[T]) Via(label string, pipe Pipeline) Pipeline {
	c.ctx.RegisterStream(label, c, pipe)
	return pipe
}

func NewChannelSource[T any](ctx *RunnerContext) *ChannelSource[T] {
	source := ChannelSource[T]{
		out: NewNodeOut(ctx),
	}

	go func() {
		<-ctx.Done()
		close(source.out.C())
	}()

	return &source
}
