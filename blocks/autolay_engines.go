package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/internal/core/node"
	"github.com/libanvl/swager/pkg/ipc"
)

func (a *Autolay) autoTiler(wct ipc.WindowChangeType, ws *ipc.Node) error {
	if core.Deny(wct, ipc.TitleWindow) {
		return nil
	}

	cwin := node.Count(ws, node.IsLeaf)
	is_even := (cwin % 2) == 0

	if a.LogLevel.Debug() {
		a.Opts.Log.Printf("autolay", "{autotiler} ws.name:  %v, cwin:  %v, is_even: %v", ws.Name, cwin, is_even)
	}

	if is_even {
		a.Command("autotiler", "splitv")
	} else {
		a.Command("autotiler", "splith")
	}
	return nil
}

func (a *Autolay) masterStack(wct ipc.WindowChangeType, ws *ipc.Node) error {
	if core.Deny(wct, ipc.TitleWindow, ipc.FocusWindow) {
		return nil
	}

	cwin := node.Count(ws, node.IsLeaf)
	switch {
	case cwin == 1:
		a.Command("masterstack", "splith")
	case cwin == 2:
		a.Command("masterstack", "splitv")
	case cwin == 3:
		a.Command("masterstack", "focus parent; splitv, focus child")
	}

	return nil
}
