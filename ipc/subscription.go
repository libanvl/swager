package ipc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
)

// Subscribe to sway-ipc events, creating a new Client.
// A single Subscription can listen for multiple event types.
// After all event handlers have been added, call Run to start
// listening for events.
func Subscribe() (*Subscription, error) {
	c, err := Connect()
	if err != nil {
		return nil, err
	}

	return SubscribeCustom(c), nil
}

// SubscribeCustom uses a custom Client instance to listen
// for events. The Client should be used only for this Subcription.
// Multiple events can be subscribed to.
// After all event handlers have been added, call Run to start
// listening for events.
func SubscribeCustom(client *Client) *Subscription {
	s := new(Subscription)
	s.client = client
	s.errors = make([]chan<- error, 0)

	return s
}

// Cookie represents a single registered event handler.
type Cookie uint32

var EmptyCookie = Cookie(0)

type Subscription struct {
	client     *Client
	errors     []chan<- error
	clientmx   sync.Mutex
	currcookie uint32
	workspaces mapSyncPair[WorkspaceChange]
	modes      mapSyncPair[ModeChange]
	windows    mapSyncPair[WindowChange]
	bindings   mapSyncPair[BindingChange]
	shutdowns  mapSyncPair[ShutdownChange]
	ticks      mapSyncPair[Tick]
}

// Errors returns the channel that subscription errors are yielded on.
// All errors from this channel are of type MonitoringError.
func (s *Subscription) Errors(ch chan<- error) {
	s.errors = append(s.errors, ch)
}

// WorkspaceChanges registers a new event handler.
func (s *Subscription) WorkspaceChanges(h func(WorkspaceChange)) (Cookie, error) {
	return register(s, &s.workspaces, WorkspaceEvent, h)
}

// ModeChanges registers a new event handler.
func (s *Subscription) ModeChanges(h func(ModeChange)) (Cookie, error) {
	return register(s, &s.modes, ModeEvent, h)
}

// WindowChanges registers a new event handler.
func (s *Subscription) WindowChanges(h func(WindowChange)) (Cookie, error) {
	return register(s, &s.windows, WindowEvent, h)
}

// BindingChanges registers a new event handler.
func (s *Subscription) BindingChanges(h func(BindingChange)) (Cookie, error) {
	return register(s, &s.bindings, BindingEvent, h)
}

// ShutdownChanges registers a new event handler.
func (s *Subscription) ShutdownChanges(h func(ShutdownChange)) (Cookie, error) {
	return register(s, &s.shutdowns, ShutdownEvent, h)
}

// Ticks registers a new event handler.
func (s *Subscription) Ticks(h func(Tick)) (Cookie, error) {
	return register(s, &s.ticks, TickEvent, h)
}

// RemoveHandler removes a registered event handler.
func (s *Subscription) RemoveHandler(c Cookie) {
	if err := s.ensureClient(); err != nil {
		return
	}

	delete(s.workspaces.handlers, c)
	delete(s.modes.handlers, c)
	delete(s.windows.handlers, c)
	delete(s.bindings.handlers, c)
	delete(s.shutdowns.handlers, c)
	delete(s.ticks.handlers, c)
}

// Run starts listening for events, calling the registered handlers
// as events come in.
func (s *Subscription) Run() {
	for {
		var h header

		if s.client == nil {
			break
		}

		s.clientmx.Lock()
		if s.client == nil {
			break
		}

		if err := binary.Read(s.client, s.client.yo, &h); err != nil {
			s.sendError(&MonitoringError{
				fmt.Errorf("run binary.Read: %s", err)})
			continue
		}

		if !validMagic(h.Magic) {
			continue
		}

		buf := make([]byte, int(h.PayloadLength))
		_, err := io.ReadFull(s.client, buf)
		if err != nil {
			s.sendError(&MonitoringError{
				fmt.Errorf("run io.ReadFull: %s", err)})
			continue
		}
		s.clientmx.Unlock()

		switch EventPayloadType(h.PayloadType) {
		case WorkspaceEvent:
			if err := handle(s.workspaces.handlers, buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("handle s.workspaces: %s", err)})
			}
			break
		case ModeEvent:
			if err := handle(s.modes.handlers, buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("handle s.modes: %s", err)})
			}
			break
		case WindowEvent:
			if err := handle(s.windows.handlers, buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("handle s.windows: %s", err)})
			}
			break
		case BindingEvent:
			if err := handle(s.bindings.handlers, buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("handle s.bindings: %s", err)})
			}
			break
		case ShutdownEvent:
			if err := handle(s.shutdowns.handlers, buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("handle s.shutdowns: %s", err)})
			}
			break
		case TickEvent:
			if err := handle(s.ticks.handlers, buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("handle s.ticks: %s", err)})
			}
			break
		default:
			s.sendError(&MonitoringError{
				errors.New("Unknown type")})
		}
	}
}

// Close removes all registered event handlers
// and closes the underlying Client.
func (s *Subscription) Close() error {
	if s.client != nil {
		s.workspaces.reset()
		s.modes.reset()
		s.windows.reset()
		s.bindings.reset()
		s.shutdowns.reset()
		s.ticks.reset()

		err := s.client.Close()
		s.client = nil
		return err
	}

	return nil
}
