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
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=EventPayloadType
type EventPayloadType PayloadType

const (
	WorkspaceEvent       EventPayloadType = 0x80000000
	ModeEvent            EventPayloadType = 0x80000002
	WindowEvent          EventPayloadType = 0x80000003
	BarconfigUpdateEvent EventPayloadType = 0x80000004
	BindingEvent         EventPayloadType = 0x80000005
	ShutdownEvent        EventPayloadType = 0x80000006
	TickEvent            EventPayloadType = 0x80000007
	BarStateUpdateEvent  EventPayloadType = 0x80000014
	InputEvent           EventPayloadType = 0x80000015
)

var magic = [6]byte{'i', '3', '-', 'i', 'p', 'c'}

// ValidMagic tests whether the byte array represents
// the ipc payload magic string
func ValidMagic(test [6]byte) bool {
	return test == magic
}

type Header struct {
	Magic         [6]byte
	PayloadLength uint32
	PayloadType   PayloadType
}

func NewHeader(pt PayloadType, plen int) *Header {
	return &Header{Magic: magic, PayloadLength: uint32(plen), PayloadType: pt}
}

func (p EventPayloadType) eventName() string {
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
	case BarStateUpdateEvent:
		return "bar_status_update"
	case InputEvent:
		return "input"
	}

	return ""
}

func eventNames(ps []EventPayloadType) []string {
	s := make([]string, len(ps))
	for i, p := range ps {
		s[i] = p.eventName()
	}

	return s
}

func ToEventPayloadType(name string) (EventPayloadType, bool) {
	switch name {
	case "workspace":
		return WorkspaceEvent, true
	case "mode":
		return ModeEvent, true
	case "window":
		return WindowEvent, true
	case "barconfig_update":
		return BarconfigUpdateEvent, true
	case "binding":
		return BindingEvent, true
	case "shutdown":
		return ShutdownEvent, true
	case "tick":
		return TickEvent, true
	case "bar_status_update":
		return BarStateUpdateEvent, true
	case "input":
		return InputEvent, true
	}
	return EventPayloadType(0), false
}
