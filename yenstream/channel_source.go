package yenstream

var _ Source = (*ChannelSource[any])(nil)

type ChannelSource[T any] struct {
	ctx *RunnerContext
	out NodeOut
	in  chan T
}

// Process implements Outlet.
func (c *ChannelSource[T]) Process() {
	out := c.out.C()
	defer close(out)
	for d := range c.in {
		out <- d
	}
}

// Out implements Source.
func (c *ChannelSource[T]) Out() NodeOut {
	return c.out
}

func (c *ChannelSource[T]) Emit(data T) T {
	go func() {
		c.in <- data
	}()
	return data
}

// Via implements Source.
func (c *ChannelSource[T]) Via(label string, pipe Pipeline) Pipeline {
	c.ctx.RegisterStream(label, c, pipe)
	return pipe
}

func NewChannelSource[T any](ctx *RunnerContext, in chan T) *ChannelSource[T] {
	if in == nil {
		in = make(chan T, 1)
	}
	source := ChannelSource[T]{
		out: NewNodeOut(ctx),
		in:  in,
		ctx: ctx,
	}

	ctx.AddProcess(source.Process)
	return &source
}
