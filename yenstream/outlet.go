package yenstream

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"golang.org/x/exp/slog"
)

type Outlet interface {
	Out() NodeOut
}

type NodeOut interface {
	Pair(label string, in chan any)
	Process()
	C() chan any
}

type nodeOutImpl struct {
	c         chan any
	nodemap   map[string]chan any
	debugMode bool
}

// Process implements NodeOut.
func (n *nodeOutImpl) Process() {
	defer func() {

		for key, mchan := range n.nodemap {
			if n.debugMode {
				slog.Info("closing node in", slog.String("key", key))
			}

			close(mchan)
		}
	}()

	for data := range n.c {
		for _, mchan := range n.nodemap {
			mchan <- data
		}
	}
}

// C implements NodeOut.
func (n *nodeOutImpl) C() chan any {
	return n.c
}

// Pair implements NodeOut.
func (n *nodeOutImpl) Pair(label string, in chan any) {
	var key string
	if n.debugMode {
		key = label
	} else {
		key = n.hash(label)
	}

	isExist := n.nodemap[key] != nil
	if isExist {
		msg := fmt.Sprintf("streaming with label %s is exist", label)
		panic(msg)
	}

	n.nodemap[key] = in

}

// Emit implements NodeOut.
func (n *nodeOutImpl) hash(data string) string {
	hash := md5.Sum([]byte(data)) // returns [16]byte
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

const DEBUG_NODE = "debug_mode"

func NewNodeOut(ctx *RunnerContext) NodeOut {
	var debugMode bool
	debugMode, _ = ctx.Value(DEBUG_NODE).(bool)

	out := &nodeOutImpl{
		nodemap:   map[string]chan any{},
		c:         make(chan any, 1),
		debugMode: debugMode,
	}

	ctx.AddProcess(out.Process)
	return out
}
