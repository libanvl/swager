package core

import (
	"flag"

	"github.com/libanvl/swager/pkg/ipc"
)

func init() {
	// core.Client must be a subset of ipc.Client
	// core.Subscription must be a subset of ipc.Subscription
	var _ Client = (*ipc.Client)(nil)
	var _ Sub = (*ipc.Subscription)(nil)
	var _ LogMessage = PrefixLogMessage{}
	var _ flag.Value = (*LogLevel)(nil)
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
	WorkspaceChanges() <-chan *ipc.WorkspaceChange
	WindowChanges() <-chan *ipc.WindowChange
	BindingChanges() <-chan *ipc.BindingChange
	ShutdownChanges() <-chan *ipc.ShutdownChange
	Ticks() <-chan *ipc.Tick
}

type ServerControlRequest int8

const (
	ReloadRequest ServerControlRequest = 0
	ExitRequest   ServerControlRequest = 1
)

type ServerControlChannel chan<- ServerControlRequest

func (scc ServerControlChannel) RequestReload() {
	go func() {
		scc <- ReloadRequest
	}()
}

func (scc ServerControlChannel) RequestExit() {
	go func() {
		scc <- ExitRequest
	}()
}

// Options are shared options for use by core.Block instances.
// Debug indicates that debug logging was requested when starting the daemon.
// Use the Log channel to send log data back to the daemon.
type Options struct {
	Log    LogChannel
	Server ServerControlChannel
}
