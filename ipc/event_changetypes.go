package ipc

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

type BindingChangeType string

const (
	RunBinding BindingChangeType = "run"
)

type InputType string

const (
	KeyboardInput InputType = "keyboard"
	MouseInput    InputType = "mouse"
)

type ShutdownChangeType string

const (
	ExitShutdown ShutdownChangeType = "exit"
)
