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

type WorkspaceChange struct {
	Change  WorkspaceChangeType
	Current Node
	Old     Node
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
	Container Node             `json:"container"`
}

type BindingModeChange struct {
	Change      string
	PangoMarkup bool `json:"pango_markup"`
}

type BindingChangeType string

const (
	RunBinding BindingChangeType = "run"
)

type InputType string

const (
	KeyboardInput InputType = "keyboard"
	MouseInput    InputType = "mouse"
)

type BindingChange struct {
	Change         BindingChangeType `json:"change"`
	Command        string            `json:"command"`
	EventStateMask []string          `json:"event_state_mask"`
	InputCode      int               `json:"input_code"`
	Symbol         string            `json:"symbol,omitempty"`
	InputType      InputType         `json:"input_type"`
}

type ShutdownChangeType string

const (
	ExitShutdown ShutdownChangeType = "exit"
)

type ShutdownChange struct {
	Change ShutdownChangeType
}

type Tick struct {
	First   bool
	Payload string
}
