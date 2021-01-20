package ipc

import (
  "encoding/binary"
  "encoding/json"
  "errors"
  "io"
  "log"

  "github.com/libanvl/swager/pkg/ipc/event"
)

type WindowSubscription struct {
  client Client
  observers []chan *event.WindowReply
}

func (s *WindowSubscription) Observe(ch chan *event.WindowReply) {
  if s.observers == nil {
    s.observers = make([]chan *event.WindowReply, 1, 5)
  }

  s.observers = append(s.observers, ch)
}

func (s *WindowSubscription) Start() error {
  client, err := Connect()
  if err != nil {
    log.Print("Failed connect")
    return err
  }

  res, err := client.Subscribe("window", "shutdown")
  if err != nil {
    log.Print("Failed subscribe call")
    return err
  }

  if !res.Success {
    return errors.New("Failed to subscribe to window events")
  }

  s.client = client

  for {
    var h header
    if err := binary.Read(s.client.Conn(), binary.LittleEndian, &h); err != nil {
      log.Printf("Failed reading event header: %v", err)
      //      for _, o := range s.observers {
      //        close(o)
      //      }
      //      return err
    }

    log.Printf("Recevied header: %#v", h)

    buf := make([]byte, int(h.PayloadLength))
    _, err := io.ReadFull(s.client.Conn(), buf)
    if err != nil {
      log.Print("Failed reading event payload")
      //      for _, o := range s.observers {
      //        close(o)
      //      }
      return err
    }

    wr := new(event.WindowReply)
    if err := json.Unmarshal(buf, wr); err != nil {
      log.Print("Failed to unmarshal response")
      //      for _, o := range s.observers {
      //        close(o)
      //      }
      return err
    }

    log.Print("Notifying observers")
    for x, observer := range s.observers {
      go func(i int, ch chan *event.WindowReply) {
        log.Printf("Observer: %v", i)
        ch <- wr
      }(x, observer)
    }
  }
}
