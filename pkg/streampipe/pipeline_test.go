package streampipe_test

import (
	"strconv"
	"testing"
	"time"

	pipeline "github.com/pdcgo/shared/pkg/streampipe"
)

func chanDumpInt() <-chan int {
	dump := make(chan int, 4)
	go func() {
		defer close(dump)
		dump <- 1
		dump <- 2
		dump <- 3
		dump <- 4
	}()

	return dump
}

func chanDumpIntSlice() <-chan []int {
	dump := make(chan []int, 4)
	go func() {
		defer close(dump)
		dump <- []int{1}
		dump <- []int{2}
		dump <- []int{3}
		dump <- []int{4}
	}()

	return dump
}

func TestMergeAndSplitPipeline(t *testing.T) {
	dump := chanDumpInt()
	splitp := pipeline.Split(dump, 2)

	endprocess := pipeline.Merge(splitp[0], splitp[1])

	pipeline.Release(endprocess)
}

func TestSinkPipeline(t *testing.T) {

	dump := chanDumpInt()
	endprocess := pipeline.Sink(dump, func(d int) {})
	pipeline.Release(endprocess)

}

func TestFilterPipeline(t *testing.T) {

	dump := chanDumpInt()
	endprocess := pipeline.Filter(dump, func(d int) bool { return true })
	pipeline.Release(endprocess)

}

func TestMapPipeline(t *testing.T) {

	t.Run("testing channel close map", func(t *testing.T) {
		dump := chanDumpInt()
		endprocess := pipeline.Map(dump, func(d int) int {
			return d * 2
		})
		pipeline.Release(endprocess)
	})

}

func TestUnslicePipeline(t *testing.T) {

	dump := chanDumpIntSlice()
	endprocess := pipeline.UnSlice(dump)
	pipeline.Release(endprocess)

}

func TestUniquePipeline(t *testing.T) {

	dump := chanDumpIntSlice()
	endprocess := pipeline.Unique(dump, func(item int) string { return strconv.FormatInt(int64(item), 10) })
	pipeline.Release(endprocess)

}

func TestTimewindowPipeline(t *testing.T) {

	dump := chanDumpInt()
	endprocess := pipeline.TimeWindow(time.Second, dump)
	pipeline.Release(endprocess)

}
