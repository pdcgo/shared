package yenstream

import (
	"sync"
	"time"
)

type WindowType string

const (
	WindowGlobal WindowType = "global"
	WindowChild  WindowType = "child"
)

type Window interface {
	Start() time.Time
	End() time.Time
	WindowType() WindowType
	// Store(key string) store.StateStore
	Emit(data *TimestampedValue)
	Close()
}

type WindowCreatePipe func(rctx *RunnerContext, window Window, source Source) Pipeline

var _ Pipeline = (*windowIntoImpl)(nil)

type TimestampedValue struct {
	Key  time.Time `json:"key"`
	Data HaveMeta  `json:"data"`
}

// GetMeta implements HaveMeta.
func (t *TimestampedValue) GetMeta(key string) any {
	return t.Data.GetMeta(key)
}

// SetMeta implements HaveMeta.
func (t *TimestampedValue) SetMeta(key string, value any) {
	t.Data.SetMeta(key, value)
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
