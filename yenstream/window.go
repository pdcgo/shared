package yenstream

import (
	"sync"
	"time"
)

type StateStore interface {
	Get(key any) any
	Set(key any, value any)
	GetAll(emitter func(key any, data any))
}

type WindowType string

const (
	WindowGlobal WindowType = "global"
	WindowChild  WindowType = "child"
)

type Window interface {
	Start() time.Time
	End() time.Time
	WindowType() WindowType
	Store(key string) StateStore
	Emit(data *TimestampedValue)
	Close()
}

type WindowCreatePipe func(rctx *RunnerContext, window Window, source Source) Pipeline

type keyMapStoreImpl struct {
	state map[any]any
}

// GetAll implements StateStore.
func (s *keyMapStoreImpl) GetAll(emitter func(key any, data any)) {
	for key, d := range s.state {
		data := d
		emitter(key, data)
	}
}

// Get implements StateStore.
func (s *keyMapStoreImpl) Get(key any) any {
	return s.state[key]
}

// Set implements StateStore.
func (s *keyMapStoreImpl) Set(key any, value any) {
	s.state[key] = value
}

var _ Pipeline = (*windowIntoImpl)(nil)

type TimestampedValue struct {
	Key  time.Time
	Data any
}

type windowIntoImpl struct {
	ctx      *RunnerContext
	in       chan any
	out      NodeOut
	label    string
	windowfn WindowFunc
	windows  sync.Map
}

func NewWindowInto(ctx *RunnerContext, windowfn WindowFunc) *windowIntoImpl {
	into := &windowIntoImpl{
		ctx:      ctx,
		in:       make(chan any, 1),
		out:      NewNodeOut(ctx),
		windowfn: windowfn,
		windows:  sync.Map{},
	}
	return into
}

// In implements Pipeline.
func (w *windowIntoImpl) In() chan any {
	return w.in
}

// Out implements Pipeline.
func (w *windowIntoImpl) Out() NodeOut {
	return w.out
}

// Process implements Pipeline.
func (w *windowIntoImpl) Process() {
	out := w.out.C()
	defer close(out)

	var wg sync.WaitGroup

Loop:
	for {
		data, ok := <-w.in
		if !ok {
			w.closeWindow()
			break Loop
		}
		tdata := data.(*TimestampedValue)
		wid := w.windowfn.WindowID(tdata)
		var wd any
		wd, ok = w.windows.Load(wid)
		if !ok {
			wd, _ = w.windows.LoadOrStore(wid, w.windowfn.CreateWindow(&wg, w.ctx, wid, out))
		}

		window := wd.(Window)
		window.Emit(tdata)
		// masukkan ke current window
		// tunggu per window
	}

	wg.Wait()

}

// SetLabel implements Pipeline.
func (w *windowIntoImpl) SetLabel(label string) {
	w.label = label
}

// Via implements Pipeline.
func (w *windowIntoImpl) Via(label string, pipe Pipeline) Pipeline {
	w.ctx.RegisterStream(label, w, pipe)
	return pipe
}

func (w *windowIntoImpl) closeWindow() {
	w.windows.Range(func(key, value any) bool {
		window := value.(Window)
		window.Close()

		return true
	})
}
