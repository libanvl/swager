package core

import (
	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/ipc/event"
	"github.com/libanvl/swager/pkg/ipc/reply"
)

func init() {
	// core.Client must be a subset of ipc.Client
	var _ Client = (ipc.Client)(nil)
	// core.Subscription must be a subset of ipc.Subscription
	var _ Sub = (ipc.Subscription)(nil)
}

type Client interface {
	ClientRaw() ipc.ClientRaw
	RunCommand(cmd string) ([]reply.Command, error)
	GetTree() (*reply.Node, error)
	Subscribe(evts ...ipc.PayloadType) (*reply.Result, error)
	Version() (*reply.Version, error)
	Workspaces() ([]reply.Workspace, error)
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
