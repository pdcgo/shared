package yenstream

import (
	"log/slog"
)

var _ Pipeline = (*flatMapImpl[any, any])(nil)

type flatMapImpl[T, R any] struct {
	ctx        *RunnerContext
	label      string
	in         chan any
	out        NodeOut
	flatmapper func(data T) ([]R, error)
}

// In implements Pipeline.
func (f *flatMapImpl[T, R]) In() chan any {
	return f.in
}

// Out implements Pipeline.
func (f *flatMapImpl[T, R]) Out() NodeOut {
	return f.out
}

// SetLabel implements Pipeline.
func (f *flatMapImpl[T, R]) SetLabel(label string) {
	f.label = label
}

// Via implements Pipeline.
func (f *flatMapImpl[T, R]) Via(label string, pipe Pipeline) Pipeline {
	f.ctx.RegisterStream(label, f, pipe)
	return pipe
}

// process implements Pipeline.
func (f *flatMapImpl[T, R]) Process() {
	out := f.out.C()
	defer close(out)

	for data := range f.in {

		datas, err := f.flatmapper(data.(T))
		for _, data := range datas {
			out <- data
		}
		if err != nil {
			slog.Error(err.Error(), slog.String("label", f.label))
			continue
		}

	}
}

func NewFlatMap[T, R any](ctx *RunnerContext, flatmapper func(data T) ([]R, error)) *flatMapImpl[T, R] {
	flatm := flatMapImpl[T, R]{
		ctx:        ctx,
		in:         make(chan any, 1),
		out:        NewNodeOut(ctx),
		flatmapper: flatmapper,
	}
	return &flatm
}
