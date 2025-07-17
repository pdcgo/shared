package yenstream

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"
)

type RunnerContext struct {
	ctx         context.Context
	processFunc []func()
	window      Window
}

var WINDOW_KEY = "state_store"

func NewRunnerContext(pctx context.Context) *RunnerContext {
	// var _ Window = (*windowImpl)(nil)
	ctx := &RunnerContext{
		ctx:         pctx,
		processFunc: []func(){},
		window:      &windowImpl{},
	}

	return ctx
}

func (rctx *RunnerContext) GetWindow() Window {
	return rctx.window
}

func (rctx *RunnerContext) CreatePipeline(handler func(ctx *RunnerContext) Pipeline) *RunnerContext {
	pipe := handler(rctx)
	rctx.runProcess()
	res := pipe.
		Out().C()

	for range res {
	}
	return rctx
}

func (rctx *RunnerContext) RegisterStream(label string, fromp Outlet, top Inlet) {
	top.SetLabel(label)
	fromp.Out().Pair(label, top.In())
	rctx.processFunc = append(rctx.processFunc, top.Process)
}

func (rctx *RunnerContext) AddProcess(handler func()) {
	rctx.processFunc = append(rctx.processFunc, handler)
}

func (rctx *RunnerContext) runProcess() {
	// dispatching stream
	for _, pfunc := range rctx.processFunc {
		go pfunc()
	}
}

// Deadline implements context.Context.
func (rctx *RunnerContext) Deadline() (deadline time.Time, ok bool) {
	return rctx.ctx.Deadline()
}

// Done implements context.Context.
func (rctx *RunnerContext) Done() <-chan struct{} {
	return rctx.ctx.Done()
}

// Err implements context.Context.
func (rctx *RunnerContext) Err() error {
	return rctx.ctx.Err()
}

// Value implements context.Context.
func (rctx *RunnerContext) Value(key any) any {
	return rctx.ctx.Value(key)
}

func (rctx *RunnerContext) Hash(data string) string {
	hash := md5.Sum([]byte(data)) // returns [16]byte
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}
