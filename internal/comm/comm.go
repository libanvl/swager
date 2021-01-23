package comm

import "math"

type SwagerMethod string

const (
	InitBlock SwagerMethod = "Swager.InitBlock"
	SendToTag SwagerMethod = "Swager.SendToTag"
	Control   SwagerMethod = "Swager.Control"
)

func (sm SwagerMethod) String() string {
	return string(sm)
}

type ServerControl int8

const (
  PingServer ServerControl = 0
	ExitServer ServerControl = math.MaxInt8
)

type InitBlockArgs struct {
	Tag    string
	Block  string
	Config []string
}

type SendToTagArgs struct {
	Tag  string
	Args []string
}

type ControlArgs struct {
	Command ServerControl
	Args    []string
}

type Reply struct {
	Success bool
}
