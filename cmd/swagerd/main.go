package main

import (
	"log"
	"net"
	"os"

	"github.com/adrg/xdg"

	"github.com/libanvl/swager/internal/blocks"
	"github.com/libanvl/swager/internal/comm"
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

func main() {
	blocks.RegisterBlocks()

	client, err := ipc.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	sub := ipc.Subscribe()

	addr, err := xdg.RuntimeFile("swager/swager.sock")
	if err != nil {
		log.Fatalln(err)
	}

	os.RemoveAll(addr)
	listener, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatalln(err)
	}

	logch := make(chan string)

	opts := core.Options{Debug: true, Log: logch}

	server, err := comm.CreateServer(core.Blocks, listener, client, sub, &opts)
	if err != nil {
		log.Fatalln(err)
	}

	exitch := make(chan bool)

	server.NotifyExit(exitch)

	go sub.Start()
	go server.Start()
	defer sub.Close()

	<-exitch
}
