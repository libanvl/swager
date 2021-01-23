package blocks

import (
	"errors"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc/event"
)

type workspace string

type InitSpawn struct {
	client        core.Client
	workspaceevts <-chan *event.WorkspaceChange
	spawns        map[workspace]string
}

func init() {
	var _ core.Block = (*InitSpawn)(nil)
	var _ core.Receiver = (*InitSpawn)(nil)
}

func (i *InitSpawn) Init(client core.Client, sub core.Sub, opts *core.Options) error {
	i.client = client
	i.workspaceevts = sub.Workspace()
	i.spawns = map[workspace]string{}
	return nil
}

func (i *InitSpawn) Configure(args []string) error {
	return nil
}

func (i *InitSpawn) Receive(args []string) error {
	if len(args) != 2 {
		return errors.New("requires two arguments")
	}

	if i.spawns == nil {
		i.spawns = map[workspace]string{workspace(args[0]): args[1]}
	} else {
		i.spawns[workspace(args[0])] = args[1]
	}

	return nil
}

func (i *InitSpawn) Close() {
}

func (i *InitSpawn) Run() {
	for evt := range i.workspaceevts {
		if evt.Change != event.InitWorkspace {
			continue
		}

		cmd, ok := i.spawns[workspace(evt.Current.Name)]
		if !ok {
			continue
		}

		i.client.Command(cmd)
	}
}
