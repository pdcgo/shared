package yenstream

var _ Pipeline = (*combinerImpl[any, any])(nil)

func NewKeyCombiner[T, R any](ctx *RunnerContext, acc Accumulator[T, R]) *combinerImpl[T, R] {

	combine := &combinerImpl[T, R]{
		ctx:      ctx,
		acc:      acc,
		in:       make(chan any, 1),
		out:      NewNodeOut(ctx),
		globally: false,
	}
	return combine
}
