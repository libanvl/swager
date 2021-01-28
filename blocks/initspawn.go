package blocks

import (
	"errors"
	"fmt"
	"sync"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type workspace string

/*
InitSpawn sends a command to sway when a workspace with a given name is created.

Registration

The block is registered with the name 'initspawn'

Configuration

InitSpawn has no init configuration:

  swayctrl init <tagname> initspawn

Send

Workspace name, command pairs are registered using the send command:

  swayctrl send <tagname> 2 "exec alacritty"
  swayctrl send <tagname> 3 "exec chromium --new-window"
*/
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

  if i.opts.Debug {
    i.opts.Log <- fmt.Sprintf("added spawn for workspace init: %s, '%s'", args[0], args[1])
  }

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
      if i.opts.Debug {
        i.opts.Log <- fmt.Sprintf("no spawn registered for workspace: '%s'", evt.Current.Name)
      }
			continue
		}

    if i.opts.Debug {
      i.opts.Log <- fmt.Sprintf("running spawn command: '%s'", cmd)
    }
		i.client.Command(cmd)
	}
}
