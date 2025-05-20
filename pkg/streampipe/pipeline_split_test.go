package streampipe_test

import (
	"testing"

	"github.com/pdcgo/shared/pkg/streampipe"
	"github.com/stretchr/testify/assert"
)

type MapDirect string

const (
	TestMap MapDirect = "test"
	SlowMap MapDirect = "slow"
)

func TestSplitMap(t *testing.T) {

	input := make(chan string)

	go func() {
		defer close(input)

		datas := []string{"all", "1", "2", "drop"}

		for _, data := range datas {
			select {
			case input <- data:
			case <-t.Context().Done():
			}
		}

	}()

	mappers := streampipe.SplitMap(input, []MapDirect{SlowMap, TestMap}, func(item string) ([]MapDirect, error) {

		switch item {
		case "all":
			return []MapDirect{SlowMap, TestMap}, nil
		case "1":
			return []MapDirect{SlowMap}, nil
		case "2":
			return []MapDirect{TestMap}, nil
		case "drop":
			return []MapDirect{}, streampipe.ErrDropFromStream
		}

		return []MapDirect{}, nil
	})

	all := <-mappers[SlowMap]
	assert.Equal(t, "all", all)
	all = <-mappers[TestMap]
	assert.Equal(t, "all", all)

	all = <-mappers[SlowMap]
	assert.Equal(t, "1", all)

	all = <-mappers[TestMap]
	assert.Equal(t, "2", all)

}
