package blocks

import (
	"errors"
	"sync"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type workspace string

type InitSpawn struct {
	client        core.Client
	opts          *core.Options
	workspaceevts <-chan *ipc.WorkspaceChange
	spawns        map[workspace]string
	spawnsmx      sync.Mutex
}

func init() {
	var _ core.Block = (*InitSpawn)(nil)
	var _ core.Receiver = (*InitSpawn)(nil)
}

func (i *InitSpawn) Init(client core.Client, sub core.Sub, opts *core.Options) error {
	i.client = client
	i.opts = opts
	i.workspaceevts = sub.WorkspaceChanges()
	i.spawns = map[workspace]string{}
	i.spawnsmx = sync.Mutex{}
	return nil
}

func (i *InitSpawn) Configure(args []string) error {
	return nil
}

func (i *InitSpawn) Receive(args []string) error {
	if len(args) != 2 {
		return errors.New("requires two arguments: <workspace> <command>")
	}

	i.spawnsmx.Lock()
	if i.spawns == nil {
		i.spawns = map[workspace]string{workspace(args[0]): args[1]}
	} else {
		i.spawns[workspace(args[0])] = args[1]
	}
	i.spawnsmx.Unlock()

	i.opts.Log.Infof("initspawn", "added spawn for workspace init: %s, '%s'", args[0], args[1])

	return nil
}

func (i *InitSpawn) Close() {
}

func (i *InitSpawn) Run() {
	for evt := range i.workspaceevts {
		if evt.Change != ipc.InitWorkspace {
			continue
		}

		i.spawnsmx.Lock()
		cmd, ok := i.spawns[workspace(evt.Current.Name)]
		i.spawnsmx.Unlock()
		if !ok {
			i.opts.Log.Debugf("initspawn", "no spawn registered for workspace: '%s'", evt.Current.Name)
			continue
		}

		i.opts.Log.Debugf("initspawn", "running spawn command: '%s'", cmd)
		i.client.Command(cmd)
	}
}
