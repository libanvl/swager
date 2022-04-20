package blocks

import (
	"sync"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/internal/core/node"
	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/stoker"
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

	parser := stoker.NewParser(
		stoker.NewFlag("autotiler", func(_ any, tl stoker.TokenList) error {
			for _, ws := range tl {
				a.workspaces[ws] = a.autoTiler
				a.Log.Infof("Managing %v with autotiler", ws)
			}

			return nil
		}),

		stoker.NewFlag("masterstack", func(_ any, tl stoker.TokenList) error {
			for _, ws := range tl {
				a.workspaces[ws] = a.masterStack
				a.Log.Infof("Managing %v with masterstack", ws)
			}

			return nil
		}),
	)

	handler := parser.Parse(args...)

	a.Log.Defaultf("executor: %#v", handler)

	if err := handler.HandleAll(nil); err != nil {
		a.Log.Debugf("Executor error: %#v", err)
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
		a.Log.Defaultf("(%v) Failed getting tree: %v", evt.Container.ID, err)
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
		a.Log.Defaultf("{%v} ipc error: %v", engine_name, err)
		return err
	}

	if a.LogLevel.Debug() {
		for _, r := range res {
			a.Log.Debugf("{%v} Command result: %v", engine_name, r)
		}
	}

	return nil
}
