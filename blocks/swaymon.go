package blocks

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/ipc"
)

type SwayMon struct {
	sub  core.Sub
	opts *core.Options
	log  core.Logger
}

func init() {
	var _ core.BlockInitializer = (*SwayMon)(nil)
}

func (m *SwayMon) Init(client core.Client, sub core.Sub, opts *core.Options, log core.Logger, args ...string) error {
	m.opts = opts
	m.sub = sub
	m.log = log

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
}

func (m *SwayMon) ShutdownChanged(evt ipc.ShutdownChange) {
	m.log.Debugf("got shutdown event: %#v", evt.Change)
	if evt.Change != ipc.ExitShutdown {
		return
	}

	m.log.Info("received shutdown event")
	m.opts.Server.RequestExit()
	return
}

func (m *SwayMon) WorkspaceChanged(evt ipc.WorkspaceChange) {
	m.log.Debugf("got workspace event: %#v", evt.Change)
	if evt.Change != ipc.ReloadWorkspace {
		return
	}

	m.log.Info("received reload event")
	m.opts.Server.RequestExit()
	return
}
