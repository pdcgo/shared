package yenstream

import (
	"sync"
)

type Inlet interface {
	In() chan<- any
	SetLabel(label string)
}

type Outlet interface {
	Out() <-chan any
}

func DoStream(label string, fromp Outlet, top Inlet) {
	top.SetLabel(label)
	go func() {
		for element := range fromp.Out() {
			top.In() <- element
		}

		close(top.In())
	}()
}

func Drain(fromps ...Outlet) {
	var wg sync.WaitGroup
	for _, d := range fromps {
		fromp := d
		wg.Add(1)
		go func() {
			defer wg.Done()
			out := fromp.Out()
			for {

				_, ok := <-out
				if !ok {
					break
				}
			}

		}()
	}

	wg.Wait()

}
