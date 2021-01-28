package core

import (
	"github.com/libanvl/swager/pkg/ipc"
)

func init() {
	// core.Client must be a subset of ipc.Client
	// core.Subscription must be a subset of ipc.Subscription
	var _ Client = (*ipc.Client)(nil)
	var _ Sub = (*ipc.Subscription)(nil)
}

// Client exports a limited set of methods for use by core.Block instances.
type Client interface {
	Command(cmd string) ([]ipc.Command, error)
	CommandRaw(cmd string) (string, error)
  Workspaces() ([]ipc.Workspace, error)
  WorkspacesRaw() (string, error)
	Tree() (*ipc.Node, error)
  TreeRaw() (string, error)
	Version() (*ipc.Version, error)
  VersionRaw() (string, error)
}

// Sub exports a limited set of methods for use by core.Block instances.
type Sub interface {
	WindowChanges() <-chan *ipc.WindowChange
	WorkspaceChanges() <-chan *ipc.WorkspaceChange
	ShutdownChanges() <-chan *ipc.ShutdownChange
}

// Options are shared options for use by core.Block instances.
// Debug indicates that debug logging was requested when starting the daemon.
// Use the Log channel to send log data back to the daemon.
type Options struct {
	Debug bool
	Log   chan<- string
}
