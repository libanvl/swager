package core

type Block interface {
	Init(client Client, sub Sub, opts *Options, args ...string) error
	SetLogLevel(level LogLevel)
}

type Runner interface {
	Run()
}

type Closer interface {
	Close()
}

type Receiver interface {
	Receive(args []string) error
}

type JsonReporter interface {
	Report(args []string) (string, error)
}

type BlockRunnerCloser interface {
	Block
	Runner
	Closer
}
