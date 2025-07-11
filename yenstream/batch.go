package yenstream

import (
	"time"
)

type batchImpl[T any] struct {
	ctx   *RunnerContext
	in    chan any
	out   NodeOut
	label string

	maxDuration time.Duration
	size        int
}

var _ Pipeline = (*batchImpl[any])(nil)

func NewBatch[T any](ctx *RunnerContext, size int, maxDuration time.Duration) *batchImpl[T] {

	batch := &batchImpl[T]{
		ctx:         ctx,
		in:          make(chan any, 1),
		out:         NewNodeOut(ctx),
		size:        size,
		maxDuration: maxDuration,
	}

	return batch
}

// In implements Pipeline.
func (b *batchImpl[T]) In() chan any {
	return b.in
}

// Out implements Pipeline.
func (b *batchImpl[T]) Out() NodeOut {
	return b.out
}

// Process implements Pipeline.
func (b *batchImpl[T]) Process() {
	out := b.out.C()
	defer close(out)

	tick := time.NewTicker(b.maxDuration)
	defer tick.Stop()

	bulk := make([]T, 0)
	c := 0
	reset := func() {
		c = 0
		bulk = make([]T, 0)

	}

Parent:
	for {
		select {
		case <-tick.C:
			if len(bulk) != 0 {
				out <- bulk
			}
			reset()
		default:
			data, ok := <-b.in
			if !ok {
				break Parent
			}

			bulk = append(bulk, data.(T))
			c += 1
			if c >= b.size {
				out <- bulk
				reset()
			}
		}
	}

	if len(bulk) != 0 {
		out <- bulk
	}
}

// SetLabel implements Pipeline.
func (b *batchImpl[T]) SetLabel(label string) {
	b.label = label
}

// Via implements Pipeline.
func (b *batchImpl[T]) Via(label string, pipe Pipeline) Pipeline {
	b.ctx.RegisterStream(label, b, pipe)
	return pipe
}
