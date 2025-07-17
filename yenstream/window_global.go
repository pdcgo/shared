package yenstream

import (
	"context"
	"time"
)

type windowImpl struct {
	ctx context.Context
}

// WindowType implements Window.
func (w *windowImpl) WindowType() WindowType {
	return WindowGlobal
}

// End implements Window.
func (w *windowImpl) End() time.Time {
	return time.Now()
}

// Start implements Window.
func (w *windowImpl) Start() time.Time {
	return time.Time{}
}

// Close implements Window.
func (w *windowImpl) Close() {}

// Emit implements Window.
func (w *windowImpl) Emit(data *TimestampedValue) {
	panic("not supported in global window")
}
