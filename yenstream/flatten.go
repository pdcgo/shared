package yenstream

import "sync"

var _ Pipeline = (*flattenImpl)(nil)

type flattenImpl struct {
	ctx       *RunnerContext
	in        chan any
	out       NodeOut
	label     string
	pipelines []Pipeline
}

func NewFlatten(ctx *RunnerContext, label string, pipelines ...Pipeline) *flattenImpl {
	flat := &flattenImpl{
		label:     label,
		ctx:       ctx,
		in:        make(chan any, 1),
		out:       NewNodeOut(ctx),
		pipelines: pipelines,
	}

	var wg sync.WaitGroup

	for _, p := range flat.pipelines {
		pipe := p
		wg.Add(1)
		pipe.Out().Pair(label, flat.In(), func() {
			wg.Done()
		})

	}

	ctx.AddProcess(func() {
		wg.Wait()
		close(flat.in)
	})
	ctx.AddProcess(flat.Process)

	// ctx.AddProcess(flat.Process)
	return flat
}

// In implements Pipeline.
func (f *flattenImpl) In() chan any {
	return f.in
}

// Out implements Pipeline.
func (f *flattenImpl) Out() NodeOut {
	return f.out
}

// Process implements Pipeline.
func (f *flattenImpl) Process() {
	out := f.out.C()
	defer close(out)

	for data := range f.in {
		out <- data
	}

}

// SetLabel implements Pipeline.
func (f *flattenImpl) SetLabel(label string) {
	f.label = label
}

// Via implements Pipeline.
func (f *flattenImpl) Via(label string, pipe Pipeline) Pipeline {
	f.ctx.RegisterStream(label, f, pipe)
	return pipe
}
