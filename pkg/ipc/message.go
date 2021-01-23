package ipc

//go:generate go run golang.org/x/tools/cmd/stringer -type=PayloadType
type PayloadType uint32

const (
	RunCommandMessage      PayloadType = 0
	GetWorkspacesMessage   PayloadType = 1
	SubscribeMessage       PayloadType = 2
	GetOutputsMessage      PayloadType = 3
	GetTreeMessage         PayloadType = 4
	GetMarksMessage        PayloadType = 5
	GetBarConfigMessage    PayloadType = 6
	GetVersionMessage      PayloadType = 7
	GetBindingModesMessage PayloadType = 8
	GetConfigMessage       PayloadType = 9
	SendTickMessage        PayloadType = 10
	SyncMessage            PayloadType = 11
	GetBindingStateMessage PayloadType = 12
	GetInputsMessage       PayloadType = 100
	GetSeatsMessage        PayloadType = 101
	WorkspaceEvent         PayloadType = 0x80000000
	ModeEvent              PayloadType = 0x80000002
	WindowEvent            PayloadType = 0x80000003
	BarconfigUpdateEvent   PayloadType = 0x80000004
	BindingEvent           PayloadType = 0x80000005
	ShutdownEvent          PayloadType = 0x80000006
	TickEvent              PayloadType = 0x80000007
	BarStatusUpdateEvent   PayloadType = 0x80000014
	InputEvent             PayloadType = 0x80000015
)

var magic = [6]byte{'i', '3', '-', 'i', 'p', 'c'}

func ValidMagic(test [6]byte) bool {
	return test == magic
}

type Header struct {
	Magic         [6]byte
	PayloadLength uint32
	PayloadType   PayloadType
}

func NewHeader(pt PayloadType, plen int) *Header {
	h := new(Header)
	h.Magic = magic
	h.PayloadLength = uint32(plen)
	h.PayloadType = pt
	return h
}

func (p PayloadType) eventName() string {
	switch p {
	case WorkspaceEvent:
		return "workspace"
	case ModeEvent:
		return "mode"
	case WindowEvent:
		return "window"
	case BarconfigUpdateEvent:
		return "barconfig_update"
	case BindingEvent:
		return "binding"
	case ShutdownEvent:
		return "shutdown"
	case TickEvent:
		return "tick"
	case BarStatusUpdateEvent:
		return "bar_status_update"
	case InputEvent:
		return "input"
	}

	return ""
}

func eventNames(ps []PayloadType) []string {
	s := make([]string, len(ps))
	for i, p := range ps {
		s[i] = p.eventName()
	}

	return s
}
