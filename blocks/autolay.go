package blocks

import (
	"sync"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/internal/core/node"
	"github.com/libanvl/swager/ipc"
	"github.com/libanvl/swager/stoker"
)

type LayoutEngine func(evt ipc.WindowChange, ws *ipc.Node) error

type Autolay struct {
	core.BasicBlock
	workspaces   map[string]LayoutEngine
	workspacesmx sync.Mutex
	eventmx      sync.Mutex
}

func init() {
	var _ core.BlockInitializer = (*Autolay)(nil)
}

func (a *Autolay) Init(client core.Client, sub core.Sub, opts *core.Options, log core.Logger, args ...string) error {
	a.Client = client
	a.Opts = opts
	a.Log = log
	a.workspaces = make(map[string]LayoutEngine)

	parser := stoker.NewParser[*Autolay](
		stoker.NewFlag("-autotiler", func(al *Autolay, tl stoker.TokenList) error {
			for _, ws := range tl {
				al.workspaces[ws] = al.autoTiler
				al.Log.Infof("Managing %v with autotiler", ws)
			}

			return nil
		}),

		stoker.NewFlag("-masterstack", func(al *Autolay, tl stoker.TokenList) error {
			for _, ws := range tl {
				al.workspaces[ws] = al.masterStack
				al.Log.Infof("Managing %v with masterstack", ws)
			}

			return nil
		}),
	)

	handler := parser.Parse(args...)

	a.Log.Defaultf("handler: %#v", handler)

	if err := handler.HandleAll(a); err != nil {
		a.Log.Debugf("handler error: %#v", err)
		return err
	}

	if _, err := sub.WindowChanges(a.WindowChanged); err != nil {
		return err
	}

	return nil
}

func (a *Autolay) SetLogLevel(level core.LogLevel) {
	a.LogLevel = level
}

func (a *Autolay) WindowChanged(evt ipc.WindowChange) {
	if evt.Container.Type == ipc.FloatingConNode {
		return
	}

	workspaces, err := a.Client.Workspaces()
	if err != nil {
		a.Log.Defaultf("(%v) Failed getting workspaces", evt.Container.ID)
		return
	}

	a.eventmx.Lock()
	defer a.eventmx.Unlock()

	focused := core.Focused(workspaces)
	if focused == nil {
		a.Log.Defaultf("(%v) Failed finding focused workspace", evt.Container.ID)
		return
	}

	eng, ok := a.workspaces[focused.Name]
	if !ok {
		a.Log.Debugf("(%v) Parent not managed: %v", evt.Container.ID, focused.Name)
		return
	}

	a.Log.Debugf("(%v) Using engine: %#v", evt.Container.ID, eng)

	root, err := a.Client.Tree()
	if err != nil {
		a.Log.Defaultf("(%v) Failed getting tree: %#v", evt.Container.ID, err)
	}

	workspace_node := node.First(
		root,
		node.MatchAnd(
			node.MatchType(ipc.WorkspaceNode),
			node.MatchName(focused.Name)))

	err = eng(evt, workspace_node)
	if err != nil {
		a.Log.Defaultf("(%v) Error executing step: %#v", err)
	}
}

func (a *Autolay) Command(engine_name string, cmd string) error {
	a.Log.Debugf("{%v} running command: %v", engine_name, cmd)

	res, err := a.Client.Command(cmd)
	if err != nil {
		a.Log.Defaultf("{%v} ipc error: %#v", engine_name, err)
		return err
	}

	if a.LogLevel.Debug() {
		for _, r := range res {
			a.Log.Debugf("{%v} Command result: %#v", engine_name, r)
		}
	}

	return nil
}

func IsTilingEligible(node *ipc.Node) bool {
	if node == nil {
		return false
	}

	if node.FullscreenMode != nil &&
		*node.FullscreenMode != ipc.NoneFullscreenMode {
		return false
	}

	if node.Type != ipc.ConNode {
		return false
	}

	return true
}
