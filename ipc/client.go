package ipc

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"os"
	"sync"
)

// Client is a sway-ipc compatible rpc client.
// Client is also an io.ReadWriteCloser.
type Client struct {
	io.ReadWriteCloser
	yo    binary.ByteOrder
	ipcmx sync.Mutex
}

func NewClient(conn io.ReadWriteCloser, yo binary.ByteOrder) *Client {
	return &Client{conn, yo, sync.Mutex{}}
}

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
func ConnectCustom(uds string, yo binary.ByteOrder) (*Client, error) {
	c, err := net.Dial("unix", uds)
	if err != nil {
		return nil, err
	}

	return NewClient(c, yo), nil
}

// Command implements the sway-ipc RUN_COMMAND message.
func (c *Client) Command(cmd string) ([]Command, error) {
	return callgetarr[Command](c, RunCommandMessage, []byte(cmd))
}

// CommandRaw implements the sway-ipc RUN_COMMAND message
// and returns a json string.
func (c *Client) CommandRaw(cmd string) (string, error) {
	return c.ipccallraw(RunCommandMessage, []byte(cmd))
}

// Workspaces implements the sway-ipc GET_WORKSPACES message.
func (c *Client) Workspaces() ([]Workspace, error) {
	return callgetarr[Workspace](c, GetWorkspacesMessage, nil)
}

// Workspaces implements the sway-ipc GET_WORKSPACES message
// and returns a json string.
func (c *Client) WorkspacesRaw() (string, error) {
	return c.ipccallraw(GetWorkspacesMessage, nil)
}

// Subscribe implements the sway-ipc SUBSCRIBE message.
func (c *Client) Subscribe(evts ...EventPayloadType) (*Result, error) {
	pbytes, err := json.Marshal(eventNames(evts))
	if err != nil {
		// panic here because this shouldn't be possible
		panic(err)
	}

	return callgetptr[Result](c, SubscribeMessage, pbytes)
}

// Outputs implements the sway-ipc GET_OUTPUTS message.
func (c *Client) Outputs() ([]Output, error) {
	return callgetarr[Output](c, GetOutputsMessage, nil)
}

// Outputs implements the sway-ipc GET_OUTPUTS message
// and returns a json string.
func (c *Client) OutputsRaw() (string, error) {
	return c.ipccallraw(GetOutputsMessage, nil)
}

// Tree implements the sway-ipc GET_TREE message.
// Returns a *Node representing the root of the tree.
func (c *Client) Tree() (*Node, error) {
	return callgetptr[Node](c, GetTreeMessage, nil)
}

// Tree implements the sway-ipc GET_TREE message
// and returns a json string.
func (c *Client) TreeRaw() (string, error) {
	return c.ipccallraw(GetTreeMessage, nil)
}

// Marks implements the sway-ipc GET_MARKS message.
func (c *Client) Marks() ([]string, error) {
	return callgetarr[string](c, GetMarksMessage, nil)
}

// Marks implements the sway-ipc GET_MARKS message
// and returns a json string.
func (c *Client) MarksRaw() (string, error) {
	return c.ipccallraw(GetMarksMessage, nil)
}

// Version implements the sway-ipc GET_VERSION message.
func (c *Client) Version() (*Version, error) {
	return callgetptr[Version](c, GetVersionMessage, nil)
}

// Version implements the sway-ipc GET_VERSION message
// and returns a json string.
func (c *Client) VersionRaw() (string, error) {
	return c.ipccallraw(GetVersionMessage, nil)
}

// BindingModes implements the sway-ipc GET_BINDING_MODES message.
func (c *Client) BindingModes() ([]string, error) {
	return callgetarr[string](c, GetBindingModesMessage, nil)
}

// BindingModes implements the sway-ipc GET_BINDING_MODES message
// and returns a json string.
func (c *Client) BindingModesRaw() (string, error) {
	return c.ipccallraw(GetBindingModesMessage, nil)
}

// Tick implements the sway-ipc SEND_TICK message.
func (c *Client) Tick(payload string) (*Result, error) {
	return callgetptr[Result](c, SendTickMessage, []byte(payload))
}

// Tick implements the sway-ipc SEND_TICK message
// and returns a json string.
func (c *Client) TickRaw(payload string) (string, error) {
	return c.ipccallraw(SendTickMessage, []byte(payload))
}

// BindingState implements the sway-ipc GET_BINDING_STATE message.
func (c *Client) BindingState() (*BindingState, error) {
	return callgetptr[BindingState](c, GetBindingStateMessage, nil)
}

// BindingState implements the sway-ipc GET_BINDING_STATE message
// and returns a json string.
func (c *Client) BindingStateRaw() (string, error) {
	return c.ipccallraw(GetBindingStateMessage, nil)
}
