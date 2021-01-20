package ipc

import (
  "encoding/json"
  "net"
  "os"

  "github.com/libanvl/swager/pkg/ipc/reply"
)

type client struct {
  conn net.Conn
}

type Client interface {
  Conn() net.Conn
  Close() error
  ClientRaw() ClientRaw
  Subscribe(evts ...string) (*reply.Success, error)
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

  return client{c}, nil
}

func (c client) Conn() net.Conn {
  return c.conn
}

func (c client) Close() error {
  return c.conn.Close()
}

func (c client) ClientRaw() ClientRaw {
  return c
}

func (c client) Subscribe(evts ...string) (*reply.Success, error) {
  pbytes, err := json.Marshal(evts)
  if err != nil {
    return nil, err
  }

  res, err := c.ipccall(SubscribePayload, pbytes)
  s := new(reply.Success)
  if err := json.Unmarshal(res, s); err != nil {
    return nil, err
  }

  return s, nil
}

func (c client) Version() (*reply.Version, error) {
  res, err := c.ipccall(GetVersionPayload, nil)
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
  res, err := c.ipccall(GetWorkspacesPayload, nil)
  if err != nil {
    return nil, err
  }
  var ws []reply.Workspace
  if err := json.Unmarshal(res, &ws); err != nil {
    return nil, err
  }

  return ws, nil
}

