package main

import (
  "log"
  "net"
  "net/rpc"
  "os"
  "os/signal"
  "time"

  "github.com/libanvl/swager/internal/comm"
)

func main() {
  addr, err := comm.GetSwagerSocket()
  if err != nil {
    log.Fatal("swager socket error:", err)
  }

  signalch := make(chan os.Signal)
  signal.Notify(signalch, os.Interrupt)
  signal.Notify(signalch, os.Kill)

  if _, err := os.Stat(addr); os.IsNotExist(err) {
    log.Fatal("swager daemon error:", err)
  }

  conn, err := net.DialTimeout("unix", addr, time.Millisecond*500)
  if err != nil {
    log.Fatal("failed dialing rpc:", err)
  }

  client := rpc.NewClient(conn)

  go func(ch chan os.Signal, client *rpc.Client) {
    _ = <-ch
    client.Close()
    os.Exit(1)
  }(signalch, client)

  reply := new(comm.Reply)

  var op comm.SwagerMethod
  var args interface{}

  switch os.Args[1] {
  case "init":
    // swagerctl init tag block arg0 arg1 arg2 ...
    op = comm.InitBlock
    args = &comm.InitBlockArgs{Tag: os.Args[2], Block: os.Args[3], Config: os.Args[4:]}
    break
  case "send":
    // swagerctl send tag arg0 arg1 arg2 ...
    op = comm.SendToTag
    args = &comm.SendToTagArgs{Tag: os.Args[2], Args: os.Args[3:]}
    break
  case "config":
    if len(os.Args) != 3 {
      log.Fatal("config requires a subcommand")
    }
    switch os.Args[2] {
    case "commit":
      // swagerctl config commit
      op = comm.Control
      args = &comm.ControlArgs{Command: comm.RunServer}
      break
    case "reset":
      // swagerctl config reset
      op = comm.Control
      args = &comm.ControlArgs{Command: comm.ResetServer}
      break
    default:
      log.Fatal("unknown method:", os.Args[1], os.Args[2])
    }
  case "server":
    if len(os.Args) != 3 {
      log.Fatal("server requires a subcommand")
    }
    switch os.Args[2]{
    case "exit":
      // swagerctl server exit
      op = comm.Control
      args = &comm.ControlArgs{Command: comm.ExitServer}
      break
    case "ping":
      // swagerctl server ping
      op = comm.Control
      args = &comm.ControlArgs{Command: comm.PingServer}
      break
    default:
      log.Fatal("unknown method:", os.Args[1], os.Args[2])
    }
  default:
    log.Fatal("unknown method:", os.Args[1])
  }

  if err := client.Call(string(op), args, reply); err != nil {
    log.Fatal("swager err:", err)
  } else {
    log.Printf("%#v\n", reply)
  }
}
