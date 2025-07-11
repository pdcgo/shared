package yenstream

var _ Source = (*sliceSource[any])(nil)

type sliceSource[T any] struct {
	ctx   *RunnerContext
	label string
	out   NodeOut
	datas []T
}

// NodeOut implements Outlet.
func (s *sliceSource[T]) Out() NodeOut {
	return s.out
}

// Process implements Outlet.
func (s *sliceSource[T]) Process() {
	in := s.out.C()
	defer close(in)

	for _, data := range s.datas {
		in <- data
	}
}

// SetLabel implements Outlet.
func (s *sliceSource[T]) SetLabel(label string) {
	s.label = label
}

// func (s *sliceSource[T]) process(datas []T) {
// 	defer close(s.out)

// 	for _, data := range datas {
// 		s.out <- data
// 	}
// }

// Via implements Source.
func (s *sliceSource[T]) Via(label string, pipe Pipeline) Pipeline {
	s.ctx.RegisterStream(label, s, pipe)
	return pipe
}

func NewSliceSource[T any](ctx *RunnerContext, datas []T) *sliceSource[T] {
	s := sliceSource[T]{
		ctx:   ctx,
		out:   NewNodeOut(ctx),
		datas: datas,
		label: "slice_source",
	}
	ctx.AddProcess(s.Process)
	return &s
}
