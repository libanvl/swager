package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type SwayMon struct {
	sub       core.Sub
	opts      *core.Options
	shutdown  <-chan *ipc.ShutdownChange
	workspace <-chan *ipc.WorkspaceChange
}

func (m *SwayMon) Init(client core.Client, sub core.Sub, opts *core.Options) error {
	m.opts = opts
	m.sub = sub
	m.shutdown = m.sub.ShutdownChanges()
	m.workspace = m.sub.WorkspaceChanges()

	return nil
}

func (m *SwayMon) Configure(args []string) error {
	return nil
}

func (m *SwayMon) Run() {
	for {
		select {
		case wsc := <-m.workspace:
			if wsc.Change != ipc.ReloadWorkspace {
				continue
			}
			m.opts.Log.Default("swaymon", "received reload event")
			m.opts.Server.RequestExit()
      return
		case sdc := <-m.shutdown:
			if sdc.Change != ipc.ExitShutdown {
				continue
			}
			m.opts.Log.Default("swaymon", "received shutdown event")
			m.opts.Server.RequestExit()
			return
		}
	}
}

func (m *SwayMon) Close() {
}
