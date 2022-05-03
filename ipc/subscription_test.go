package ipc_test

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/libanvl/swager/ipc"
	"github.com/libanvl/swager/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionSubscribe(t *testing.T) {
	tmpsocket := t.TempDir() + "/uds"
	net.Listen("unix", tmpsocket)
	t.Setenv("SWAYSOCK", tmpsocket)
	sub, err := ipc.Subscribe()

	assert.Nil(t, err)
	assert.NotNil(t, sub)
}

func TestSubscriptionSubscribeNoSwaysock(t *testing.T) {
	t.Setenv("SWAYSOCK", "")
	_, err := ipc.Subscribe()
	assert.NotNil(t, err)

	os.Unsetenv("SWAYSOCK")
	_, err = ipc.Subscribe()
	assert.NotNil(t, err)
}

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

func TestRemoveHandler(t *testing.T) {
	tests := map[string]func(*ipc.Subscription) (ipc.Cookie, error){
		"BindingChanges":   func(s *ipc.Subscription) (ipc.Cookie, error) { return s.BindingChanges(nil) },
		"ModeChanges":      func(s *ipc.Subscription) (ipc.Cookie, error) { return s.ModeChanges(nil) },
		"ShutdownChanges":  func(s *ipc.Subscription) (ipc.Cookie, error) { return s.ShutdownChanges(nil) },
		"Ticks":            func(s *ipc.Subscription) (ipc.Cookie, error) { return s.Ticks(nil) },
		"WindowChanges":    func(s *ipc.Subscription) (ipc.Cookie, error) { return s.WindowChanges(nil) },
		"WorkspaceChanges": func(s *ipc.Subscription) (ipc.Cookie, error) { return s.WorkspaceChanges(nil) },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			conn := test.NewMockConnection(t)
			client := ipc.NewClient(conn, binary.LittleEndian)
			sub := ipc.SubscribeCustom(client)

			response, err := json.Marshal(ipc.Result{Success: true})
			require.NotNil(t, response)
			require.Nil(t, err)

			conn.PushPayloadForRead(uint32(ipc.SubscribeMessage), response, binary.LittleEndian)

			cookie, err := tc(sub)
			assert.Nil(t, err)
			assert.NotZero(t, cookie)
			sub.RemoveHandler(cookie)
		})
	}
}

func TestClose(t *testing.T) {
	tests := map[string]func(*ipc.Subscription) (ipc.Cookie, error){
		"BindingChanges":   func(s *ipc.Subscription) (ipc.Cookie, error) { return s.BindingChanges(nil) },
		"ModeChanges":      func(s *ipc.Subscription) (ipc.Cookie, error) { return s.ModeChanges(nil) },
		"ShutdownChanges":  func(s *ipc.Subscription) (ipc.Cookie, error) { return s.ShutdownChanges(nil) },
		"Ticks":            func(s *ipc.Subscription) (ipc.Cookie, error) { return s.Ticks(nil) },
		"WindowChanges":    func(s *ipc.Subscription) (ipc.Cookie, error) { return s.WindowChanges(nil) },
		"WorkspaceChanges": func(s *ipc.Subscription) (ipc.Cookie, error) { return s.WorkspaceChanges(nil) },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			conn := test.NewMockConnection(t)
			client := ipc.NewClient(conn, binary.LittleEndian)
			sub := ipc.SubscribeCustom(client)

			response, err := json.Marshal(ipc.Result{Success: true})
			require.NotNil(t, response)
			require.Nil(t, err)

			conn.PushPayloadForRead(uint32(ipc.SubscribeMessage), response, binary.LittleEndian)

			cookie, err := tc(sub)
			assert.Nil(t, err)
			assert.NotZero(t, cookie)
			err = sub.Close()
			assert.Nil(t, err)
			conn.AssertCalled("Close")
		})
	}
}

func TestHandlerRegistration(t *testing.T) {
	type action func(*ipc.Subscription) (ipc.Cookie, error)
	tests := map[string]struct {
		expected string
		action   action
	}{
		"BindingChanges": {
			"binding",
			func(s *ipc.Subscription) (ipc.Cookie, error) { return s.BindingChanges(nil) },
		},
		"ModeChanges": {
			"mode",
			func(s *ipc.Subscription) (ipc.Cookie, error) { return s.ModeChanges(nil) },
		},
		"ShutdownChanges": {
			"shutdown",
			func(s *ipc.Subscription) (ipc.Cookie, error) { return s.ShutdownChanges(nil) },
		},
		"Ticks": {
			"tick",
			func(s *ipc.Subscription) (ipc.Cookie, error) { return s.Ticks(nil) },
		},
		"WindowChanges": {
			"window",
			func(s *ipc.Subscription) (ipc.Cookie, error) { return s.WindowChanges(nil) },
		},
		"WorkspaceChanges": {
			"workspace",
			func(s *ipc.Subscription) (ipc.Cookie, error) { return s.WorkspaceChanges(nil) },
		},
	}

	for name, tc := range tests {
		for _, yo := range []binary.ByteOrder{binary.LittleEndian, binary.BigEndian} {
			t.Run(name+"_"+yo.String(), func(t *testing.T) {
				conn := test.NewMockConnection(t)
				client := ipc.NewClient(conn, yo)
				sub := ipc.SubscribeCustom(client)

				result := ipc.Result{
					Success: true,
				}

				result_json, err := json.Marshal(result)
				require.Nilf(t, err, "Got: %#v", err)

				conn.PushPayloadForRead(uint32(ipc.SubscribeMessage), result_json, yo)

				cookie, err := tc.action(sub)
				assert.Nil(t, err)
				assert.NotZero(t, cookie)

				assert.Contains(t, string(conn.WriteAt(1)), tc.expected)
			})
		}
	}
}

func TestRun(t *testing.T) {
	type action func(*ipc.Subscription, *assert.Assertions, any)
	tests := map[string]struct {
		payload ipc.EventPayloadType
		message any
		action  action
	}{
		"BindingChanges": {
			ipc.BindingEvent,
			ipc.BindingChange{
				Change:    ipc.RunBinding,
				Command:   "test command",
				InputType: ipc.KeyboardInput,
			},
			func(s *ipc.Subscription, a *assert.Assertions, x any) {
				s.BindingChanges(func(bc ipc.BindingChange) {
					a.EqualValues(x, bc)
					s.Close()
				})
			},
		},
		"ModeChanges": {
			ipc.ModeEvent,
			ipc.ModeChange{
				Change: "test mode",
			},
			func(s *ipc.Subscription, a *assert.Assertions, x any) {
				s.ModeChanges(func(mc ipc.ModeChange) {
					a.EqualValues(x, mc)
					s.Close()
				})
			},
		},
		"ShutdownChanges": {
			ipc.ShutdownEvent,
			ipc.ShutdownChange{
				Change: ipc.ExitShutdown,
			},
			func(s *ipc.Subscription, a *assert.Assertions, x any) {
				s.ShutdownChanges(func(sc ipc.ShutdownChange) {
					a.EqualValues(x, sc)
					s.Close()
				})
			},
		},
		"Ticks": {
			ipc.TickEvent,
			ipc.Tick{
				First:   true,
				Payload: "test payload",
			},
			func(s *ipc.Subscription, a *assert.Assertions, x any) {
				s.Ticks(func(t ipc.Tick) {
					a.EqualValues(x, t)
					s.Close()
				})
			},
		},
		"WindowChanges": {
			ipc.WindowEvent,
			ipc.WindowChange{
				Change:    ipc.FocusWindow,
				Container: ipc.Node{},
			},
			func(s *ipc.Subscription, a *assert.Assertions, x any) {
				s.WindowChanges(func(wc ipc.WindowChange) {
					a.EqualValues(x, wc)
					s.Close()
				})
			},
		},
		"WorkspaceChanges": {
			ipc.WorkspaceEvent,
			ipc.WorkspaceChange{
				Change:  ipc.InitWorkspace,
				Current: nil,
				Old:     nil,
			},
			func(s *ipc.Subscription, a *assert.Assertions, x any) {
				s.WorkspaceChanges(func(wc ipc.WorkspaceChange) {
					a.EqualValues(x, wc)
					s.Close()
				})
			},
		},
	}

	for name, tc := range tests {
		for _, yo := range []binary.ByteOrder{binary.LittleEndian, binary.BigEndian} {
			t.Run(name+"_"+yo.String(), func(t *testing.T) {
				a := assert.New(t)
				conn := test.NewMockConnection(t)
				client := ipc.NewClient(conn, yo)
				sub := ipc.SubscribeCustom(client)

				result := ipc.Result{
					Success: true,
				}

				result_json, err := json.Marshal(result)
				require.Nil(t, err)
				conn.PushPayloadForRead(uint32(ipc.SubscribeMessage), result_json, yo)

				tc.action(sub, a, tc.message)

				message_json, err := json.Marshal(tc.message)
				require.Nil(t, err)
				conn.PushPayloadForRead(uint32(tc.payload), message_json, yo)

				go sub.Run()
			})
		}
	}
}

func TestMonitoringError(t *testing.T) {
	x := ipc.NewMonitoringError(assert.AnError)
	assert.Contains(t, x.Error(), assert.AnError.Error())
	assert.Equal(t, assert.AnError, x.Unwrap())
}
