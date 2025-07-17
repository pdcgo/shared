package store

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
	// print(value)
	s.state[key] = value
}
