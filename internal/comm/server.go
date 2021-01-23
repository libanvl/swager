package comm

import (
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

	if err := block.Configure(args.Config); err != nil {
		return &BlockInitializationError{err, args.Block}
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
  default:
    s.cfg.Ctrl <- args
  }
	reply.Success = true
	return nil
}
