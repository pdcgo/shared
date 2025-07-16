package yenstream

import "log/slog"

var _ Pipeline = (*filterImpl[any])(nil)

type filterImpl[T any] struct {
	ctx     *RunnerContext
	label   string
	out     NodeOut
	in      chan any
	handler func(data T) (bool, error)
}

// Process implements Outlet.
func (f *filterImpl[T]) Process() {
	out := f.out.C()
	defer close(out)

	for data := range f.in {
		notskip, err := f.handler(data.(T))
		if err != nil {
			slog.Error(err.Error(), slog.String("label", f.label))
			continue
		}
		if !notskip {
			continue
		}
		out <- data.(T)
	}
}

// SetLabel implements Pipeline.
func (f *filterImpl[T]) SetLabel(label string) {
	f.label = label
}

// In implements Pipeline.
func (f *filterImpl[T]) In() chan any {
	return f.in
}

// Out implements Pipeline.
func (f *filterImpl[T]) Out() NodeOut {
	return f.out
}

// Via implements Pipeline.
func (f *filterImpl[T]) Via(label string, pipe Pipeline) Pipeline {
	f.ctx.RegisterStream(label, f, pipe)

	return pipe
}

func NewFilter[T any](ctx *RunnerContext, mapper func(data T) (bool, error)) *filterImpl[T] {
	pipe := &filterImpl[T]{
		ctx:     ctx,
		out:     NewNodeOut(ctx),
		in:      make(chan any, 1),
		handler: mapper,
	}
	return pipe
}
