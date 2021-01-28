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

  return SubscribeCustom(c), nil
}

func SubscribeCustom(client *Client) *Subscription {
	s := new(Subscription)
	s.evts = make([]EventPayloadType, 0)
	s.errors = make(chan error, 3)
	return s
}

type Subscription struct {
	client    *Client
	evts      []EventPayloadType
	errors    chan error
	startmx   sync.Mutex
	workspace chan *WorkspaceChange
	window    chan *WindowChange
	shutdown  chan *ShutdownChange
}

func (s *Subscription) Close() error {
	if s.client != nil {
		if s.workspace != nil {
			close(s.workspace)
		}

		if s.window != nil {
			close(s.window)
		}

		if s.shutdown != nil {
			close(s.shutdown)
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

func (s *Subscription) IsStarted() bool {
	s.startmx.Lock()
	defer s.startmx.Unlock()

	return s.client != nil
}

// WindowChanges returns the channel that WindowChange events are yielded on.
func (s *Subscription) WindowChanges() <-chan *WindowChange {
	if !s.IsStarted() && s.window == nil {
		s.window = make(chan *WindowChange)
		s.evts = append(s.evts, WindowEvent)
	}
	return s.window
}

// WorkspaceChanges returns the channel that WorkspaceChange events are yielded on.
func (s *Subscription) WorkspaceChanges() <-chan *WorkspaceChange {
	if !s.IsStarted() && s.workspace == nil {
		s.workspace = make(chan *WorkspaceChange)
		s.evts = append(s.evts, WorkspaceEvent)
	}
	return s.workspace
}

// ShutdownChanges returns the channel that ShutdownChange events are yielded on.
func (s *Subscription) ShutdownChanges() <-chan *ShutdownChange {
	if !s.IsStarted() && s.shutdown == nil {
		s.shutdown = make(chan *ShutdownChange)
		s.evts = append(s.evts, ShutdownEvent)
	}
	return s.shutdown
}

// Start starts the Subscription monitoring for subscribed events.
// Before calling Start, call the method for each event type this Subscription
// instance should monitor, at least once.
// After Start is called, calling an event method for the first time will return nil.
func (s *Subscription) Run() error {
	s.startmx.Lock()
	defer s.startmx.Unlock()

	if len(s.evts) < 1 {
		return errors.New("No events subscribed")
	}

	client, err := Connect()
	if err != nil {
		return fmt.Errorf("Failed to connect client: %v", err)
	}

	res, err := client.Subscribe(s.evts...)
	if err != nil {
		return fmt.Errorf("Failed to subscribe to : %v", err)
	}

	if !res.Success {
		return errors.New("sway: failed to subscribe to ")
	}

	s.client = client

	for s.client != nil {
		var h header
		if err := binary.Read(s.client, binary.LittleEndian, &h); err != nil {
			s.errors <- &MonitoringError{err}
		}

		if !validMagic(h.Magic) {
			return nil
		}

		buf := make([]byte, int(h.PayloadLength))
		_, err := io.ReadFull(s.client, buf)
		if err != nil {
			s.errors <- &MonitoringError{err}
			continue
		}

		switch EventPayloadType(h.PayloadType) {
		case WindowEvent:
			if err := s.handleWindow(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		case WorkspaceEvent:
			if err := s.handleWorkspace(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		case ShutdownEvent:
			if err := s.handleShutdown(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		default:
			s.errors <- &MonitoringError{
				errors.New("Unknown type")}
		}
	}

	return nil
}
