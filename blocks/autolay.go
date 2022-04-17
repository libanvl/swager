package blocks

import (
	"sync"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/internal/core/node"
	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/stoker"
)

type LayoutEngine func(wct ipc.WindowChangeType, ws *ipc.Node) error

type Autolay struct {
	core.BasicBlock
	workspaces   map[string]LayoutEngine
	workspacesmx sync.Mutex
	eventmx      sync.Mutex
}

func init() {
	var _ core.BlockInitializer = (*Autolay)(nil)
}

func (a *Autolay) Init(client core.Client, sub core.Sub, opts *core.Options, args ...string) error {
	a.Client = client
	a.Opts = opts
	a.workspaces = make(map[string]LayoutEngine)

	parser := stoker.Parser(
		stoker.Def("--autotiler"),
		stoker.Def("--masterstack"),
	)

	tokenmap := parser.Parse(args...)

	tokenmap.ProcessSet("--autotiler", func(ts stoker.TokenSet) error {
		for _, tokenlist := range ts {
			for _, ws := range tokenlist {
				a.workspaces[ws] = a.autoTiler
			}
		}
		return nil
	})

	tokenmap.ProcessSet("--masterstack", func(ts stoker.TokenSet) error {
		for _, tokenlist := range ts {
			for _, ws := range tokenlist {
				a.workspaces[ws] = a.masterStack
			}
		}
		return nil
	})

	sub.WindowChanges(a.WindowChanged)

	return nil
}

func (a *Autolay) SetLogLevel(level core.LogLevel) {
	a.LogLevel = level
}

func (a *Autolay) WindowChanged(evt ipc.WindowChange) {
	workspaces, err := a.Client.Workspaces()
	if err != nil {
		a.Opts.Log.Printf("autolay", "(%v) Failed getting workspaces", evt.Container.ID)
		return
	}

	a.eventmx.Lock()
	defer a.eventmx.Unlock()

	focused := core.Focused(workspaces)
	if focused == nil {
		a.Opts.Log.Printf("autolay", "(%v) Failed finding focused workspace", evt.Container.ID)
	}

	eng, ok := a.workspaces[focused.Name]
	if !ok {
		if a.LogLevel.Debug() {
			a.Opts.Log.Printf("autolay", "(%v) Parent not managed: %v", evt.Container.ID, focused.Name)
		}
		return
	}

	if a.LogLevel.Debug() {
		a.Opts.Log.Printf("autolay", "(%v) Using engine: %#v", evt.Container.ID, eng)
	}

	root, err := a.Client.Tree()
	if err != nil {
		a.Opts.Log.Printf("autolay", "(%v) Failed getting tree: %v", evt.Container.ID, err)
	}

	workspace_node := node.First(
		root,
		node.MatchAnd(
			node.MatchType(ipc.WorkspaceNode),
			node.MatchName(focused.Name)))

	err = eng(evt.Change, workspace_node)
	if err != nil {
		a.Opts.Log.Printf("autolay", "(%v) Error executing step: %#v", err)
	}
}

func (a *Autolay) Command(engine_name string, cmd string) error {
	if a.LogLevel.Debug() {
		a.Opts.Log.Printf("autolay", "{%v} running command: %v", engine_name, cmd)
	}

	res, err := a.Client.Command(cmd)
	if err != nil {
		a.Opts.Log.Printf("autolay", "{%v} ipc error: %v", engine_name, err)
		return err
	}

	if a.LogLevel.Debug() {
		for _, r := range res {
			a.Opts.Log.Printf("autolay", "{%v} Command result: %v", engine_name, r)
		}
	}

	return nil
}
