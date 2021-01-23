package main

import (
	"flag"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/libanvl/swager/internal/blocks"
	"github.com/libanvl/swager/internal/comm"
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/ipc/event"
)

func main() {
  var addr string

  flag.StringVar(&addr, "socket", "", "Path to the unix domain socket")
  flag.Parse()

  if addr == "" {
    log.Fatalln("A value for socket is required.")
  }

  blocks.RegisterBlocks()

	sub := event.Subscribe()

  client, err := ipc.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	os.RemoveAll(addr)
	listener, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatalln(err)
	}

  defer os.RemoveAll(addr)

	logch := make(chan string)
	opts := core.Options{Debug: true, Log: logch}

  ctrlch := make(chan *comm.ControlArgs)
  config := comm.ServerConfig {
    Blocks: core.Blocks,
    Client: client,
    Sub: sub,
    Ctrl: ctrlch,
  }

	server := comm.CreateServer(&config, &opts)

  rpc.Register(server)

	go sub.Start()
  defer sub.Close()

  go rpc.Accept(listener)

  log.Printf("swagerd listening on %s", addr)

  for cargs := range ctrlch {
    if cargs.Command == comm.ExitServer {
      log.Println("Exiting server. Goodbye!")
      break
    }
  }
}
