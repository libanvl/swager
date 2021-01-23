package core

import (
	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/ipc/event"
)

func init() {
	// core.Client must be a subset of ipc.Client
	// core.Subscription must be a subset of ipc.Subscription
	var _ Client = (ipc.Client)(nil)
	var _ Sub = (event.Subscription)(nil)
}

type Client interface {
	ClientRaw() ipc.ClientRaw
	Command(cmd string) ([]ipc.Command, error)
	Workspaces() ([]ipc.Workspace, error)
	Subscribe(evts ...ipc.PayloadType) (*ipc.Result, error)
	Tree() (*ipc.Node, error)
	Version() (*ipc.Version, error)
}

type Sub interface {
	Window() <-chan *event.WindowChange
	Workspace() <-chan *event.WorkspaceChange
	Shutdown() <-chan *event.ShutdownChange
}

type Options struct {
	Debug bool
	Log   chan<- string
}
