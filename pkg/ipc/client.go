package ipc

import (
	"encoding/json"
	"io"
	"net"
	"os"
)

type Client interface {
	io.ReadWriteCloser
	Command(cmd string) ([]Command, error)
	Workspaces() ([]Workspace, error)
	Subscribe(evts ...PayloadType) (*Result, error)
	Outputs() ([]Output, error)
  Tree() (*Node, error)
	Version() (*Version, error)
	ClientRaw() ClientRaw
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

func (c client) Command(cmd string) ([]Command, error) {
	res, err := c.ipccall(RunCommandMessage, []byte(cmd))
	if err != nil {
		return nil, err
	}

	var ss []Command
	if err := json.Unmarshal(res, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (c client) Workspaces() ([]Workspace, error) {
	res, err := c.ipccall(GetWorkspacesMessage, nil)
	if err != nil {
		return nil, err
	}
	var ws []Workspace
	if err := json.Unmarshal(res, &ws); err != nil {
		return nil, err
	}

	return ws, nil
}

func (c client) Subscribe(evts ...PayloadType) (*Result, error) {
	pbytes, err := json.Marshal(eventNames(evts))
	if err != nil {
		return nil, err
	}

	res, err := c.ipccall(SubscribeMessage, pbytes)
	s := new(Result)
	if err := json.Unmarshal(res, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (c client) Outputs() ([]Output, error) {
  res, err := c.ipccall(GetOutputsMessage, nil)
  if err != nil {
    return nil, err
  }

  var os []Output
  if err := json.Unmarshal(res, &os); err != nil {
    return nil, err
  }

  return os, nil
}

func (c client) Tree() (*Node, error) {
	res, err := c.ipccall(GetTreeMessage, nil)
	if err != nil {
		return nil, err
	}

	n := new(Node)
	if err := json.Unmarshal(res, n); err != nil {
		return nil, err
	}

	return n, nil
}

func (c client) Version() (*Version, error) {
	res, err := c.ipccall(GetVersionMessage, nil)
	if err != nil {
		return nil, err
	}

	v := new(Version)
	if err := json.Unmarshal(res, v); err != nil {
		return nil, err
	}
	return v, nil
}
