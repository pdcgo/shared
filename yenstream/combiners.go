package yenstream

import (
	"log/slog"
	"reflect"
)

var GLOBAL_COMBINE = "global"

var _ Pipeline = (*combinerImpl[any, HaveMeta])(nil)

type CombinerValue[T HaveMeta] struct {
	Key    any
	Data   T    `json:"data"`
	Final  bool `json:"final"`
	window Window
}

// GetMeta implements HaveMeta.
func (c *CombinerValue[T]) GetMeta(key string) any {
	return c.Data.GetMeta(key)
}

// SetMeta implements HaveMeta.
func (c *CombinerValue[T]) SetMeta(key string, value any) {
	c.Data.SetMeta(key, value)
}

func (c *CombinerValue[T]) GetWindow() Window {
	return c.window
}

type TriggerCombiner interface {
	Emit(key, data any)
	Process()
	Close()
}

type TriggerFunc func(store StateStore, window Window, sendItem func(key, data any)) TriggerCombiner

type Accumulator[T, R any] interface {
	CreateAccumulator() R
	AddInput(item T, acc R) R
}
type combinerImpl[T any, R HaveMeta] struct {
	ctx      *RunnerContext
	acc      Accumulator[T, R]
	in       chan any
	out      NodeOut
	label    string
	globally bool
	trigger  TriggerFunc
}

func NewCombiner[T any, R HaveMeta](ctx *RunnerContext, acc Accumulator[T, R], trigger TriggerFunc) *combinerImpl[T, R] {
	if trigger == nil {
		trigger = NewEmptyTrigger
	}

	combine := &combinerImpl[T, R]{
		ctx:      ctx,
		acc:      acc,
		in:       make(chan any, 1),
		out:      NewNodeOut(ctx),
		globally: true,
		trigger:  trigger,
	}
	return combine
}

// Process implements Pipeline.
func (c *combinerImpl[T, R]) Process() {
	out := c.out.C()
	defer close(out)
	// getting accumulate store
	window := c.ctx.GetWindow()
	store := window.Store(c.ctx.hash(c.label))

	sendItemFinal := func(key, data any) {
		var dsend CombinerValue[R]
		dsend.Key = key
		dsend.Data = data.(R)
		dsend.Final = true
		dsend.window = window
		out <- &dsend
	}

	sendItem := func(key, data any) {
		var dsend CombinerValue[R]
		dsend.Key = key
		dsend.Data = data.(R)
		dsend.window = window
		out <- &dsend
	}
	defer store.GetAll(sendItemFinal)

	trigger := c.trigger(store, window, sendItem)
	defer trigger.Close()
	go trigger.Process()

Loop:
	for {
		data, ok := <-c.in
		if !ok {
			break Loop
		}
		// getting accumulator
		var sacc any
		var key any
		var pdata any

		tdata, ok := data.(*TimestampedValue)
		if ok {
			pdata = tdata.Data
		} else {
			pdata = data
		}

		if c.globally {
			key = GLOBAL_COMBINE
			if window.WindowType() != WindowGlobal {
				key = window.Start().UnixMicro()
			}
		} else {

			dkey, ok := pdata.(KeyedItem[T])
			if !ok {
				name := reflect.TypeOf(pdata).Elem().Name()
				var need T
				needname := reflect.TypeOf(need).Elem().Name()
				slog.Error("type error", slog.String("type", name), slog.String("type_need", needname), slog.String("pipe", c.label))
				panic("pipe " + c.label + " item not keyed data")
			}

			key = dkey.Key()
		}

		sacc = store.Get(key)

		if sacc == nil {
			sacc = c.acc.CreateAccumulator()
		}

		accu := sacc.(R)
		if c.globally {
			accu = c.acc.AddInput(pdata.(T), accu)
		} else {
			kdata := pdata.(KeyedItem[T])
			accu = c.acc.AddInput(kdata.Data(), accu)
		}

		store.Set(key, accu)
		trigger.Emit(key, accu)
		// sending accumulator

	}

}

// In implements Pipeline.
func (c *combinerImpl[T, R]) In() chan any {
	return c.in
}

// Out implements Pipeline.
func (c *combinerImpl[T, R]) Out() NodeOut {
	return c.out
}

// SetLabel implements Pipeline.
func (c *combinerImpl[T, R]) SetLabel(label string) {
	c.label = label
}

// Via implements Pipeline.
func (c *combinerImpl[T, R]) Via(label string, pipe Pipeline) Pipeline {
	c.ctx.RegisterStream(label, c, pipe)
	return pipe
}
