package comm

import (
	"fmt"

	"github.com/libanvl/swager/internal/core"
	"github.com/libanvl/swager/pkg/ipc"
	"github.com/libanvl/swager/pkg/ipc/event"
)

type Swager struct {
	cfg        *ServerConfig
	opts       *core.Options
	initalized map[string]core.Block
}

type ServerConfig struct {
	Blocks core.BlockRegistry
	Client ipc.Client
	Sub    event.Subscription
	Ctrl   chan<- *ControlArgs
}

func CreateServer(cfg *ServerConfig, opts *core.Options) *Swager {
	swager := new(Swager)
	swager.cfg = cfg
  swager.opts = opts
	return swager
}

func (s *Swager) InitBlock(args *InitBlockArgs, reply *Reply) error {
	blockfac, ok := s.cfg.Blocks[args.Block]
	if !ok {
		return &BlockNotFoundError{args.Block}
	}

	block := blockfac()
	if err := block.Init(s.cfg.Client, s.cfg.Sub, s.opts); err != nil {
		return &BlockInitializationError{err, args.Block}
	}

  if s.opts.Debug {
    s.opts.Log <- fmt.Sprintf("[%s](%s) initalized", args.Block, args.Tag)
  }

	if err := block.Configure(args.Config); err != nil {
		return &BlockInitializationError{err, args.Block}
	}

  if s.opts.Debug {
    s.opts.Log <- fmt.Sprintf("[%s](%s) configured", args.Block, args.Tag)
  }

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

	reply.Success = true
	return nil
}

func (s *Swager) Control(args *ControlArgs, reply *Reply) error {
  switch args.Command {
  case PingServer:
    reply.Success = true
    return nil
  case RunServer:
    if s.opts.Debug {
      s.opts.Log <- "running initalized blocks"
    }
    for _, block := range s.initalized {
      go block.Run()
    }
    go s.cfg.Sub.Start()
    reply.Success = true
    return nil
  case ExitServer:
    s.cfg.Sub.Close()
    fallthrough
  default:
    s.cfg.Ctrl <- args
  }
	reply.Success = true
	return nil
}
