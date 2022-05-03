package test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/libanvl/swager/ipc"
	"github.com/stretchr/testify/assert"
)

func init() {
	var _ io.ReadWriteCloser = (*MockConnection)(nil)
}

type MockConnection struct {
	t               *testing.T
	log             bool
	callCounts      map[string]uint
	writes          [][]byte
	nextRead        []ReadValue
	nextWriteResult *struct {
		n int
		e error
	}
}

func NewMockConnection(t *testing.T) *MockConnection {
	return &MockConnection{
		t:          t,
		callCounts: make(map[string]uint)}
}

func (mc *MockConnection) Logf(msg string, args ...any) {
	if mc.log {
		mc.t.Logf(msg, args...)
	}
}

// Read implements io.ReadWriteCloser
func (mc *MockConnection) Read(p []byte) (n int, err error) {
	mc.callCounts["Read"]++

	size := len(mc.nextRead)
	if size > 0 {
		next := mc.nextRead[0]
		mc.nextRead = mc.nextRead[1:]

		if b := next.Bytes(); b != nil {
			mc.Logf("Read: %s", b)
			n := copy(p, b)
			return n, nil
		} else {
			e := next.Error()
			mc.Logf("Read: %#v", e)
			return 1, e
		}
	}

	mc.Logf("Read: 1, nil")
	return 1, nil
}

// Write implements io.ReadWriteCloser
func (mc *MockConnection) Write(p []byte) (n int, err error) {
	mc.callCounts["Write"]++
	if mc.writes == nil {
		mc.writes = make([][]byte, 0)
	}

	mc.writes = append(mc.writes, p)
	mc.Logf("Write: %s", p)

	if mc.nextWriteResult != nil {
		mc.Logf("Returned: %#v", mc.nextWriteResult)
		return mc.nextWriteResult.n, mc.nextWriteResult.e
	}

	return len(p), nil
}

// Close implements io.ReadWriteCloser
func (mc *MockConnection) Close() error {
	mc.callCounts["Close"]++
	return nil
}

func (mc *MockConnection) WriteAt(n int) []byte {
	assert.NotNil(mc.t, mc.writes)
	value := mc.writes[n]
	assert.NotNil(mc.t, value)

	return value
}

func (mc *MockConnection) SetNextWriteResult(n int, err error) {
	mc.nextWriteResult = &struct {
		n int
		e error
	}{n, err}
}

func (mc *MockConnection) PushNextReadBytes(p []byte) {
	if mc.nextRead == nil {
		mc.nextRead = make([]ReadValue, 0)
	}

	mc.nextRead = append(mc.nextRead, ReadValue{}.WithBytes(p))
}

func (mc *MockConnection) PushNextReadError(err error) {
	if mc.nextRead == nil {
		mc.nextRead = make([]ReadValue, 0)
	}

	mc.nextRead = append(mc.nextRead, ReadValue{}.WithError(err))
}

func (mc *MockConnection) PushPayloadForRead(payloadType uint32, p []byte, yo binary.ByteOrder) {
	header := ipc.NewHeader(ipc.PayloadType(payloadType), len(p))
	var buffer bytes.Buffer
	binary.Write(&buffer, yo, header)
	mc.PushNextReadBytes(buffer.Bytes())

	var buffer2 bytes.Buffer
	binary.Write(&buffer2, yo, p)
	mc.PushNextReadBytes(buffer2.Bytes())
}

func (mc *MockConnection) AssertCalled(method string) bool {
	return assert.True(mc.t, mc.callCounts[method] >= 1)
}

func (mc *MockConnection) AssertNotCalled(method string) bool {
	return assert.Zero(mc.t, mc.callCounts[method])
}

func (mc *MockConnection) AssertNumberOfCalls(method string, calls uint) bool {
	return assert.Equal(mc.t, calls, mc.callCounts[method])
}
