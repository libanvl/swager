package core

type Block interface {
	Init(client Client, sub Sub, opts *Options) error
	Configure(args []string) error
	Run()
	Close()
}

type Receiver interface {
	Receive(args []string) error
}
