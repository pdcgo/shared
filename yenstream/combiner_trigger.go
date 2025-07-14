package yenstream

import (
	"time"
)

type durationTriggerImpl struct {
	sendItem func(key, data any)
	data     map[any]any
	duration time.Duration
	close    chan int
}

func NewDurationTrigger(duration time.Duration) func(store StateStore, window Window, sendItem func(key, data any)) TriggerCombiner {
	return func(store StateStore, window Window, sendItem func(key any, data any)) TriggerCombiner {
		durTrigger := &durationTriggerImpl{
			sendItem: sendItem,
			data:     map[any]any{},
			duration: duration,
			close:    make(chan int, 1),
		}

		return durTrigger
	}
}

// Close implements TriggerCombiner.
func (d *durationTriggerImpl) Close() {
	d.close <- 1
	close(d.close)
}

// Emit implements TriggerCombiner.
func (d *durationTriggerImpl) Emit(key any, data any) {
	d.data[key] = data

}

// Process implements TriggerCombiner.
func (d *durationTriggerImpl) Process() {
	tick := time.NewTicker(d.duration)
	for {
		select {
		case <-tick.C:
			for key, data := range d.data {
				d.sendItem(key, data)
			}
			d.data = map[any]any{}
		case _, ok := <-d.close:
			if !ok {
				return
			}
		}
	}
}

type EmptyTrigger struct{}

func NewEmptyTrigger(store StateStore, window Window, sendItem func(key, data any)) TriggerCombiner {
	return &EmptyTrigger{}
}

// Emit implements TriggerCombiner.
func (e *EmptyTrigger) Emit(key any, data any) {
}

// Close implements TriggerCombiner.
func (e *EmptyTrigger) Close() {

}

// Process implements TriggerCombiner.
func (e *EmptyTrigger) Process() {

}
