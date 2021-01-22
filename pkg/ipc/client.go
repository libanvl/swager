package ipc

import (
	"encoding/json"
	"io"
	"net"
	"os"

	"github.com/libanvl/swager/pkg/ipc/reply"
)

type Client interface {
	io.ReadWriteCloser
	ClientRaw() ClientRaw
	RunCommand(cmd string) ([]reply.Command, error)
	GetTree() (*reply.Node, error)
	Subscribe(evts ...PayloadType) (*reply.Result, error)
	Version() (*reply.Version, error)
	Workspaces() ([]reply.Workspace, error)
}

func Connect() (Client, error) {
	addr, present := os.LookupEnv("SWAYSOCK")
	if !present {
		return nil, os.ErrNotExist
	}

	c, err := net.Dial("unix", addr)
	if err != nil {
		return nil, err
	}

	return &client{c}, nil
}

type client struct {
	io.ReadWriteCloser
}

func (c client) ClientRaw() ClientRaw {
	return c
}

func (c client) RunCommand(cmd string) ([]reply.Command, error) {
	res, err := c.ipccall(RunCommandMessage, []byte(cmd))
	if err != nil {
		return nil, err
	}

	var ss []reply.Command
	if err := json.Unmarshal(res, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (c client) GetTree() (*reply.Node, error) {
	res, err := c.ipccall(GetTreeMessage, nil)
	if err != nil {
		return nil, err
	}

	n := new(reply.Node)
	if err := json.Unmarshal(res, n); err != nil {
		return nil, err
	}

	return n, nil
}

func (c client) Subscribe(evts ...PayloadType) (*reply.Result, error) {
	pbytes, err := json.Marshal(eventNames(evts))
	if err != nil {
		return nil, err
	}

	res, err := c.ipccall(SubscribeMessage, pbytes)
	s := new(reply.Result)
	if err := json.Unmarshal(res, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (c client) Version() (*reply.Version, error) {
	res, err := c.ipccall(GetVersionMessage, nil)
	if err != nil {
		return nil, err
	}

	v := new(reply.Version)
	if err := json.Unmarshal(res, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (c client) Workspaces() ([]reply.Workspace, error) {
	res, err := c.ipccall(GetWorkspacesMessage, nil)
	if err != nil {
		return nil, err
	}
	var ws []reply.Workspace
	if err := json.Unmarshal(res, &ws); err != nil {
		return nil, err
	}

	return ws, nil
}
