package yenstream

var _ Pipeline = (*combinerImpl[any, any])(nil)

type combinerImpl[T, R any] struct {
}

// In implements Pipeline.
func (c *combinerImpl[T, R]) In() chan<- any {
	panic("unimplemented")
}

// Out implements Pipeline.
func (c *combinerImpl[T, R]) Out() <-chan any {
	panic("unimplemented")
}

// SetLabel implements Pipeline.
func (c *combinerImpl[T, R]) SetLabel(label string) {
	panic("unimplemented")
}

// Via implements Pipeline.
func (c *combinerImpl[T, R]) Via(label string, pipe Pipeline) Pipeline {
	panic("unimplemented")
}

func NewCombiner[T, R any]() *combinerImpl[T, R] {
	panic("unimplemented")
}
