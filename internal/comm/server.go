package comm

import (
  "net"
  "net/rpc"

  "github.com/libanvl/swager/internal/core"
  "github.com/libanvl/swager/pkg/ipc"
)

type SwagerServer interface {
  Start()
  NotifyExit(chan<- bool)
}

type Swager struct {
  l            net.Listener
  blocks       core.BlockRegistry
  opts         *core.Options
  initalized   map[string]core.Block
  client       ipc.Client
  subscription ipc.Subscription
  exitch       chan<- bool
}

func CreateServer(b core.BlockRegistry, l net.Listener, c ipc.Client, s ipc.Subscription, o *core.Options) (SwagerServer, error) {
  swager := new(Swager)
  swager.l = l
  swager.blocks = b
  swager.opts = o
  swager.client = c
  swager.subscription = s

  if err := rpc.Register(swager); err != nil {
    return nil, err
  }

  return swager, nil
}

func (s *Swager) NotifyExit(exit chan<- bool) {
  // move this to options
  s.exitch = exit
}

func (s *Swager) Start() {
  // do this from the caller
  rpc.Accept(s.l)
}

func (s *Swager) InitBlock(args *InitBlockArgs, reply *Reply) error {
  blockfac, ok := s.blocks[args.Block]
  if !ok {
    return &BlockNotFoundError{args.Block}
  }

  block := blockfac()
  if err := block.Init(s.client, s.subscription, s.opts); err != nil {
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

func (s *Swager) ExitServer(restart bool, reply *Reply) error {
  if !restart {
    reply.Success = false
    return nil
  }

  if s.exitch != nil {
    s.exitch <- restart
  }

  return nil
}
