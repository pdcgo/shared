package yenstream

import (
	"log/slog"
)

type Source interface {
	Out() NodeOut
	Via(label string, pipe Pipeline) Pipeline
}

type Pipeline interface {
	Out() NodeOut
	In() chan any
	Via(label string, pipe Pipeline) Pipeline
	SetLabel(label string)
	Process()
}

type Sink interface {
	Out() NodeOut
	Drain()
	Label() string
}

// --------------------------stream ops -------------------------------------

var _ Pipeline = (*mapImpl[any, any])(nil)

type mapImpl[T, R any] struct {
	ctx     *RunnerContext
	label   string
	out     NodeOut
	in      chan any
	handler func(data T) (R, error)
}

// Process implements Outlet.
func (m *mapImpl[T, R]) Process() {
	out := m.out.C()
	defer close(out)

	for data := range m.in {
		res, err := m.handler(data.(T))
		if err != nil {
			slog.Error(err.Error(), slog.String("label", m.label))
			continue
		}
		out <- res
	}
}

// SetLabel implements Pipeline.
func (m *mapImpl[T, R]) SetLabel(label string) {
	m.label = label
}

// In implements Pipeline.
func (m *mapImpl[T, R]) In() chan any {
	return m.in
}

// Out implements Pipeline.
func (m *mapImpl[T, R]) Out() NodeOut {
	return m.out
}

// Via implements Pipeline.
func (m *mapImpl[T, R]) Via(label string, pipe Pipeline) Pipeline {
	m.ctx.RegisterStream(label, m, pipe)

	return pipe
}

func NewMap[T, R any](ctx *RunnerContext, mapper func(data T) (R, error)) *mapImpl[T, R] {
	pipe := &mapImpl[T, R]{
		ctx:     ctx,
		out:     NewNodeOut(ctx),
		in:      make(chan any, 1),
		handler: mapper,
	}
	return pipe
}
