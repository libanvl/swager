package event

import (
	"github.com/libanvl/swager/pkg/ipc/reply"
)

type WorkspaceChangeType string

const (
	InitWorkspace   WorkspaceChangeType = "init"
	EmptyWorkspace  WorkspaceChangeType = "empty"
	FocusWorkspace  WorkspaceChangeType = "focus"
	MoveWorkspace   WorkspaceChangeType = "move"
	RenameWorkspace WorkspaceChangeType = "rename"
	UrgentWorkspace WorkspaceChangeType = "urgent"
	ReloadWorkspace WorkspaceChangeType = "reload"
)

type WorkspaceChange struct {
	Change  WorkspaceChangeType
	Current reply.Node
	Old     reply.Node
}

type WindowChangeType string

const (
	NewWindow            WindowChangeType = "new"
	CloseWindow          WindowChangeType = "close"
	FocusWindow          WindowChangeType = "focus"
	TitleWindow          WindowChangeType = "title"
	FullscreenModeWindow WindowChangeType = "fullscreen_mode"
	MoveWindow           WindowChangeType = "move"
	FloatingWindow       WindowChangeType = "floating"
	UrgentWindow         WindowChangeType = "urgent"
	MarkWindow           WindowChangeType = "mark"
)

type WindowChange struct {
	Change    WindowChangeType `json:"change"`
	Container reply.Node       `json:"container"`
}

type ShutdownChangeType string

const (
	ExitShutdown ShutdownChangeType = "exit"
)

type ShutdownChange struct {
	Change ShutdownChangeType
}
