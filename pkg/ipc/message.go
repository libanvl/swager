package ipc

type PayloadType uint32

const (
  RunCommandPayload      PayloadType =   0
  GetWorkspacesPayload   PayloadType =   1
  SubscribePayload       PayloadType =   2
  GetOutputsPayload      PayloadType =   3
  GetTreePayload         PayloadType =   4
  GetMarksPayload        PayloadType =   5
  GetBarConfigPayload    PayloadType =   6
  GetVersionPayload      PayloadType =   7
  GetBindingModesPayload PayloadType =   8
  GetConfigPayload       PayloadType =   9
  SendTickPayload        PayloadType =  10
  SyncPayload            PayloadType =  11
  GetBindingStatePayload PayloadType =  12
  GetInputsPayload       PayloadType = 100
  GetSeatsPayload        PayloadType = 101
)

var magic = [6]byte {'i', '3', '-', 'i', 'p', 'c'}

type header struct {
  Magic [6]byte
  PayloadLength uint32
  PayloadType PayloadType
}

type message struct {
  header
  payload []byte
}

func newMessage(pt PayloadType, p []byte) *message {
  m := new(message)
  m.Magic = magic
  m.PayloadLength = uint32(len(p))
  m.PayloadType = pt
  m.payload = p
  return m
}
