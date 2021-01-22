package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/ipc/event"
)

func main() {
	client()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)

	sub := ipc.Subscribe()
	windows, err := sub.Window()
	if err != nil {
		log.Fatal(err)
	}

	workspaces, err := sub.Workspace()
	if err != nil {
		log.Fatal(err)
	}

	suberrs := sub.Errors()
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(ch <-chan *event.WindowChange, wg *sync.WaitGroup) {
		for change := range ch {
			fmt.Printf("WINDOW CHANGE:\n%#v", change)
		}
		fmt.Println("WINDOW CHANGE CLOSED")
		wg.Done()
	}(windows, wg)

	wg.Add(1)
	go func(ch <-chan *event.WorkspaceChange, wg *sync.WaitGroup) {
		for change := range ch {
			fmt.Printf("WORKSPACE CHANGE:\n%#v", change)
		}
		fmt.Println("WORKSPACE CHANGE CLOSED")
		wg.Done()
	}(workspaces, wg)

	go func(ch <-chan *ipc.MonitoringError) {
		for err := range ch {
			fmt.Printf("ERROR: %v", err)
		}
	}(suberrs)

	go func() {
		if err := sub.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	<-killSignal
	fmt.Println("Closing...")
	if err := sub.Close(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}

func client() {
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

	ss, err := c.RunCommand("exec alacritty")
	if err != nil {
		log.Fatalf("Failed running command: %v", err)
	}

	log.Println(ss)
}
