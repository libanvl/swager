package main

import (
  "fmt"
  "log"

  "github.com/libanvl/swager/pkg/ipc"
  "github.com/libanvl/swager/pkg/ipc/event"
)

func main() {
  c, err := ipc.Connect()
  if err != nil {
    log.Fatalf("Could not connect: %v", err)
  }

  defer c.Close()
  v, err := c.Version()
  if err != nil {
    log.Fatalf("Failed getting version: %v", err)
  }

  log.Println("Version")
  log.Printf("%#v", v)

  ws, err := c.Workspaces()
  if err != nil {
    log.Fatalf("Failed getting workspaces: %v", err)
  }

  log.Println("Workspaces")
  log.Printf("%#v", ws)

  cr := c.ClientRaw()
  vr, err := cr.VersionRaw()
  if err != nil {
    log.Fatalf("Failed version raw: %v", err)
  }

  log.Print("Version Raw")
  log.Print(vr)

  wr, err := cr.WorkspacesRaw()
  if err != nil {
    log.Fatalf("Failed workspaces raw: %v", err)
  }

  log.Print("Workspaces Raw")
  log.Print(wr)

  winevts := make(chan *event.WindowReply)
  sub := new(ipc.WindowSubscription)
  sub.Observe(winevts)

  go receiveWindows(winevts)

  if err := sub.Start(); err != nil {
    log.Printf("Failed starting sub: %v", err)
  }
}

func receiveWindows(ch chan *event.WindowReply) {
  for evt := range ch {
    fmt.Printf("EVENT! %#v", evt)
  }
  fmt.Print("CLOSED!")
}
