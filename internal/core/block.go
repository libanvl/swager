package core

import "io"

type Logger interface {
	Default(msg string)
	Defaultf(format string, args ...any)
	Info(msg string)
	Infof(format string, args ...any)
	Debug(msg string)
	Debugf(format string, args ...any)
}

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
