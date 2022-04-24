package ipc

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestConnect(t *testing.T) {
	tmpsocket := t.TempDir() + "/uds"
	net.Listen("unix", tmpsocket)
	t.Setenv("SWAYSOCK", tmpsocket)
	client, err := Connect()

	if err != nil {
		t.Fatal(err)
	}

	if client.yo != binary.LittleEndian {
		t.Errorf("client byteorder %v not equal to expected %v", client.yo, binary.LittleEndian)
	}
}

func TestConnectNoSwaysock(t *testing.T) {
	t.Setenv("SWAYSOCK", "")
	_, err := Connect()

	if err == nil {
		t.Errorf("Did not get error as expected")
	}
}
