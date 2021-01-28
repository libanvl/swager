package ipc

//go:generate go run golang.org/x/tools/cmd/stringer -type=payloadType
type payloadType uint32
const (
	RunCommandMessage      payloadType = 0
	GetWorkspacesMessage   payloadType = 1
	SubscribeMessage       payloadType = 2
	GetOutputsMessage      payloadType = 3
	GetTreeMessage         payloadType = 4
	GetMarksMessage        payloadType = 5
	GetBarConfigMessage    payloadType = 6
	GetVersionMessage      payloadType = 7
	GetBindingModesMessage payloadType = 8
	GetConfigMessage       payloadType = 9
	SendTickMessage        payloadType = 10
	SyncMessage            payloadType = 11
	GetBindingStateMessage payloadType = 12
	GetInputsMessage       payloadType = 100
	GetSeatsMessage        payloadType = 101
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=EventPayloadType
type EventPayloadType payloadType
const (
	WorkspaceEvent         EventPayloadType = 0x80000000
	ModeEvent              EventPayloadType = 0x80000002
	WindowEvent            EventPayloadType = 0x80000003
	BarconfigUpdateEvent   EventPayloadType = 0x80000004
	BindingEvent           EventPayloadType = 0x80000005
	ShutdownEvent          EventPayloadType = 0x80000006
	TickEvent              EventPayloadType = 0x80000007
	BarStatusUpdateEvent   EventPayloadType = 0x80000014
	InputEvent             EventPayloadType = 0x80000015
)

var magic = [6]byte{'i', '3', '-', 'i', 'p', 'c'}

// validMagic tests whether the byte array represents
// the ipc payload magic string
func validMagic(test [6]byte) bool {
	return test == magic
}

type header struct {
	Magic         [6]byte
	PayloadLength uint32
	PayloadType   payloadType
}

func newHeader(pt payloadType, plen int) *header {
	h := new(header)
	h.Magic = magic
	h.PayloadLength = uint32(plen)
	h.PayloadType = pt
	return h
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
	case BarStatusUpdateEvent:
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
