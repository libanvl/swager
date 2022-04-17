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

var loglevel core.LogLevel

func main() {
	flag.Var(&loglevel, "log", "the log level")
	flag.Parse()

	log.Println("log level: ", loglevel)

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

	logch := make(chan core.LogMessage, 10)
	creqch := make(chan core.ServerControlRequest)
	ctrlch := make(chan *comm.ControlArgs)
	opts := core.Options{Server: creqch}
	config := comm.ServerConfig{
		Blocks: core.Blocks,
		Ctrl:   ctrlch,
		Log:    logch,
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
		case r := <-creqch:
			if r == core.ExitRequest {
				go server.Control(&comm.ControlArgs{Command: comm.ExitServer}, &comm.Reply{})
				goto cleanup
			}
			if r == core.ReloadRequest {
				go server.Control(&comm.ControlArgs{Command: comm.ResetServer}, &comm.Reply{})
			}
		case l := <-logch:
			if loglevel >= l.Level() {
				log.Println(l)
			}
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
	fmt.Println("...GOODBYE")
}
