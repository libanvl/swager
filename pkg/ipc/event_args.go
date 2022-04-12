package ipc

type EventArgs interface {
	WorkspaceChange | ModeChange | WindowChange | BindingChange | ShutdownChange | Tick
}

type WorkspaceChange struct {
	Change  WorkspaceChangeType `json:"change"`
	Current *Node               `json:"current"`
	Old     *Node               `json:"old"`
}

type ModeChange struct {
	// Change is the custom name of the activated mode
	Change      string `json:"change"`
	PangoMarkup bool   `json:"pango_markup"`
}

type WindowChange struct {
	Change    WindowChangeType `json:"change"`
	Container Node             `json:"container"`
}

// Barconfig_Update is not implemented

type BindingChange struct {
	Change         BindingChangeType `json:"change"`
	Command        string            `json:"command"`
	EventStateMask []string          `json:"event_state_mask"`
	InputCode      int               `json:"input_code"`
	Symbol         *string           `json:"symbol"`
	InputType      InputType         `json:"input_type"`
}

type ShutdownChange struct {
	Change ShutdownChangeType `json:"change"`
}

type Tick struct {
	First   bool   `json:"first"`
	Payload string `json:"payload"`
}

// BarState_Update is not implemented
// Input is not implemented
