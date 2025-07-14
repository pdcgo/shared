package yenstream

var _ Pipeline = (*combinerImpl[any, any])(nil)

func NewKeyCombiner[T, R any](ctx *RunnerContext, acc Accumulator[T, R], trigger TriggerFunc) *combinerImpl[T, R] {
	if trigger == nil {
		trigger = NewEmptyTrigger
	}

	combine := &combinerImpl[T, R]{
		ctx:      ctx,
		acc:      acc,
		in:       make(chan any, 1),
		out:      NewNodeOut(ctx),
		globally: false,
		trigger:  trigger,
	}
	return combine
}
