package yenstream

var _ Source = (*sliceSource[any])(nil)

type sliceSource[T any] struct {
	label string
	out   chan any
}

// SetLabel implements Outlet.
func (s *sliceSource[T]) SetLabel(label string) {
	s.label = label
}

func (s *sliceSource[T]) process(datas []T) {
	defer close(s.out)

	for _, data := range datas {
		s.out <- data
	}
}

// Out implements Source.
func (s *sliceSource[T]) Out() <-chan any {
	return s.out
}

// Via implements Source.
func (s *sliceSource[T]) Via(label string, pipe Pipeline) Pipeline {
	DoStream(label, s, pipe)
	return pipe
}

func NewSliceSource[T any](datas []T) *sliceSource[T] {
	s := sliceSource[T]{
		out: make(chan any, 1),
	}
	go s.process(datas)
	return &s
}
