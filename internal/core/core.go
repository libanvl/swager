package core

import (
	"fmt"

	"github.com/libanvl/swager/pkg/ipc"
)

func init() {
	// core.Client must be a subset of ipc.Client
	// core.Subscription must be a subset of ipc.Subscription
	var _ Client = (*ipc.Client)(nil)
	var _ Sub = (*ipc.Subscription)(nil)
  var _ BlockLogMessage = defaultLogMessage{}
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

type BlockLogMessage interface {
  String() string
}

type defaultLogMessage struct {
  prefix string
  message string
}

func (dlm defaultLogMessage) String() string {
  return fmt.Sprintf("[%s] %s", dlm.prefix, dlm.message)
}

type BlockLogChannel chan<- BlockLogMessage

func (blc BlockLogChannel) Message(prefix string, msg string){
  dlm := defaultLogMessage{prefix, msg}
  blc<- dlm
}

func (blc BlockLogChannel) Messagef(prefix string, format string, args...interface{}) {
  blc.Message(prefix, fmt.Sprintf(format, args...))
}

// Options are shared options for use by core.Block instances.
// Debug indicates that debug logging was requested when starting the daemon.
// Use the Log channel to send log data back to the daemon.
type Options struct {
	Debug bool
	Log   BlockLogChannel
}
