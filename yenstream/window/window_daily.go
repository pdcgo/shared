package window

import "github.com/pdcgo/shared/yenstream"

type dailyWindowImpl struct {
}

// In implements yenstream.Pipeline.
func (d *dailyWindowImpl) In() chan any {
	panic("unimplemented")
}

// Out implements yenstream.Pipeline.
func (d *dailyWindowImpl) Out() yenstream.NodeOut {
	panic("unimplemented")
}

// Process implements yenstream.Pipeline.
func (d *dailyWindowImpl) Process() {
	panic("unimplemented")
}

// SetLabel implements yenstream.Pipeline.
func (d *dailyWindowImpl) SetLabel(label string) {
	panic("unimplemented")
}

// Via implements yenstream.Pipeline.
func (d *dailyWindowImpl) Via(label string, pipe yenstream.Pipeline) yenstream.Pipeline {
	panic("unimplemented")
}
