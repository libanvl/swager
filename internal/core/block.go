package core

import (
  "go.i3wm.org/i3/v4"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=WinMgr
type WinMgr int

const (
  WinMgrSway WinMgr = iota
  WinMgrI3
)

type Block interface {
  Init(mgr WinMgr) error
  Configure(args []string) error
  Close()
}

type ChangeEventBlock interface {
  Block
  Event() []i3.EventType
  MatchChange(change string) bool
  OnEvent(evt interface{}) error
}
