package ipc

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

func (c *Client) ipccall(pt PayloadType, payload []byte) ([]byte, error) {
	c.ipcmx.Lock()
	defer c.ipcmx.Unlock()

	if err := c.write(pt, payload); err != nil {
		return nil, err
	}

	return c.read()
}

func (c *Client) ipccallraw(pt PayloadType, payload []byte) (string, error) {
	res, err := c.ipccall(pt, payload)
	if err != nil {
		return "", nil
	}
	return string(res), nil
}

func (c *Client) write(pt PayloadType, payload []byte) error {
	h := NewHeader(pt, len(payload))
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

func (c *Client) read() ([]byte, error) {
	var h Header
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

func callgetptr[T interface{}](c *Client, pt PayloadType, payload []byte) (*T, error) {
	res, err := c.ipccall(pt, payload)
	if err != nil {
		return nil, err
	}

	t := new(T)
	if err := json.Unmarshal(res, t); err != nil {
		return nil, err
	}

	return t, nil
}

func callgetarr[T interface{}](c *Client, pt PayloadType, payload []byte) ([]T, error) {
	res, err := c.ipccall(pt, payload)
	if err != nil {
		return nil, err
	}

	var values []T
	if err := json.Unmarshal(res, &values); err != nil {
		return nil, err
	}

	return values, nil
}
