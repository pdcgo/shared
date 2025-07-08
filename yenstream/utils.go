package yenstream

type Inlet interface {
	In() chan<- any
}

type Outlet interface {
	Out() <-chan any
	SetLabel(label string)
}

func DoStream(label string, fromp Outlet, top Inlet) {
	fromp.SetLabel(label)
	go func() {
		for element := range fromp.Out() {
			top.In() <- element
		}

		close(top.In())
	}()
}
