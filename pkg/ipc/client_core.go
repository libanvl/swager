package ipc

import (
	"encoding/binary"
	"io"
)

func (c Client) ipccall(pt payloadType, payload []byte) ([]byte, error) {
	if err := c.write(pt, payload); err != nil {
		return nil, err
	}

	return c.read()
}

func (c Client) ipccallraw(pt payloadType, payload []byte) (string, error) {
	res, err := c.ipccall(pt, payload)
	if err != nil {
		return "", nil
	}
	return string(res), nil
}

func (c Client) write(pt payloadType, payload []byte) error {
	h := newHeader(pt, len(payload))
	if err := binary.Write(c, c.yo, h); err != nil {
		return err
	}

	if h.PayloadLength > 0 {
		_, err := c.Write(payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Client) read() ([]byte, error) {
	var h header
	if err := binary.Read(c, c.yo, &h); err != nil {
		return nil, err
	}

	buf := make([]byte, int(h.PayloadLength))
	_, err := io.ReadFull(c, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
