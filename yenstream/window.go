package yenstream

import "sync"

type StateStore interface {
	Get(key any) any
	Set(key any, value any)
}

type Window interface {
	Store(key string) StateStore
}

type windowImpl struct {
	sync.Mutex
	stateStore map[string]StateStore
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

type keyMapStoreImpl struct {
	state map[any]any
}

// Get implements StateStore.
func (s *keyMapStoreImpl) Get(key any) any {
	return s.state[key]
}

// Set implements StateStore.
func (s *keyMapStoreImpl) Set(key any, value any) {
	s.state[key] = value
}
