package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"

	"github.com/libanvl/swager/internal/blocks"
	"github.com/libanvl/swager/internal/comm"
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/ipc/event"
)

func main() {
  addr, err := comm.GetSwagerSocket()
  if err != nil {
    log.Fatal("swager socket error:", addr)
  }

  if _, err := os.Stat(addr); !os.IsNotExist(err) {
    if err := os.RemoveAll(addr); err != nil {
      log.Fatal("socket already exists:", err)
    }
  }

  blocks.RegisterBlocks()
  log.Printf("registered blocks: %v\n", len(core.Blocks))

	sub := event.Subscribe()
  client, err := ipc.Connect()
	if err != nil {
    log.Fatal("failed getting sway client:", err)
	}

  signalch := make(chan os.Signal)
  signal.Notify(signalch, os.Interrupt)
  signal.Notify(signalch, os.Kill)

	listener, err := net.Listen("unix", addr)
	if err != nil {
    log.Fatal("failed listening on socket:", err)
	}
  defer os.RemoveAll(addr)

  logch := make(chan string, 10)
  ctrlch := make(chan *comm.ControlArgs)
	opts := core.Options{Debug: true, Log: logch}
  config := comm.ServerConfig {
    Blocks: core.Blocks,
    Client: client,
    Sub: sub,
    Ctrl: ctrlch,
  }

	server := comm.CreateServer(&config, &opts)
  rpc.Register(server)

  log.Println("rpc server registered")

  go rpc.Accept(listener)

  fmt.Println("SOCKET:", addr)

  for {
    select {
    case l := <-logch:
      log.Println(l)
    case cmdargs := <-ctrlch:
      if cmdargs.Command != comm.ExitServer {
        continue
      }
      goto done
    case _ = <-signalch:
      goto done
    }
  }

done:
  fmt.Println("SHUTTING DOWN...")
  close(ctrlch)
  close(signalch)
  close(logch)
  fmt.Println("GOODBYE")
}
