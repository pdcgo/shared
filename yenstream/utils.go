package yenstream

type Inlet interface {
	In() chan any
	SetLabel(label string)
	Process()
}

// func DoStream(label string, fromp Outlet, top Inlet) {
// 	top.SetLabel(label)
// 	go func() {
// 		for element := range fromp.Out() {
// 			top.In() <- element
// 		}

// 		close(top.In())
// 	}()
// }
