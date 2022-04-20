package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/internal/core/node"
	"github.com/libanvl/swager/pkg/ipc"
)

func (a *Autolay) autoTiler(evt ipc.WindowChange, ws *ipc.Node) error {
	if core.Deny(evt.Change, ipc.TitleWindow, ipc.NewWindow) {
		a.Log.Debugf("{autotiler} Denying change type: %v", evt.Change)
		return nil
	}

	focused := node.First(ws, node.MatchAnd(node.IsFocused, node.IsLeaf))

	if !IsTilingEligible(focused) {
		a.Log.Debug("{autotiler} focused is not eligible for tiling")
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
	if core.Deny(evt.Change, ipc.TitleWindow, ipc.FocusWindow) {
		a.Log.Debugf("{masterstack} Denying change type: %v", evt.Change)
		return nil
	}

	focused := node.First(ws, node.MatchAnd(node.IsFocused, node.IsLeaf))

	if !IsTilingEligible(focused) {
		a.Log.Debug("{masterstack} focused is not eligible for tiling")
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
