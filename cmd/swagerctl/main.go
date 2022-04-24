package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libanvl/swager/internal/comm"
	"github.com/libanvl/swager/stoker"
)

func main() {
	var ctimeout = time.Duration(2 * time.Second)
	var print_usage = func(_ any, _ stoker.TokenList) error {
		fmt.Println(usage())
		os.Exit(0)
		return nil
	}

	if len(os.Args) < 1 {
		print_usage(nil, nil)
	}

	flag_parse := stoker.NewParser[any](
		stoker.NewFlag("-h", print_usage),
		stoker.NewFlag("--help", print_usage),
		stoker.NewFlag("-c", func(_ any, tl stoker.TokenList) error {
			var err error
			ctimeout, err = time.ParseDuration(tl[0])
			if err != nil {
				return err
			}

			return nil
		}),
	)

	if err := flag_parse.Parse(os.Args...).HandleAll(nil); err != nil {
		log.Fatal("argument parsing error: ", err)
	}

	parser := stoker.NewParser[*rpc.Client](
		stoker.NewFlag("--server", listHandler(comm.Control)),
		stoker.NewFlag("--init", listHandler(comm.InitBlock)),
		stoker.NewFlag("--log", listHandler(comm.SetTagLog)),
		stoker.NewFlag("--send", listHandler(comm.SendToTag)),
	)

	handler := parser.Parse(os.Args...)

	addr := getSwagerSocketAddress(&ctimeout)

	if _, err := os.Stat(addr); os.IsNotExist(err) {
		log.Fatal("swager daemon error: ", err)
	}

	conn, err := net.DialTimeout("unix", addr, 1*time.Second)
	if err != nil {
		log.Fatal("failed dialing ipc: ", err)
	}

	client := rpc.NewClient(conn)
	attachSignals(client)

	if err = handler.HandleAll(client); err != nil {
		log.Fatal(err)
	}
}

func attachSignals(client *rpc.Client) {
	signalch := make(chan os.Signal, 3)
	signal.Notify(signalch, os.Interrupt)
	signal.Notify(signalch, syscall.SIGTERM)

	go func(ch chan os.Signal, client *rpc.Client) {
		<-ch
		client.Close()
		os.Exit(1)
	}(signalch, client)
}

func getSwagerSocketAddress(flagCtimeout *time.Duration) string {
	sokch := make(chan string)
	go func(ch chan<- string) {
		for {
			addr, err := comm.GetSwagerSocket()
			if err != nil {
				log.Print("swager socket retry: ", err)
				continue
			}
			ch <- addr
			break
		}
	}(sokch)

	var addr string
	select {
	case addr = <-sokch:
		close(sokch)
		break
	case <-time.After(*flagCtimeout):
		log.Fatal("swager socket timeout error")
	}

	return addr
}

func listHandler(op comm.SwagerMethod) stoker.TokenListHandler[*rpc.Client] {
	return func(c *rpc.Client, tokenlist stoker.TokenList) error {
		args, err := toSwagerArgs(op, tokenlist)
		if err != nil {
			return err
		}
		if err := call(c, op, args, new(comm.Reply)); err != nil {
			return err
		}
		return nil
	}
}

func toSwagerArgs(op comm.SwagerMethod, tl stoker.TokenList) (comm.SwagerArgs, error) {
	switch op {
	case comm.InitBlock:
		return comm.ToInitBlockArgs(tl)
	case comm.SendToTag:
		return comm.ToSendToTagArgs(tl)
	case comm.Control:
		return comm.ToServerArgs(tl)
	case comm.SetTagLog:
		return comm.ToSetTagLogArgs(tl)
	}

	return nil, errors.New("unknown error")
}

func call(client *rpc.Client, op comm.SwagerMethod, a interface{}, reply *comm.Reply) error {
	if err := client.Call(string(op), a, reply); err != nil {
		return errors.New(fmt.Sprintf("swager err: %#v", err))
	} else {
		log.Printf("%#v\n", reply)
	}

	return nil
}

func usage() string {
	help := `swagerctl [<flags>] <method> [args...]

  flags:
  -c <duration> time to wait for swagerd connection
  -h help

  methods:
  --server - send a server control command
  --init   - initialze a new block instance
  --log    - set log level on a block
  --send   - send a command to an initialized block instance

  --init <tagname> <blockname> [args...]

    <tagname> is a user-provided name for a specific block instance
    <blockname> is the registered name for a block type
    [args...] are the initialization arguments for the block instance
    see the block documentation for the supported args

    examples:
      --init myauto autolay -masterstack 1 2 3 4
      --init myexecnew execnew 1 10

  --log <tagname> <loglevel>

    <tagname> is a user-provided name for a specific block instance
    <loglevel> is one of: default, info, debug

    examples:
      --init myauto debug
      --init myauto default

  --send <tagname> arg0 [args...]

    <tagname> is the user-provided name for a block instance
    arg0 [args...] are the arguments to send to the block instance
    not all block types support receiving arguments using send
    see the block documentation for the supported args

    examples:
      --send myexecnew "exec alacritty"

  --server <submethod>

    server method should be the only method in a call to swagerctl

    submethods:
      listen - start monitoring events, no more blocks can be initialized
      reset  - stop monitoring events, close all initialized blocks
      ping   - ping the server
      exit   - notify the server to shutdown`

	return help
}
