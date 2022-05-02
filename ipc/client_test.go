package ipc_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/libanvl/swager/ipc"
	"github.com/libanvl/swager/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	tmpsocket := t.TempDir() + "/uds"
	net.Listen("unix", tmpsocket)
	t.Setenv("SWAYSOCK", tmpsocket)
	client, err := ipc.Connect()

	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestConnectNoSwaysock(t *testing.T) {
	t.Setenv("SWAYSOCK", "")
	_, err := ipc.Connect()
	assert.NotNil(t, err)

	os.Unsetenv("SWAYSOCK")
	_, err = ipc.Connect()
	assert.NotNil(t, err)
}

func TestPayloadWrites(t *testing.T) {
	tests := map[string]struct {
		action  func(*ipc.Client, string)
		payload string
	}{
		"Command":               {func(c *ipc.Client, s string) { c.Command(s) }, "testpayload"},
		"CommandWithSemiColons": {func(c *ipc.Client, s string) { c.Command(s) }, "testpayload; testpayload2"},
		"CommandWithBrackets":   {func(c *ipc.Client, s string) { c.Command(s) }, "[app_id=testpayload] testpayload2"},
		"CommandRaw":            {func(c *ipc.Client, s string) { c.CommandRaw(s) }, "[app_id=testpayload] testpayload2"},
		"Tick":                  {func(c *ipc.Client, s string) { c.Tick(s) }, "testpayload"},
		"TickRaw":               {func(c *ipc.Client, s string) { c.TickRaw(s) }, "testpayload"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := require.New(t)
			conn := test.NewMockConnection(t)
			client := ipc.NewClient(conn, binary.LittleEndian)
			tc.action(client, tc.payload)
			conn.AssertNumberOfCalls("Write", 2)
			conn.AssertCalled("Read")
			conn.AssertNotCalled("Close")

			var got_header ipc.Header
			r := bytes.NewReader(conn.WriteAt(0))
			err := binary.Read(r, binary.LittleEndian, &got_header)
			assert.Nil(err)
			assert.True(ipc.ValidMagic(got_header.Magic), "Magic should be valid")
			assert.Equal(uint32(len(tc.payload)), got_header.PayloadLength)
			assert.Equal(tc.payload, string(conn.WriteAt(1)))
		})
	}
}

func TestPayloadWritesWithWriteError(t *testing.T) {
	tests := map[string]struct {
		action  func(*ipc.Client, string) (any, error)
		payload string
	}{
		"Command":    {func(c *ipc.Client, s string) (any, error) { return c.Command(s) }, "commandpayload"},
		"CommandRaw": {func(c *ipc.Client, s string) (any, error) { return c.CommandRaw(s) }, "commandrawpayload"},
		"Tick":       {func(c *ipc.Client, s string) (any, error) { return c.Tick(s) }, "tickpayload"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := require.New(t)
			conn := test.NewMockConnection(t)
			client := ipc.NewClient(conn, binary.LittleEndian)

			conn.SetNextWriteResult(1, net.ErrClosed)
			_, err := tc.action(client, tc.payload)
			assert.NotNil(err)

			conn.AssertCalled("Write")
			conn.AssertNotCalled("Read")
			conn.AssertNotCalled("Close")
		})
	}
}
func TestNoPayloadWrites(t *testing.T) {
	tests := map[string]struct {
		action func(*ipc.Client)
	}{
		"Workspaces":      {func(c *ipc.Client) { c.Workspaces() }},
		"WorkspacesRaw":   {func(c *ipc.Client) { c.WorkspacesRaw() }},
		"Outputs":         {func(c *ipc.Client) { c.Outputs() }},
		"OutputsRaw":      {func(c *ipc.Client) { c.OutputsRaw() }},
		"Tree":            {func(c *ipc.Client) { c.Tree() }},
		"TreeRaw":         {func(c *ipc.Client) { c.TreeRaw() }},
		"Marks":           {func(c *ipc.Client) { c.Marks() }},
		"MarksRaw":        {func(c *ipc.Client) { c.MarksRaw() }},
		"Version":         {func(c *ipc.Client) { c.Version() }},
		"VersionRaw":      {func(c *ipc.Client) { c.VersionRaw() }},
		"BindingModes":    {func(c *ipc.Client) { c.BindingModes() }},
		"BindingModesRaw": {func(c *ipc.Client) { c.BindingModesRaw() }},
		"BindingState":    {func(c *ipc.Client) { c.BindingState() }},
		"BindingStateRaw": {func(c *ipc.Client) { c.BindingStateRaw() }},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			conn := test.NewMockConnection(t)
			client := ipc.NewClient(conn, binary.LittleEndian)
			tc.action(client)
			conn.AssertNumberOfCalls("Write", 1)
			conn.AssertCalled("Read")
			conn.AssertNotCalled("Close")

			var got_header ipc.Header
			r := bytes.NewReader(conn.WriteAt(0))
			err := binary.Read(r, binary.LittleEndian, &got_header)
			assert.Nil(err)
			assert.True(ipc.ValidMagic(got_header.Magic), "Magic should be valid")
			assert.Equal(uint32(0), got_header.PayloadLength)
		})
	}
}

func TestNoPayloadReads(t *testing.T) {
	type action func(*ipc.Client) (any, error)
	tests := map[string]struct {
		payload  ipc.PayloadType
		action   action
		expected any
	}{
		"Workspaces": {
			ipc.GetWorkspacesMessage,
			func(c *ipc.Client) (any, error) { return c.Workspaces() },
			[]ipc.Workspace{
				{
					Num:     1,
					Name:    "TestWorkspaceOne",
					Visible: true,
					Focused: true,
					Urgent:  false,
					Rect:    ipc.Rect{X: 500, Y: 500, Width: 500, Height: 500},
					Output:  "TestOutput",
				},
				{
					Num:     2,
					Name:    "TestWorkspaceTwo",
					Visible: false,
					Focused: false,
					Urgent:  true,
					Rect:    ipc.Rect{X: 1000, Y: 1000, Width: 1000, Height: 1000},
					Output:  "TestOutput2",
				},
			},
		},
		"Outputs": {
			ipc.GetOutputsMessage,
			func(c *ipc.Client) (any, error) { return c.Outputs() },
			[]ipc.Output{
				{
					Name:             "TestOutput",
					Make:             "TestMake",
					Model:            "TestModel",
					Serial:           "TestSerial",
					Active:           true,
					Dpms:             true,
					Primary:          true,
					Scale:            5.0,
					SubpixelHinting:  "",
					Transform:        "",
					CurrentWorkspace: "TestWorkspace",
					CurrentMode:      ipc.Mode{100, 100, 60},
					Rect:             ipc.Rect{X: 500, Y: 500, Width: 500, Height: 500},
				},
			},
		},
		"Version": {
			ipc.GetVersionMessage,
			func(c *ipc.Client) (any, error) { return c.Version() },
			&ipc.Version{
				Major:                100,
				Minor:                200,
				Patch:                3000,
				HumanReadable:        "TestHumanReadable",
				LoadedConfigFileName: "TestConfigFileName",
			},
		},
	}

	for name, tc := range tests {
		for _, yo := range []binary.ByteOrder{binary.LittleEndian, binary.BigEndian} {
			t.Run(name+"_"+yo.String(), func(t *testing.T) {
				conn := test.NewMockConnection(t)
				client := ipc.NewClient(conn, yo)
				expected_json, err := json.Marshal(tc.expected)
				require.Nil(t, err)
				require.NotNil(t, expected_json)
				conn.PushPayloadForRead(uint32(tc.payload), expected_json, yo)
				actual, err := tc.action(client)
				assert.Nil(t, err)
				assert.Equal(t, tc.expected, actual)
				conn.AssertNumberOfCalls("Read", 2)
				conn.AssertCalled("Write")
				conn.AssertNotCalled("Close")
			})
		}
	}
}

func TestSubscribe(t *testing.T) {
	conn := test.NewMockConnection(t)
	client := ipc.NewClient(conn, binary.LittleEndian)
	result, err := client.Subscribe()
	assert.Nil(t, result)
	assert.NotNil(t, err)
}
