package event

import (
  "github.com/libanvl/swager/pkg/ipc/reply"
)

type EventType uint32

const (
  WorkspaceEvent       EventType = 0x80000000
  ModeEvent            EventType = 0x80000002
  WindowEvent          EventType = 0x80000003
  BarconfigUpdateEvent EventType = 0x80000004
  BindingEvent         EventType = 0x80000005
  ShutdownEvent        EventType = 0x80000006
  TickEvent            EventType = 0x80000007
  BarStatusUpdateEvent EventType = 0x80000014
  InputEvent           EventType = 0x80000015
)

type WindowChangeType string

const (
  NewWindow             WindowChangeType = "new"
  CloseWindow           WindowChangeType = "close"
  FocusWindow           WindowChangeType = "focus"
  TitleWindow           WindowChangeType = "title"
  FullscreenModeWindow  WindowChangeType = "fullscreen_mode"
  MoveWindow            WindowChangeType = "move"
  FloatingWindow        WindowChangeType = "floating"
  UrgentWindow          WindowChangeType = "urgent"
  MarkWindow            WindowChangeType = "mark"
)

type WindowReply struct {
  Change    WindowChangeType `json:"change"`
  Container reply.Node       `json:"container"`
}
