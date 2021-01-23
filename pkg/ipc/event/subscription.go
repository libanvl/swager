package event

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/libanvl/swager/pkg/ipc"
)

type Subscription interface {
	Errors() <-chan *MonitoringError
	Start() error
	Close() error
	Window() <-chan *WindowChange
	Workspace() <-chan *WorkspaceChange
	Shutdown() <-chan *ShutdownChange
}

func Subscribe() Subscription {
	s := new(subscription)
	s.evts = make([]ipc.PayloadType, 0)
	s.errors = make(chan *MonitoringError, 3)
	return s
}

type subscription struct {
	client    ipc.Client
	evts      []ipc.PayloadType
	errors    chan *MonitoringError
	startmx   sync.Mutex
	workspace chan *WorkspaceChange
	window    chan *WindowChange
	shutdown  chan *ShutdownChange
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

func (s *subscription) Window() <-chan *WindowChange {
	if !s.IsStarted() && s.window == nil {
		s.window = make(chan *WindowChange)
		s.evts = append(s.evts, ipc.WindowEvent)
	}
	return s.window
}

func (s *subscription) Workspace() <-chan *WorkspaceChange {
	if !s.IsStarted() && s.workspace == nil {
		s.workspace = make(chan *WorkspaceChange)
		s.evts = append(s.evts, ipc.WorkspaceEvent)
	}
	return s.workspace
}

func (s *subscription) Shutdown() <-chan *ShutdownChange {
	if !s.IsStarted() && s.shutdown == nil {
		s.shutdown = make(chan *ShutdownChange)
		s.evts = append(s.evts, ipc.ShutdownEvent)
	}
	return s.shutdown
}

func (s *subscription) Start() error {
	s.startmx.Lock()
	defer s.startmx.Unlock()

	if len(s.evts) < 1 {
		return errors.New("No  subscribed")
	}

	client, err := ipc.Connect()
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
		var h ipc.Header
		if err := binary.Read(s.client, binary.LittleEndian, &h); err != nil {
			s.errors <- &MonitoringError{err}
		}

		if !ipc.ValidMagic(h.Magic) {
			return nil
		}

		buf := make([]byte, int(h.PayloadLength))
		_, err := io.ReadFull(s.client, buf)
		if err != nil {
			s.errors <- &MonitoringError{err}
			continue
		}

		switch h.PayloadType {
		case ipc.WindowEvent:
			if err := s.handleWindowEvent(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		case ipc.WorkspaceEvent:
			if err := s.handleWorkspaceEvent(buf); err != nil {
				s.errors <- &MonitoringError{err}
			}
			break
		case ipc.ShutdownEvent:
			if err := s.handleShutdownEvent(buf); err != nil {
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
