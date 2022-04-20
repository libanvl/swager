package comm

import (
	"io"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type Swager struct {
	Client     *ipc.Client
	Sub        *ipc.Subscription
	cfg        *ServerConfig
	opts       *core.Options
	initalized map[string]core.BlockInitializer
}

type ServerConfig struct {
	Blocks core.BlockRegistry
	Ctrl   chan<- *ControlArgs
	Log    core.LogChannel
}

func CreateServer(cfg *ServerConfig, opts *core.Options) (*Swager, error) {
	client, err := ipc.Connect()
	if err != nil {
		return nil, err
	}

	sub, err := ipc.Subscribe()
	if err != nil {
		return nil, err
	}

	suberrors := make(chan error, 3)
	go func() {
		for serr := range suberrors {
			cfg.Log.Sendf(core.DefaultLog, "server", "%s", serr)
		}
	}()

	sub.Errors(suberrors)

	swager := new(Swager)
	swager.Client = client
	swager.Sub = sub
	swager.opts = opts
	swager.cfg = cfg

	return swager, nil
}

func (s *Swager) InitBlock(args *InitBlockArgs, reply *Reply) error {
	blockfac, ok := s.cfg.Blocks[args.Block]
	if !ok {
		return &BlockNotFoundError{args.Block}
	}

	block := blockfac()
	if err := block.Init(s.Client, s.Sub, s.opts, core.NewPrefixLogger(args.Block, s.cfg.Log), args.Args...); err != nil {
		return &BlockInitializationError{err, args.Block}
	}

	s.cfg.Log.Sendf(core.InfoLog, "server", "<%s>(%s) configured", args.Block, args.Tag)

	if s.initalized == nil {
		s.initalized = map[string]core.BlockInitializer{args.Tag: block}
	} else {
		s.initalized[args.Tag] = block
	}

	reply.Args = args
	reply.Success = true
	return nil
}

func (s *Swager) SendToTag(args *SendToTagArgs, reply *Reply) error {
	block, ok := s.initalized[args.Tag]
	if !ok {
		return &TagNotFoundError{args.Tag}
	}

	rcv, ok := block.(core.Receiver)
	if !ok {
		return &TagCannotReceiveError{args.Tag}
	}

	if err := rcv.Receive(args.Args); err != nil {
		return &TagReceiveError{err, args.Tag}
	}

	s.cfg.Log.Sendf(core.InfoLog, "server", "(%s) received args: %v", args.Tag, args.Args)
	reply.Success = true
	return nil
}

func (s *Swager) SetTagLog(args *SetTagLogArgs, reply *Reply) error {
	block, ok := s.initalized[args.Tag]
	if !ok {
		return &TagNotFoundError{args.Tag}
	}

	block.SetLogLevel(args.Level)

	s.cfg.Log.Sendf(core.InfoLog, "server", "(%s) set log level: %v", args.Tag, args.Level)

	reply.Args = args
	reply.Success = true
	return nil
}

func (s *Swager) Control(args *ControlArgs, reply *Reply) error {
	switch args.Command {
	case PingServer:
		s.cfg.Log.Send(core.DefaultLog, "server", "pong")
		reply.Success = true
		return nil
	case RunServer:
		s.cfg.Log.Send(core.DefaultLog, "server", "running initalized blocks")
		for _, block := range s.initalized {
			runner, ok := block.(core.Runner)
			if ok {
				go runner.Run()
			}
		}
		go s.Sub.Run()
		reply.Args = args
		reply.Success = true
		return nil
	case ResetServer:
		s.cfg.Log.Send(core.DefaultLog, "server", "resetting initalized blocks")
		closeAllBlocks(s)
		s.Sub.Close()
		sub, err := ipc.Subscribe()
		if err != nil {
			return err
		}
		s.Sub = sub
		reply.Args = args
		reply.Success = true
		return nil
	case ExitServer:
		closeAllBlocks(s)
		s.Sub.Close()
		fallthrough
	default:
		s.cfg.Ctrl <- args
	}

	reply.Args = args
	reply.Success = true
	return nil
}

func closeAllBlocks(s *Swager) {
	for tag, block := range s.initalized {
		closer, ok := block.(io.Closer)
		if ok {
			closer.Close()
		}
		delete(s.initalized, tag)
		s.cfg.Log.Sendf(core.DefaultLog, "server", "(%s) closed", tag)
	}
}
