package ipc_test

import (
	"testing"

	"github.com/libanvl/swager/ipc"
	"github.com/stretchr/testify/assert"
)

func TestEventPayloadType(t *testing.T) {
	for _, ept := range []ipc.EventPayloadType{
		ipc.BarStateUpdateEvent,
		ipc.BarconfigUpdateEvent,
		ipc.InputEvent,
		ipc.ModeEvent,
		ipc.WorkspaceEvent,
		ipc.WindowEvent,
		ipc.TickEvent,
	} {
		assert.NotEmpty(t, ept.String())
	}
}

func TestFullScreenModeType(t *testing.T) {
	for _, fmt := range []ipc.FullscreenModeType{
		ipc.GlobalFullscreenMode,
		ipc.NoneFullscreenMode,
		ipc.WorkspaceFullscreenMode,
	} {
		assert.NotEmpty(t, fmt.String())
	}
}

func TestPayloadType(t *testing.T) {
	for _, pt := range []ipc.PayloadType{
		ipc.GetBarConfigMessage,
		ipc.SendTickMessage,
		ipc.RunCommandMessage,
		ipc.GetWorkspacesMessage,
		ipc.GetVersionMessage,
		ipc.GetTreeMessage,
		ipc.GetSeatsMessage,
		ipc.GetOutputsMessage,
		ipc.GetConfigMessage,
		ipc.GetSeatsMessage,
		ipc.SubscribeMessage,
	} {
		assert.NotEmpty(t, pt.String())
	}
}
