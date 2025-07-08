package main

import (
	"fmt"
	"log"

	"github.com/pdcgo/shared/yenstream"
)

func main() {
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
			}

			return result, nil
		})).
		Via("Mapping To String", yenstream.NewMap(func(data uint) (string, error) {
			return fmt.Sprintf("asdasd-%d", data), nil
		}))

	for data := range source.Out() {
		log.Println(data)
	}
}
