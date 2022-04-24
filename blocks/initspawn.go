package blocks

import (
	"errors"
	"sync"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/ipc"
)

type workspace string

type InitSpawn struct {
	client        core.Client
	opts          *core.Options
	workspaceevts ipc.Cookie
	spawns        map[workspace]string
	spawnsmx      sync.Mutex
	log           core.Logger
	loglevel      core.LogLevel
}

func init() {
	var _ core.BlockInitializer = (*InitSpawn)(nil)
	var _ core.Receiver = (*InitSpawn)(nil)
}

func (i *InitSpawn) Init(client core.Client, sub core.Sub, opts *core.Options, log core.Logger, args ...string) error {
	i.client = client
	i.opts = opts
	i.log = log
	i.spawns = map[workspace]string{}
	i.spawnsmx = sync.Mutex{}
	cookie, err := sub.WorkspaceChanges(i.WorkspaceChanged)
	if err != nil {
		return err
	}

	i.workspaceevts = cookie
	return nil
}

func (i *InitSpawn) SetLogLevel(level core.LogLevel) {
	i.loglevel = level
}

func (i *InitSpawn) WorkspaceChanged(evt ipc.WorkspaceChange) {
	i.log.Debugf("got workspace event: %#v, %s", evt.Change, evt.Current.Name)
	if evt.Change != ipc.InitWorkspace {
		return
	}

	i.spawnsmx.Lock()
	cmd, ok := i.spawns[workspace(evt.Current.Name)]
	i.spawnsmx.Unlock()
	if !ok {
		i.log.Defaultf("no spawn registered for workspace: '%s'", evt.Current.Name)
		return
	}

	if i.loglevel.Debug() {
		i.log.Debugf("nodes count: %d", len(evt.Current.Nodes))
	}

	if len(evt.Current.Nodes) < 1 {

		i.log.Defaultf("running spawn command: '%s'", cmd)

		res, err := i.client.Command(cmd)
		if err != nil {
			i.log.Default(err.Error())
		}

		for _, r := range res {
			if !r.Success {
				i.log.Default(r.Error)
			}
		}
	}
}

func (i *InitSpawn) Receive(args []string) error {
	if len(args) != 2 {
		return errors.New("requires two arguments: <workspace> <command>")
	}

	i.spawnsmx.Lock()
	defer i.spawnsmx.Unlock()

	if i.spawns == nil {
		i.spawns = map[workspace]string{workspace(args[0]): args[1]}
	} else {
		i.spawns[workspace(args[0])] = args[1]
	}

	i.log.Infof("added spawn for workspace init: %s, '%s'", args[0], args[1])
	return nil
}
