package yenstream

import (
	"sync"
	"time"
)

type windowImpl struct {
	sync.Mutex
	stateStore map[string]StateStore
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

// AccumulateStore implements Window.
func (w *windowImpl) Store(key string) StateStore {
	w.Lock()
	defer w.Unlock()

	if w.stateStore[key] == nil {
		w.stateStore[key] = &keyMapStoreImpl{
			state: map[any]any{},
		}
	}

	return w.stateStore[key]
}
