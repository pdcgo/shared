package combiner

import (
	"github.com/pdcgo/shared/yenstream"
	"github.com/pdcgo/shared/yenstream/store"
)

type Combiner[T any, R any] interface {
	CreateAccumulator() *yenstream.Row[R]
	AddInput(item T, acc *yenstream.Row[R]) *yenstream.Row[R]
	Key(item T) string
}

func NewKeyCombiner[T any, R any](
	rctx *yenstream.RunnerContext,
	combiner Combiner[T, R],
) *combinerImpl[T, R] {

	return &combinerImpl[T, R]{
		ctx:      rctx,
		combiner: combiner,
		in:       make(chan any, 1),
		out:      yenstream.NewNodeOut(rctx),
	}
}

type combinerImpl[T any, R any] struct {
	label    string
	ctx      *yenstream.RunnerContext
	combiner Combiner[T, R]
	in       chan any
	out      yenstream.NodeOut
}

// In implements yenstream.Pipeline.
func (c *combinerImpl[T, R]) In() chan any {
	return c.in
}

// Out implements yenstream.Pipeline.
func (c *combinerImpl[T, R]) Out() yenstream.NodeOut {
	return c.out
}

// Process implements yenstream.Pipeline.
func (c *combinerImpl[T, R]) Process() {
	state := store.CreateStoreFromCtx(c.ctx, c.ctx.Hash(c.label), c.combiner.CreateAccumulator)
	out := c.out.C()
	defer close(out)

Loop:
	for {
		data, ok := <-c.in
		if !ok {
			break Loop
		}

		tdata := data.(T)

		key := c.combiner.Key(tdata)

		var sacc any
		sacc = state.Get(key)
		if sacc == nil {
			sacc = c.combiner.CreateAccumulator()
		}

		accu := sacc.(*yenstream.Row[R])
		accu = c.combiner.AddInput(tdata, accu)
		accu.SetMeta("key", key)
		state.Set(key, accu)
	}

	state.GetAll(func(key, val any) {
		// log.Println(key, val)
		dsend := val.(*yenstream.Row[R])
		out <- dsend
	})
}

// SetLabel implements yenstream.Pipeline.
func (c *combinerImpl[T, R]) SetLabel(label string) {
	c.label = label
}

// Via implements yenstream.Pipeline.
func (c *combinerImpl[T, R]) Via(label string, pipe yenstream.Pipeline) yenstream.Pipeline {
	c.ctx.RegisterStream(label, c, pipe)
	return pipe
}
