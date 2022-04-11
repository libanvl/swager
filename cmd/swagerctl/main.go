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
	"github.com/libanvl/swager/pkg/stoker"
)

func main() {
	var flagCtimeout time.Duration = 2 * time.Second

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println(usage())
		os.Exit(0)
	}

	parser := stoker.Parser(
		stoker.DefArgs("-h", false),
		stoker.Def("-c"),
		stoker.Def("--server"),
		stoker.Def("--init"),
		stoker.Def("--log"),
		stoker.Def("--send"),
		stoker.Def("--config"),
	)

	tokenmap := parser.Parse(args...)

	if parser.Present.Contains("-h") {
		fmt.Println(usage())
		os.Exit(0)
	}

	if parser.Present.Contains("-c") {
		var err error
		flagCtimeout, err = time.ParseDuration(tokenmap["-c"][0][0])
		if err != nil {
			fmt.Println(usage())
			os.Exit(1)
		}
	}

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
	case <-time.After(flagCtimeout):
		log.Fatal("swager socket timeout error")
	}

	if _, err := os.Stat(addr); os.IsNotExist(err) {
		log.Fatal("swager daemon error: ", err)
	}

	conn, err := net.DialTimeout("unix", addr, 1*time.Second)
	if err != nil {
		log.Fatal("failed dialing rpc: ", err)
	}

	signalch := make(chan os.Signal, 3)
	signal.Notify(signalch, os.Interrupt)
	signal.Notify(signalch, syscall.SIGTERM)

	client := rpc.NewClient(conn)

	go func(ch chan os.Signal, client *rpc.Client) {
		<-ch
		client.Close()
		os.Exit(1)
	}(signalch, client)

	reply := new(comm.Reply)

	if err := tokenmap.ProcessSet("--server", func(ts stoker.TokenSet) error {
		for _, sc := range ts {
			args, err := comm.ToServerArgs(sc)
			if err != nil {
				return err
			}
			if err := call(client, comm.Control, args, reply); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Fatal("error: ", "--server ", err)
	}

	if err := tokenmap.ProcessSet("--init", processTokenSet(client, reply, comm.InitBlock)); err != nil {
		log.Fatal("error: ", "--init ", err)
	}

	if err := tokenmap.ProcessSet("--log", processTokenSet(client, reply, comm.SetTagLog)); err != nil {
		log.Fatal("error: ", "--log ", err)
	}

	if err := tokenmap.ProcessSet("--send", processTokenSet(client, reply, comm.SendToTag)); err != nil {
		log.Fatal("error: ", "--send ", err)
	}

	if err := tokenmap.ProcessSet("--config", processTokenSet(client, reply, comm.Control)); err != nil {
		log.Fatal("error: ", "--config ", err)
	}
}

func processTokenSet(client *rpc.Client, reply *comm.Reply, op comm.SwagerMethod) stoker.TokenSetProcessor {
	return func(ts stoker.TokenSet) error {
		for _, tokenlist := range ts {
			args, err := toSwagerArgs(op, tokenlist)
			if err != nil {
				return err
			}
			if err := call(client, op, args, reply); err != nil {
				return err
			}
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
		return comm.ToConfigArgs(tl)
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
	help := `swagerctl [<flags>] <method> [<submethod>] [args...]

  flags:
  -c duration - time to wait for a connection to the swager daemon (default 2s)
  -h help

  methods:
  --server - send a server control command
  --init   - initialze a new block instance
  --log    - set log level on a block
  --send   - send a command to an initialized block instance
  --config - send a configuration command

  --init <tagname> <blockname> [args...]

    <tagname> is a user-provided name for a specific block instance
    <blockname> is the registered name for a block type
    [args...] are the initialization arguments for the block instance
    see the block documentation for the supported args

    examples:
      --init mytiler tiler
      --init myexecnew execnew 1 10

  --log <tagname> <loglevel>

    <tagname> is a user-provided name for a specific block instance

    examples:
      --init mytiler debug
      --init mytiler default

  --send <tagname> arg0 [args...]

    <tagname> is the user-provided name for a block instance
    arg0 [args...] are the arguments to send to the block instance
    not all block types support receiving arguments using send
    see the block documentation for the supported args

    examples:
      --send myexecnew "exec alacritty"

  --config <submethod>

    submethods:
      commit - notify the server that user configuration is complete, and to start monitoring events
      reset  - notify the server to stop monitoring events, and close all configured blocks

  --server <submethod>

    submethods:
      ping - ping the server
      exit - notify the server to shutdown`

	return help
}
