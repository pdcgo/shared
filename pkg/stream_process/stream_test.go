package stream_process_test

import (
	"context"
	"testing"

	"github.com/pdcgo/shared/pkg/stream_process"
)

func TestStreamChaining(t *testing.T) {
	logseq := stream_process.NewSequence[float32](t.Context(), "debugging",
		stream_process.NewMap("log", func(ctx context.Context, data *stream_process.SequenceItem[float32]) (float32, error) {
			t.Log(data.Data)
			return data.Data, nil
		}),
	)

	logseq2 := stream_process.NewSequence[float32](t.Context(), "debugging",
		stream_process.NewMap("log", func(ctx context.Context, data *stream_process.SequenceItem[float32]) (float32, error) {
			t.Log(data.Data, "asdasdasdasd")
			return data.Data, nil
		}),
	)

	seq := stream_process.NewSequence[int](t.Context(), "test",
		stream_process.NewMap("kali sepuluh", func(ctx context.Context, data *stream_process.SequenceItem[int]) (float32, error) {
			return float32(data.Data) * 10, nil
		}),
		stream_process.NewFilter("filter kurang dari 4", func(ctx context.Context, data *stream_process.SequenceItem[float32]) (bool, error) {
			return data.Data > 40, nil
		}),
		stream_process.NewMap("tambah lima", func(ctx context.Context, data *stream_process.SequenceItem[float32]) (float32, error) {
			return data.Data + 5, nil
		}),
		stream_process.NewPararelSequence[float32]("pararel",
			logseq,
			logseq2,
		),
	)

	for _, c := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11} {
		seq.Process(c)
	}
}
