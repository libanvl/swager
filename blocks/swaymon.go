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
	loglevel  core.LogLevel
}

func init() {
	var _ core.Block = (*SwayMon)(nil)
}

func (m *SwayMon) Init(client core.Client, sub core.Sub, opts *core.Options, args ...string) error {
	m.opts = opts
	m.sub = sub
	m.shutdown = m.sub.ShutdownChanges()
	m.workspace = m.sub.WorkspaceChanges()

	return nil
}

func (m *SwayMon) SetLogLevel(level core.LogLevel) {
	m.loglevel = level
}

func (m *SwayMon) Run() {
	for {
		select {
		case wsc := <-m.workspace:
			if m.loglevel.Debug() {
				m.opts.Log.Printf("swaymon", "got workspace event: %#v", wsc.Change)
			}
			if wsc.Change != ipc.ReloadWorkspace {
				continue
			}
			m.opts.Log.Print("swaymon", "received reload event")
			m.opts.Server.RequestExit()
			return
		case sdc := <-m.shutdown:
			if m.loglevel.Debug() {
				m.opts.Log.Printf("swaymon", "got shutdown event: %#v", sdc.Change)
			}
			if sdc.Change != ipc.ExitShutdown {
				continue
			}

			m.opts.Log.Print("swaymon", "received shutdown event")
			m.opts.Server.RequestExit()
			return
		}
	}
}
