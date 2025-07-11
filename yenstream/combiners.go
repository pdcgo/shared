package yenstream

var GLOBAL_COMBINE = "global"

var _ Pipeline = (*combinerImpl[any, any])(nil)

type Accumulator[T, R any] interface {
	CreateAccumulator() R
	AddInput(item T, acc R) R
}
type combinerImpl[T, R any] struct {
	ctx   *RunnerContext
	acc   Accumulator[T, R]
	in    chan any
	out   NodeOut
	label string
}

func NewCombiner[T, R any](ctx *RunnerContext, acc Accumulator[T, R]) *combinerImpl[T, R] {

	combine := &combinerImpl[T, R]{
		ctx: ctx,
		acc: acc,
		in:  make(chan any, 1),
		out: NewNodeOut(ctx),
	}
	return combine
}

// Process implements Pipeline.
func (c *combinerImpl[T, R]) Process() {
	out := c.out.C()
	defer close(out)

	// getting accumulate store
	window := c.ctx.GetWindow()
	store := window.Store(c.ctx.hash(c.label))

	sacc := store.Get(GLOBAL_COMBINE)
	if sacc == nil {
		sacc = c.acc.CreateAccumulator()
	}

	accu := sacc.(R)
	for data := range c.in {
		accu = c.acc.AddInput(data.(T), accu)
	}

	out <- accu
}

// In implements Pipeline.
func (c *combinerImpl[T, R]) In() chan any {
	return c.in
}

// Out implements Pipeline.
func (c *combinerImpl[T, R]) Out() NodeOut {
	return c.out
}

// SetLabel implements Pipeline.
func (c *combinerImpl[T, R]) SetLabel(label string) {
	c.label = label
}

// Via implements Pipeline.
func (c *combinerImpl[T, R]) Via(label string, pipe Pipeline) Pipeline {
	c.ctx.RegisterStream(label, c, pipe)
	return pipe
}
