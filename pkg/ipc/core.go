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
  msg := newMessage(pt, payload)
  if err := binary.Write(c.conn, binary.LittleEndian, msg.header); err != nil {
    return err
  }

  if msg.PayloadLength > 0 {
    _, err := c.conn.Write(msg.payload)
    if err != nil {
      return err
    }
  }

  return nil
}

func (c client) read() ([]byte, error) {
  var h header
  if err := binary.Read(c.conn, binary.LittleEndian, &h); err != nil {
    return nil, err
  }

  buf := make([]byte, int(h.PayloadLength))
  _, err := io.ReadFull(c.conn, buf)
  if err != nil {
    return nil, err
  }

  return buf, nil
}

