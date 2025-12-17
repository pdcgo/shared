package stream_process

import "context"

type PararelSequence interface {
	getNextHandler() SequenceNext
}

func NewPararelSequence[Param any](name string, sequences ...PararelSequence) SequenceHandler[Param] {
	handlers := []SequenceNext{}
	for _, seq := range sequences {
		handlers = append(handlers, seq.getNextHandler())
	}

	return func(ctx context.Context, next SequenceNext) SequenceNext {
		return func(ctx context.Context, data ProcessItem) (ProcessItem, error) {
			var err error
			for _, handler := range handlers {
				_, err = handler(ctx, data)
				if err != nil {
					return nil, NewStreamErr(name, err)
				}
			}
			return nil, nil
		}
	}
}
