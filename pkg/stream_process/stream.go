package stream_process

import (
	"context"
	"fmt"
)

type Stream[Param, Result any] interface {
	Process(ctx context.Context, T Param) (Result, error)
}

type ProcessItem interface {
	Any() any
}

type SequenceNext func(ctx context.Context, data ProcessItem) (ProcessItem, error)
type SequenceHandler[Result any] func(ctx context.Context, next SequenceNext) SequenceNext

type SequenceItem[T any] struct {
	Data T
}

func (s SequenceItem[T]) Any() any {
	return s.Data
}

type Sequence[Param, Result any] struct {
	ctx         context.Context
	Name        string
	handler     []SequenceHandler[Result]
	nextHandler SequenceNext
}

func (seq *Sequence[Param, Result]) getNextHandler() SequenceNext {
	return seq.nextHandler
}

func (seq *Sequence[Param, Result]) Process(data Param) (*SequenceItem[Result], error) {
	if seq.handler == nil {
		return nil, fmt.Errorf("sequence %s have no handler", seq.Name)
	}

	res, err := seq.nextHandler(seq.ctx, &SequenceItem[Param]{
		Data: data,
	})

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	resdata := res.Any()
	return &SequenceItem[Result]{
		Data: resdata.(Result),
	}, err
}

func NewSequence[Param, Result interface{}](
	ctx context.Context,
	name string,
	handlers ...SequenceHandler[Result],
) *Sequence[Param, Result] {
	Reverse(handlers)

	var next SequenceNext = func(ctx context.Context, data ProcessItem) (ProcessItem, error) {
		return data, nil
	}

	for _, h := range handlers {
		next = h(ctx, next)
	}

	return &Sequence[Param, Result]{
		ctx,
		name,
		handlers,
		next,
	}
}

func NewMap[Param, Result any](name string, mapfunc func(ctx context.Context, data *SequenceItem[Param]) (Result, error)) SequenceHandler[Result] {
	return func(ctx context.Context, next SequenceNext) SequenceNext {

		return func(ctx context.Context, data ProcessItem) (ProcessItem, error) {
			res, err := mapfunc(ctx, data.(*SequenceItem[Param]))
			if err != nil {
				return nil, NewStreamErr(name, err)
			}

			return next(ctx, &SequenceItem[Result]{
				Data: res,
			})
		}
	}
}

func NewFilter[Result any](name string, mapfunc func(ctx context.Context, data *SequenceItem[Result]) (bool, error)) SequenceHandler[Result] {
	return func(ctx context.Context, next SequenceNext) SequenceNext {

		return func(ctx context.Context, data ProcessItem) (ProcessItem, error) {
			item := data.(*SequenceItem[Result])
			ok, err := mapfunc(ctx, item)
			if err != nil {
				return nil, NewStreamErr(name, err)
			}
			if !ok {
				return nil, nil
			}

			return next(ctx, item)
		}
	}
}

type StreamErr struct {
	name string
	err  error
}

// Error implements error.
func (s *StreamErr) Error() string {
	return fmt.Sprintf("streaming error on %s with %s", s.name, s.err)
}

func NewStreamErr(name string, err error) *StreamErr {
	return &StreamErr{
		name,
		err,
	}
}

func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
