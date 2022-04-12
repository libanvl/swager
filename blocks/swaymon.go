package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type SwayMon struct {
	sub      core.Sub
	opts     *core.Options
	loglevel core.LogLevel
}

func init() {
	var _ core.BlockInitializer = (*SwayMon)(nil)
}

func (m *SwayMon) Init(client core.Client, sub core.Sub, opts *core.Options, args ...string) error {
	m.opts = opts
	m.sub = sub
	_, err := m.sub.ShutdownChanges(m.ShutdownChanged)
	if err != nil {
		return err
	}
	_, err = m.sub.WorkspaceChanges(m.WorkspaceChanged)
	if err != nil {
		return err
	}

	return nil
}

func (m *SwayMon) SetLogLevel(level core.LogLevel) {
	m.loglevel = level
}

func (m *SwayMon) ShutdownChanged(evt ipc.ShutdownChange) {
	if m.loglevel.Debug() {
		m.opts.Log.Printf("swaymon", "got shutdown event: %#v", evt.Change)
	}
	if evt.Change != ipc.ExitShutdown {
		return
	}

	m.opts.Log.Print("swaymon", "received shutdown event")
	m.opts.Server.RequestExit()
	return
}

func (m *SwayMon) WorkspaceChanged(evt ipc.WorkspaceChange) {
	if m.loglevel.Debug() {
		m.opts.Log.Printf("swaymon", "got workspace event: %#v", evt.Change)
	}
	if evt.Change != ipc.ReloadWorkspace {
		return
	}

	m.opts.Log.Print("swaymon", "received reload event")
	m.opts.Server.RequestExit()
	return
}
