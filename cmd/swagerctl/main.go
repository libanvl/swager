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

  go func(ch chan os.Signal) {
    _ = <-ch
    os.Exit(1)
  }(signalch)

  if _, err := os.Stat(addr); os.IsNotExist(err) {
    log.Fatal("swager daemon error:", err)
  }

  conn, err := net.DialTimeout("unix", addr, time.Millisecond * 500)
  if err != nil {
    log.Fatal("failed dialing rpc:", err)
  }

  client := rpc.NewClient(conn)

  reply := new(comm.Reply)
  if err := client.Call(string(comm.Control),&comm.ControlArgs{Command: comm.PingServer}, reply); err != nil {
    log.Fatal("swagrd did not respond:", err)
  }

  var op comm.SwagerMethod
  var args interface{}

  switch os.Args[1] {
  case "init":
    // swagerctl init tag block config0 config1 config2 ...
    op = comm.InitBlock
    args = &comm.InitBlockArgs{Tag: os.Args[2], Block: os.Args[3], Config: os.Args[4:]}
    break
  case "send":
    // swagerctl send tag arg0 arg1 arg2 ...
    op = comm.SendToTag
    args = &comm.SendToTagArgs{Tag: os.Args[2], Args: os.Args[3:]}
    break
  case "run":
    op = comm.Control
    args = &comm.ControlArgs{Command: comm.RunServer}
    break
  case "exit":
    // swagerctl exit
    op = comm.Control
    args = &comm.ControlArgs{Command: comm.ExitServer}
    break
  case "ping":
    // swagerctl ping
    op = comm.Control
    args = &comm.ControlArgs{Command: comm.PingServer}
    break
  default:
    log.Fatal("unknown method:", os.Args[1])
  }

  if err := client.Call(string(op), args, reply); err != nil {
    log.Fatal("swager err:", err)
  } else {
    log.Printf("%#v\n", reply)
  }
}

