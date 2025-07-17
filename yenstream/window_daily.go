package yenstream

import (
	"fmt"
	"sync"
	"time"
)

type WindowFunc interface {
	WindowID(data *TimestampedValue) int64
	CreateWindow(wg *sync.WaitGroup, wctx *RunnerContext, key int64, in chan any) Window
}

type dailyWindowFunc struct {
	phandler WindowCreatePipe
}

// CreateWindow implements WindowFunc.
func (d *dailyWindowFunc) CreateWindow(wg *sync.WaitGroup, rctx *RunnerContext, key int64, out chan any) Window {
	inChan := make(chan any, 1)

	window := &dailyWindow{
		in:    inChan,
		label: fmt.Sprintf("%d", key),
		start: time.Unix(0, key*1000),
	}

	// var ww Window = window

	ctx := &RunnerContext{
		ctx:         rctx,
		processFunc: []func(){},
		window:      window,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx.CreatePipeline(func(ctx *RunnerContext) Pipeline {

			source := NewChannelSource(ctx, inChan)
			return d.
				phandler(ctx, window, source).
				Via("daily_window_sendout", NewMap(ctx, func(data any) (any, error) {
					out <- data
					return data, nil
				}))
			// return source.
			// 	Via("asdd", NewMap(ctx, func(data any) (any, error) {
			// 		return data, nil
			// 	})).
			// 	Via("loasdasdg", NewMap(ctx, func(data any) (any, error) {
			// 		// log.Println("logging from window"+window.label, data)
			// 		out <- data
			// 		return data, nil
			// 	}))
		})
	}()

	return window
}

func DailyWindow(phandler WindowCreatePipe) *dailyWindowFunc {
	return &dailyWindowFunc{
		phandler: phandler,
	}
}

// WindowID implements WindowFunc.
func (d *dailyWindowFunc) WindowID(data *TimestampedValue) int64 {
	t := data.Key.UTC()
	year, month, day := t.Date()

	nt := time.Time{}
	nt = nt.AddDate(year, int(month), day).UTC()
	return nt.UnixMicro()
}

type dailyWindow struct {
	start time.Time
	in    chan any
	label string
}

// WindowType implements Window.
func (d *dailyWindow) WindowType() WindowType {
	return WindowChild
}

// End implements Window.
func (d *dailyWindow) End() time.Time {
	return d.start.AddDate(0, 0, 1)
}

// Start implements Window.
func (d *dailyWindow) Start() time.Time {
	return d.start
}

// Close implements Window.
func (d *dailyWindow) Close() {
	close(d.in)
}

// Emit implements Window.
func (d *dailyWindow) Emit(data *TimestampedValue) {
	d.in <- data
}
