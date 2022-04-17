package core

import "io"

type BasicBlock struct {
	Client   Client
	Opts     *Options
	LogLevel LogLevel
	Log      Logger
}

type BlockInitializer interface {
	Init(client Client, sub Sub, opts *Options, log Logger, args ...string) error
	SetLogLevel(level LogLevel)
}

type Runner interface {
	Run()
}

type Receiver interface {
	Receive(args []string) error
}

type BlockRunnerCloser interface {
	BlockInitializer
	Runner
	io.Closer
}
