package ipc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/libanvl/swager/pkg/ipc/event"
)

type Subscription interface {
	Errors() <-chan *MonitoringError
	Start() error
	Close() error
	Window() <-chan *event.WindowChange
	Workspace() <-chan *event.WorkspaceChange
	Shutdown() <-chan *event.ShutdownChange
}

func Subscribe() Subscription {
	s := new(subscription)
	s.evts = make([]PayloadType, 0)
	s.errors = make(chan *MonitoringError, 3)
	return s
}

type subscription struct {
	client    Client
	evts      []PayloadType
	errors    chan *MonitoringError
	startmx   sync.Mutex
	workspace chan *event.WorkspaceChange
	window    chan *event.WindowChange
	shutdown  chan *event.ShutdownChange
}

func (s *subscription) Close() error {
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

func (s *subscription) Errors() <-chan *MonitoringError {
	return s.errors
}

func (s *subscription) IsStarted() bool {
	s.startmx.Lock()
	defer s.startmx.Unlock()

	return s.client != nil
}

func (s *subscription) Window() <-chan *event.WindowChange {
	if !s.IsStarted() && s.window == nil {
		s.window = make(chan *event.WindowChange)
		s.evts = append(s.evts, WindowEvent)
	}
	return s.window
}

func (s *subscription) Workspace() <-chan *event.WorkspaceChange {
	if !s.IsStarted() && s.workspace == nil {
		s.workspace = make(chan *event.WorkspaceChange)
		s.evts = append(s.evts, WorkspaceEvent)
	}
	return s.workspace
}

func (s *subscription) Shutdown() <-chan *event.ShutdownChange {
	if !s.IsStarted() && s.shutdown == nil {
		s.shutdown = make(chan *event.ShutdownChange)
		s.evts = append(s.evts, ShutdownEvent)
	}
	return s.shutdown
}

func (s *subscription) Start() error {
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
		return fmt.Errorf("Failed to subscribe to events: %v", err)
	}

	if !res.Success {
		return errors.New("sway: failed to subscribe to events")
	}

	s.client = client

	for s.client != nil {
		var h header
		if err := binary.Read(s.client, binary.LittleEndian, &h); err != nil {
			s.errors <- &MonitoringError{err}
		}

		if h.Magic != magic {
			return nil
		}

		buf := make([]byte, int(h.PayloadLength))
		_, err := io.ReadFull(s.client, buf)
		if err != nil {
			s.errors <- &MonitoringError{err}
			continue
		}

		switch h.PayloadType {
		case WindowEvent:
			if err := s.handleWindowEvent(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		case WorkspaceEvent:
			if err := s.handleWorkspaceEvent(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		case ShutdownEvent:
			if err := s.handleShutdownEvent(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		default:
			s.errors <- &MonitoringError{
				errors.New("Unknown event type")}
		}
	}

	return nil
}

func (s *subscription) handleWindowEvent(buf []byte) error {
	wc := new(event.WindowChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	go func(ch chan<- *event.WindowChange) {
		ch <- wc
	}(s.window)

	return nil
}

func (s *subscription) handleWorkspaceEvent(buf []byte) error {
	wc := new(event.WorkspaceChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	go func(ch chan<- *event.WorkspaceChange) {
		ch <- wc
	}(s.workspace)

	return nil
}

func (s *subscription) handleShutdownEvent(buf []byte) error {
	sc := new(event.ShutdownChange)
	if err := json.Unmarshal(buf, sc); err != nil {
		return err
	}

	go func(ch chan<- *event.ShutdownChange) {
		ch <- sc
	}(s.shutdown)

	return nil
}
