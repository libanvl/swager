package comm

import (
	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
)

type Swager struct {
	Client     *ipc.Client
	Sub        *ipc.Subscription
	cfg        *ServerConfig
	opts       *core.Options
	initalized map[string]core.Block
}

type ServerConfig struct {
	Blocks core.BlockRegistry
	Ctrl   chan<- *ControlArgs
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
	if err := block.Init(s.Client, s.Sub, s.opts); err != nil {
		return &BlockInitializationError{err, args.Block}
	}

	s.opts.Log.Infof("server", "<%s>(%s) initalized", args.Block, args.Tag)

	if err := block.Configure(args.Config); err != nil {
		return &BlockInitializationError{err, args.Block}
	}

	s.opts.Log.Defaultf("server", "<%s>(%s) configured", args.Block, args.Tag)

	if s.initalized == nil {
		s.initalized = map[string]core.Block{args.Tag: block}
	} else {
		s.initalized[args.Tag] = block
	}

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

	s.opts.Log.Debugf("server", "(%s) received args: %v", args.Tag, args.Args)

	reply.Success = true
	return nil
}

func (s *Swager) Control(args *ControlArgs, reply *Reply) error {
	switch args.Command {
	case PingServer:
		s.opts.Log.Default("server", "pong")
		reply.Success = true
		return nil
	case RunServer:
		s.opts.Log.Default("server", "running initalized blocks")
		for _, block := range s.initalized {
			go block.Run()
		}
		go s.Sub.Run()
		reply.Success = true
		return nil
	case ResetServer:
		s.opts.Log.Default("server", "resetting initalized blocks")
		s.Sub.Close()
		sub, err := ipc.Subscribe()
		if err != nil {
			return err
		}
		s.Sub = sub
		closeAllBlocks(s)
		reply.Success = true
		return nil
	case ExitServer:
		s.Sub.Close()
		closeAllBlocks(s)
		fallthrough
	default:
		s.cfg.Ctrl <- args
	}
	reply.Success = true
	return nil
}

func closeAllBlocks(s *Swager) {
	for tag, block := range s.initalized {
		block.Close()
		delete(s.initalized, tag)
		s.opts.Log.Defaultf("server", "(%s) closed", tag)
	}
}
