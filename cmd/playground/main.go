package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pdcgo/shared/yenstream"
)

func main() {
	sourceErr := yenstream.
		NewChannelSource[error](context.Background())

	sourceLog := sourceErr.
		Via("log_error", yenstream.NewMap(func(err error) (error, error) {
			log.Println(err)
			return err, nil
		}))

	go yenstream.Drain(sourceLog)

	source := yenstream.
		NewSliceSource([]uint{
			1,
			2,
			3,
			4,
			5,
			6,
			7,
		}).
		Via("flatmapping", yenstream.NewFlatMap(func(data uint) ([]uint, error) {
			datas := make([]uint, data)
			var c uint = 0
			result := []uint{}
			for range datas {
				c += 1
				result = append(result, c)

				if c == 4 {
					sourceErr.Emit(fmt.Errorf("mock 4 error"))
				}
			}

			return result, nil
		})).
		Via("Kali 2", yenstream.NewMap(func(data uint) (uint, error) {
			// log.Println(data)
			return data * 2, nil
		})).
		Via("Mapping To String", yenstream.NewMap(func(data uint) (string, error) {
			// log.Println(data)
			return fmt.Sprintf("asdasd-%d", data), nil
		})).
		Via("Mapping To String", yenstream.NewMap(func(data string) (string, error) {
			log.Println(data)
			return data, nil
		}))

	yenstream.Drain(source)

}
