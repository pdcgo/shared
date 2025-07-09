package yenstream

import "log/slog"

type Source interface {
	Out() <-chan any
	Via(label string, pipe Pipeline) Pipeline
}

type Pipeline interface {
	Out() <-chan any
	In() chan<- any
	Via(label string, pipe Pipeline) Pipeline
	SetLabel(label string)
}

type Sink interface {
	Out() <-chan any
	Drain()
	Label() string
}

// --------------------------stream ops -------------------------------------

var _ Pipeline = (*mapImpl[any, any])(nil)

type mapImpl[T, R any] struct {
	label   string
	out     chan any
	in      chan any
	handler func(data T) (R, error)
}

// SetLabel implements Pipeline.
func (m *mapImpl[T, R]) SetLabel(label string) {
	m.label = label
}

func (m *mapImpl[T, R]) process() {
	defer close(m.out)
	for data := range m.in {
		res, err := m.handler(data.(T))
		if err != nil {
			slog.Error(err.Error(), slog.String("label", m.label))
			continue
		}
		m.out <- res
	}
}

// In implements Pipeline.
func (m *mapImpl[T, R]) In() chan<- any {
	return m.in
}

// Out implements Pipeline.
func (m *mapImpl[T, R]) Out() <-chan any {
	return m.out
}

// Via implements Pipeline.
func (m *mapImpl[T, R]) Via(label string, pipe Pipeline) Pipeline {
	DoStream(label, m, pipe)

	return pipe
}

func NewMap[T, R any](mapper func(data T) (R, error)) *mapImpl[T, R] {
	pipe := &mapImpl[T, R]{
		out:     make(chan any, 1),
		in:      make(chan any, 1),
		handler: mapper,
	}
	go pipe.process()
	return pipe
}

var _ Pipeline = (*flatMapImpl[any, any])(nil)

type flatMapImpl[T, R any] struct {
	label      string
	in         chan any
	out        chan any
	flatmapper func(data T) ([]R, error)
}

// In implements Pipeline.
func (f *flatMapImpl[T, R]) In() chan<- any {
	return f.in
}

// Out implements Pipeline.
func (f *flatMapImpl[T, R]) Out() <-chan any {
	return f.out
}

// SetLabel implements Pipeline.
func (f *flatMapImpl[T, R]) SetLabel(label string) {
	f.label = label
}

// Via implements Pipeline.
func (f *flatMapImpl[T, R]) Via(label string, pipe Pipeline) Pipeline {
	DoStream(label, f, pipe)
	return pipe
}

// process implements Pipeline.
func (f *flatMapImpl[T, R]) process() {
	defer close(f.out)
	for data := range f.in {
		datas, err := f.flatmapper(data.(T))
		for _, data := range datas {
			f.out <- data
		}
		if err != nil {
			slog.Error(err.Error(), slog.String("label", f.label))
			continue
		}

	}
}

func NewFlatMap[T, R any](flatmapper func(data T) ([]R, error)) *flatMapImpl[T, R] {
	flatm := flatMapImpl[T, R]{
		in:         make(chan any, 1),
		out:        make(chan any, 1),
		flatmapper: flatmapper,
	}
	go flatm.process()
	return &flatm
}
