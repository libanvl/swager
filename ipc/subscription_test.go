package ipc_test

import (
	"testing"

	"github.com/libanvl/swager/ipc"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionWithNilClient(t *testing.T) {
	tests := map[string]func(*ipc.Subscription) (any, error){
		"BindingChanges":   func(s *ipc.Subscription) (any, error) { return s.BindingChanges(nil) },
		"ModeChanges":      func(s *ipc.Subscription) (any, error) { return s.ModeChanges(nil) },
		"ShutdownChanges":  func(s *ipc.Subscription) (any, error) { return s.ShutdownChanges(nil) },
		"Ticks":            func(s *ipc.Subscription) (any, error) { return s.Ticks(nil) },
		"WindowChanges":    func(s *ipc.Subscription) (any, error) { return s.WindowChanges(nil) },
		"WorkspaceChanges": func(s *ipc.Subscription) (any, error) { return s.WorkspaceChanges(nil) },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sub := ipc.SubscribeCustom(nil)
			_, err := tc(sub)
			assert.NotNil(t, err)
		})
	}
}
