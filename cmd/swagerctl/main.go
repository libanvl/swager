package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"time"

	"github.com/libanvl/swager/internal/comm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage())
		return
	}

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

	conn, err := net.DialTimeout("unix", addr, time.Second * 1)
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
		if len(os.Args) < 3 {
			log.Fatal("config requires a subcommand\n", usage())
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
			log.Fatalf("unknown method: %s %s\n %s", os.Args[1], os.Args[2], usage())
		}
	case "server":
		if len(os.Args) < 3 {
			log.Fatal("server requires a subcommand\n", usage())
		}
		switch os.Args[2] {
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
			log.Fatalf("unknown method: %s %s\n%s", os.Args[1], os.Args[2], usage())
		}
	default:
		log.Fatalf("unknown method: %s\n%s", os.Args[1], usage())
	}

	if err := client.Call(string(op), args, reply); err != nil {
		log.Fatal("swager err:", err)
	} else {
		log.Printf("%#v\n", reply)
	}
}

func usage() string {
	help := `swayctl <method> [<submethod>] [args...]

  methods:
    init   - initialze a new block instance
    send   - send a command to an initialized block instance
    config - send a configuration command
    server - send a server control command

  init <tagname> <blockname> [args...]

    <tagname> is a user-provided name for a specific block instance
    <blockname> is the registered name for a block type
    [args...] are the initialization arguments for the block instance
      see the block documentation for the supported args

    examples:
      init mytiler tiler
      init myexecnew execnew 1 10

  send <tagname> arg0 [args...]

    <tagname> is the user-provided name for a block instance
    arg0 [args...] are the argument to send to the block instance
      not all block types support receiving arguments using send
      see the block documentation for the supported args

    examples:
      send myexecnew "exec alacritty"

  config <submethod>

    submethods:
      commit - notify the server that user configuration is complete, and to start monitoring events
      reset  - notify the server to stop monitoring events, and close all configured blocks

  server <submethod>

    submethods:
      ping - ping the server
      exit - notify the server to shutdown`

	return help
}
