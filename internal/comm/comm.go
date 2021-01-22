package comm

type SwagerMethod string

const (
	InitBlock  SwagerMethod = "Swager.InitBlock"
	SendToTag  SwagerMethod = "Swager.SendToTag"
	ExitServer SwagerMethod = "Swager.ExitServer"
)

func (sm SwagerMethod) String() string {
	return string(sm)
}

type InitBlockArgs struct {
	Tag    string
	Block  string
	Config []string
}

type SendToTagArgs struct {
	Tag  string
	Args []string
}

type Reply struct {
	Success bool
}
