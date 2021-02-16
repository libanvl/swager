package ipc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
)

// Subscribe to sway-ipc events.
// A single Subscription can listen for multiple event types,
// but each event payload is yielded on the corresponding
// typed channel.
func Subscribe() (*Subscription, error) {
	c, err := Connect()
	if err != nil {
		return nil, err
	}

	return SubscribeCustom(c, 10), nil
}

func SubscribeCustom(client *Client, chbufsize int) *Subscription {
	s := new(Subscription)
	s.errors = make(chan error, 3)
	s.client = client
	s.running = false
	s.chbufsize = chbufsize

	return s
}

type Subscription struct {
	client      *Client
	errors      chan error
	clientmx    sync.Mutex
	running     bool
	chbufsize   int
	workspace   chan *WorkspaceChange
	bindingmode chan *BindingModeChange
	window      chan *WindowChange
	binding     chan *BindingChange
	shutdown    chan *ShutdownChange
	tick        chan *Tick
}

func (s *Subscription) Close() error {
	if s.client != nil {
		if s.workspace != nil {
			close(s.workspace)
		}

		if s.bindingmode != nil {
			close(s.bindingmode)
		}

		if s.window != nil {
			close(s.window)
		}

		if s.binding != nil {
			close(s.binding)
		}

		if s.shutdown != nil {
			close(s.shutdown)
		}

		if s.tick != nil {
			close(s.tick)
		}

		err := s.client.Close()
		s.client = nil
		return err
	}

	return nil
}

// Errors returns the channel that subscription errors are yielded on.
// All errors from this channel are of type MonitoringError.
func (s *Subscription) Errors() <-chan error {
	return s.errors
}

// WorkspaceChanges returns the channel that WorkspaceChange events are yielded on.
func (s *Subscription) WorkspaceChanges() <-chan *WorkspaceChange {
	if !s.running && s.workspace == nil {
		s.workspace = make(chan *WorkspaceChange, s.chbufsize)
		s.subscribeEvent(WorkspaceEvent)
	}
	return s.workspace
}

func (s *Subscription) BindingModeChanges() <-chan *BindingModeChange {
	if !s.running && s.bindingmode == nil {
		s.bindingmode = make(chan *BindingModeChange, s.chbufsize)
		s.subscribeEvent(ModeEvent)
	}
	return s.bindingmode
}

// WindowChanges returns the channel that WindowChange events are yielded on.
func (s *Subscription) WindowChanges() <-chan *WindowChange {
	if !s.running && s.window == nil {
		s.window = make(chan *WindowChange, s.chbufsize)
		s.subscribeEvent(WindowEvent)
	}
	return s.window
}

func (s *Subscription) BindingChanges() <-chan *BindingChange {
	if !s.running && s.binding == nil {
		s.binding = make(chan *BindingChange, s.chbufsize)
		s.subscribeEvent(BindingEvent)
	}
	return s.binding
}

// ShutdownChanges returns the channel that ShutdownChange events are yielded on.
func (s *Subscription) ShutdownChanges() <-chan *ShutdownChange {
	if !s.running && s.shutdown == nil {
		s.shutdown = make(chan *ShutdownChange, 1)
		s.subscribeEvent(ShutdownEvent)
	}
	return s.shutdown
}

func (s *Subscription) Ticks() <-chan *Tick {
	if !s.running && s.tick == nil {
		s.tick = make(chan *Tick, s.chbufsize)
		s.subscribeEvent(TickEvent)
	}
	return s.tick
}

func (s *Subscription) Run() {
	s.clientmx.Lock()
	defer s.clientmx.Unlock()
	for s.client != nil {
		var h header
		if err := binary.Read(s.client, binary.LittleEndian, &h); err != nil {
			s.errors <- &MonitoringError{
				fmt.Errorf("run binary.Read: %s", err)}
		}

		if !validMagic(h.Magic) {
			continue
		}

		buf := make([]byte, int(h.PayloadLength))
		_, err := io.ReadFull(s.client, buf)
		if err != nil {
			s.errors <- &MonitoringError{
				fmt.Errorf("run io.ReadFull: %s", err)}
			continue
		}

		switch EventPayloadType(h.PayloadType) {
		case WorkspaceEvent:
			if err := s.handleWorkspace(buf); err != nil {
				s.errors <- &MonitoringError{
					fmt.Errorf("run s.handleWorkspace: %s", err)}
			}
			break
		case ModeEvent:
			if err := s.handleBindingMode(buf); err != nil {
				s.errors <- &MonitoringError{
					fmt.Errorf("run s.handleBindingMode: %s", err)}
			}
			break
		case WindowEvent:
			if err := s.handleWindow(buf); err != nil {
				s.errors <- &MonitoringError{
					fmt.Errorf("run s.handleWindow: %s", err)}
			}
			break
		case BindingEvent:
			if err := s.handleBinding(buf); err != nil {
				s.errors <- &MonitoringError{
					fmt.Errorf("run s.handleBinding: %s", err)}
			}
			break
		case ShutdownEvent:
			if err := s.handleShutdown(buf); err != nil {
				s.errors <- &MonitoringError{
					fmt.Errorf("run s.handleShutdown: %s", err)}
			}
			break
		case TickEvent:
			if err := s.handleTick(buf); err != nil {
				s.errors <- &MonitoringError{
					fmt.Errorf("run s.handleTick: %s", err)}
			}
			break
		default:
			s.errors <- &MonitoringError{
				errors.New("Unknown type")}
		}
	}
}

func (s *Subscription) subscribeEvent(event EventPayloadType) {
	s.clientmx.Lock()
	defer s.clientmx.Unlock()
	res, err := s.client.Subscribe(event)
	if err != nil {
		s.errors <- &MonitoringError{fmt.Errorf("subscribeEvent s.client.Subscribe: %s", err)}
	}
	if !res.Success {
		s.errors <- &MonitoringError{errors.New("sway error: could not subscribe to event")}
	}
}
