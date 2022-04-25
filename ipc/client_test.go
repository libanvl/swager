package ipc_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"testing"

	"github.com/libanvl/swager/ipc"
)

func init() {
	var _ io.ReadWriteCloser = (*MockConnection)(nil)
}

type MockConnection struct {
	t          *testing.T
	Written    [][]byte
	ReadCount  uint
	CloseCount uint
}

// Read implements io.ReadWriteCloser
func (mc *MockConnection) Read(p []byte) (n int, err error) {
	mc.ReadCount++
	return len(p), nil
}

// Write implements io.ReadWriteCloser
func (mc *MockConnection) Write(p []byte) (n int, err error) {
	mc.Written = append(mc.Written, p)
	mc.t.Logf("Received: %#v", string(p))
	return len(p), nil
}

// Close implements io.ReadWriteCloser
func (mc *MockConnection) Close() error {
	mc.CloseCount++
	return nil
}

func TestConnect(t *testing.T) {
	tmpsocket := t.TempDir() + "/uds"
	net.Listen("unix", tmpsocket)
	t.Setenv("SWAYSOCK", tmpsocket)
	client, err := ipc.Connect()

	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatal()
	}
}

func TestConnectNoSwaysock(t *testing.T) {
	t.Setenv("SWAYSOCK", "")
	_, err := ipc.Connect()

	if err == nil {
		t.Errorf("Did not get error as expected")
	}
}

func TestPayloadCalls(t *testing.T) {
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
			conn := MockConnection{t, make([][]byte, 0), 0, 0}
			client := ipc.NewClient(&conn, binary.LittleEndian)

			tc.action(client, tc.payload)

			var got_header ipc.Header
			r := bytes.NewReader(conn.Written[0])
			if err := binary.Read(r, binary.LittleEndian, &got_header); err != nil {
				t.Error(err)
			}

			if !ipc.ValidMagic(got_header.Magic) {
				t.Errorf("Magic bytes in header not valid %#v", got_header.Magic)
			}

			if got_header.PayloadLength != uint32(len(tc.payload)) {
				t.Errorf("Wrong Payload Length in Header. Want: %v. Got: %v", len(tc.payload), got_header.PayloadLength)
			}

			got_payload := string(conn.Written[1])
			if tc.payload != got_payload {
				t.Errorf("Wrong Payload. Want: %v, Got: %v", tc.payload, got_payload)
			}

			t.Logf("Connection Read %v times", conn.ReadCount)
			t.Logf("Connection Close %v times", conn.CloseCount)
		})
	}
}

func TestNoPayloadCalls(t *testing.T) {
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
			conn := MockConnection{t, make([][]byte, 0), 0, 0}
			client := ipc.NewClient(&conn, binary.LittleEndian)

			tc.action(client)

			var got_header ipc.Header
			r := bytes.NewReader(conn.Written[0])
			if err := binary.Read(r, binary.LittleEndian, &got_header); err != nil {
				t.Error(err)
			}

			if !ipc.ValidMagic(got_header.Magic) {
				t.Errorf("Magic bytes in header not valid %#v", got_header.Magic)
			}

			if got_header.PayloadLength != 0 {
				t.Errorf("Wrong Payload Length in Header. Want: %v. Got: %v", 0, got_header.PayloadLength)
			}

			if len(conn.Written) != 1 {
				t.Errorf("Wrong Write Count. Want: %v, Got: %v", 1, len(conn.Written))
			}

			t.Logf("Connection Read %v times", conn.ReadCount)
			t.Logf("Connection Close %v times", conn.CloseCount)
		})
	}
}
