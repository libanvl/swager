package main

import (
  "flag"
  "fmt"
  "log"
  "net"
  "net/rpc"
  "os"
  "os/signal"

  "github.com/libanvl/swager/blocks"
  "github.com/libanvl/swager/internal/comm"
  "github.com/libanvl/swager/internal/core"
)

func main() {
  debug := flag.Bool("debug", false, "Whether to log debug messages")
  flag.Parse()

  addr, err := comm.GetSwagerSocket()
  if err != nil {
    log.Fatal("swager socket error:", addr)
  }

  if _, err := os.Stat(addr); !os.IsNotExist(err) {
    if err := os.RemoveAll(addr); err != nil {
      log.Fatal("socket cannot be reset:", addr)
    }
  }

  blocks.RegisterBlocks()
  log.Printf("registered blocks: %v\n", len(core.Blocks))

  logch := make(chan core.BlockLogMessage, 10)
  ctrlch := make(chan *comm.ControlArgs)
  opts := core.Options{Debug: *debug, Log: logch}
  config := comm.ServerConfig{
    Blocks: core.Blocks,
    Ctrl:   ctrlch,
  }

  server, err := comm.CreateServer(&config, &opts)
  if err != nil {
    log.Fatal("failed creating server:", err)
  }

  signalch := make(chan os.Signal)
  signal.Notify(signalch, os.Interrupt)
  signal.Notify(signalch, os.Kill)

  listener, err := net.Listen("unix", addr)
  if err != nil {
    log.Fatal("failed listening on socket:", err)
  }
  defer os.RemoveAll(addr)

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
      goto cleanup
    case _ = <-signalch:
      goto cleanup
    }
  }

  cleanup:
  fmt.Println("SHUTTING DOWN...")
  close(ctrlch)
  close(signalch)
  close(logch)
  fmt.Println("GOODBYE")
}
