package ipc

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"os"
	"sync"
)

// Connect returns a Client connected to the UDS exported
// to the environment variable SWAYSOCK, using LittleEndian
// byte order.
func Connect() (*Client, error) {
	uds, present := os.LookupEnv("SWAYSOCK")
	if !present {
		return nil, os.ErrNotExist
	}

	return ConnectCustom(uds, binary.LittleEndian)
}

// ConnectCustom returns a Client connected to the UDS
// path specified by the uds parameter, with your choice of byte order.
func ConnectCustom(uds string, byteorder binary.ByteOrder) (*Client, error) {
	c, err := net.Dial("unix", uds)
	if err != nil {
		return nil, err
	}

	return &Client{c, byteorder, sync.Mutex{}}, nil
}

// Client is a sway-ipc compatible rpc client.
// Client is also an io.ReadWriteCloser
type Client struct {
	io.ReadWriteCloser
	yo    binary.ByteOrder
	ipcmx sync.Mutex
}

// Command implements the sway-ipc RUN_COMMAND message.
func (c *Client) Command(cmd string) ([]Command, error) {
	res, err := c.ipccall(runCommandMessage, []byte(cmd))
	if err != nil {
		return nil, err
	}

	var ss []Command
	if err := json.Unmarshal(res, &ss); err != nil {
		return nil, err
	}

	return ss, nil
}

func (c *Client) CommandRaw(cmd string) (string, error) {
	return c.ipccallraw(runCommandMessage, []byte(cmd))
}

// Workspaces implements the sway-ipc GET_WORKSPACES message.
func (c *Client) Workspaces() ([]Workspace, error) {
	res, err := c.ipccall(getWorkspacesMessage, nil)
	if err != nil {
		return nil, err
	}
	var ws []Workspace
	if err := json.Unmarshal(res, &ws); err != nil {
		return nil, err
	}

	return ws, nil
}

func (c *Client) WorkspacesRaw() (string, error) {
	return c.ipccallraw(getWorkspacesMessage, nil)
}

// Subscribe implements the sway-ipc SUBSCRIBE message.
func (c *Client) Subscribe(evts ...EventPayloadType) (*Result, error) {
	pbytes, err := json.Marshal(eventNames(evts))
	if err != nil {
		return nil, err
	}

	res, err := c.ipccall(subscribeMessage, pbytes)
	s := new(Result)
	if err := json.Unmarshal(res, s); err != nil {
		return nil, err
	}

	return s, nil
}

// Outputs implements the sway-ipc GET_OUTPUTS message.
func (c *Client) Outputs() ([]Output, error) {
	res, err := c.ipccall(getOutputsMessage, nil)
	if err != nil {
		return nil, err
	}

	var os []Output
	if err := json.Unmarshal(res, &os); err != nil {
		return nil, err
	}

	return os, nil
}

func (c *Client) OutputsRaw() (string, error) {
	return c.ipccallraw(getOutputsMessage, nil)
}

// Tree implements the sway-ipc GET_TREE message.
// Returns a *Node representing the root of the tree.
func (c *Client) Tree() (*Node, error) {
	res, err := c.ipccall(getTreeMessage, nil)
	if err != nil {
		return nil, err
	}

	n := new(Node)
	if err := json.Unmarshal(res, n); err != nil {
		return nil, err
	}

	return n, nil
}

func (c *Client) TreeRaw() (string, error) {
	return c.ipccallraw(getTreeMessage, nil)
}

// Version implements the sway-ipc GET_VERSION message.
func (c *Client) Version() (*Version, error) {
	res, err := c.ipccall(getVersionMessage, nil)
	if err != nil {
		return nil, err
	}

	v := new(Version)
	if err := json.Unmarshal(res, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Client) VersionRaw() (string, error) {
	return c.ipccallraw(getVersionMessage, nil)
}
