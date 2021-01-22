package ipc

import (
	"encoding/binary"
	"io"
)

func (c client) ipccall(pt PayloadType, payload []byte) ([]byte, error) {
	if err := c.write(pt, payload); err != nil {
		return nil, err
	}

	return c.read()
}

func (c client) write(pt PayloadType, payload []byte) error {
	h := Header(pt, len(payload))
	if err := binary.Write(c, binary.LittleEndian, h); err != nil {
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

func (c client) read() ([]byte, error) {
	var h header
	if err := binary.Read(c, binary.LittleEndian, &h); err != nil {
		return nil, err
	}

	buf := make([]byte, int(h.PayloadLength))
	_, err := io.ReadFull(c, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
