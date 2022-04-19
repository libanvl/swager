package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/internal/core/node"
	"github.com/libanvl/swager/pkg/ipc"
)

func (a *Autolay) autoTiler(evt ipc.WindowChange, ws *ipc.Node) error {
	if evt.Container.Type == ipc.FloatingConNode {
		return nil
	}

	if core.Deny(evt.Change, ipc.TitleWindow) {
		return nil
	}

	cwin := node.Count(ws, node.IsLeaf)
	is_even := (cwin % 2) == 0

	if is_even {
		a.Command("autotiler", "splitv")
	} else {
		a.Command("autotiler", "splith")
	}
	return nil
}

func (a *Autolay) masterStack(evt ipc.WindowChange, ws *ipc.Node) error {
	if evt.Container.Type == ipc.FloatingConNode {
		return nil
	}

	if core.Deny(evt.Change, ipc.TitleWindow, ipc.FocusWindow) {
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
