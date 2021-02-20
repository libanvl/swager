package ipc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
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
	s.client = client

	return s
}

type Cookie uint32

var EmptyCookie = Cookie(0)

type WorkspaceChangeHandler interface {
	WorkspaceChange(WorkspaceChange)
}

type BindingModeChangeHandler interface {
	BindingModeChange(BindingModeChange)
}

type WindowChangeHandler interface {
	WindowChange(WindowChange)
}

type BindingChangeHandler interface {
	BindingChange(BindingChange)
}

type ShutdownChangeHandler interface {
	ShutdownChange(ShutdownChange)
}

type TickHandler interface {
	Tick(Tick)
}

type Subscription struct {
	client         *Client
	errors         []chan<- error
	clientmx       sync.Mutex
	currcookie     uint32
	workspaces     map[Cookie]WorkspaceChangeHandler
	bindingmodes   map[Cookie]BindingModeChangeHandler
	windows        map[Cookie]WindowChangeHandler
	bindings       map[Cookie]BindingChangeHandler
	shutdowns      map[Cookie]ShutdownChangeHandler
	ticks          map[Cookie]TickHandler
	workspacesmx   sync.Mutex
	bindingmodesmx sync.Mutex
	windowsmx      sync.Mutex
	bindingsmx     sync.Mutex
	shutdownsmx    sync.Mutex
	ticksmx        sync.Mutex
}

func (s *Subscription) Close() error {
	if s.client != nil {

		doLocked(&s.workspacesmx, func() {
			for k := range s.workspaces {
				delete(s.workspaces, k)
			}
		})

		doLocked(&s.bindingmodesmx, func() {
			for k := range s.bindingmodes {
				delete(s.bindingmodes, k)
			}
		})

		doLocked(&s.windowsmx, func() {
			for k := range s.windows {
				delete(s.windows, k)
			}
		})

		doLocked(&s.bindingsmx, func() {
			for k := range s.bindings {
				delete(s.bindings, k)
			}
		})

		doLocked(&s.shutdownsmx, func() {
			for k := range s.shutdowns {
				delete(s.shutdowns, k)
			}
		})

		doLocked(&s.ticksmx, func() {
			for k := range s.ticks {
				delete(s.ticks, k)
			}
		})

		err := s.client.Close()
		s.client = nil
		return err
	}

	return nil
}

// Errors returns the channel that subscription errors are yielded on.
// All errors from this channel are of type MonitoringError.
func (s *Subscription) Errors(ch chan<- error) {
	if s.errors == nil {
		s.errors = []chan<- error{ch}
	} else {
		s.errors = append(s.errors, ch)
	}
}

func (s *Subscription) RemoveHandler(c Cookie) {
	s.clientmx.Lock()
	defer s.clientmx.Unlock()

	if err := s.checkClient(); err != nil {
		return
	}

	delete(s.workspaces, c)
	delete(s.bindingmodes, c)
	delete(s.windows, c)
	delete(s.bindings, c)
	delete(s.shutdowns, c)
	delete(s.ticks, c)
}

// WorkspaceChanges returns the channel that WorkspaceChange events are yielded on.
func (s *Subscription) WorkspaceChanges(h WorkspaceChangeHandler) (Cookie, error) {
	//	s.clientmx.Lock()
	//	defer s.clientmx.Lock()

	if err := s.checkClient(); err != nil {
		return EmptyCookie, err
	}

	cookie := Cookie(atomic.AddUint32(&s.currcookie, 1))

	doLocked(&s.workspacesmx, func() {
		if s.workspaces == nil {
			s.workspaces = map[Cookie]WorkspaceChangeHandler{cookie: h}
			s.subscribeEvent(WorkspaceEvent)
		} else {
			s.workspaces[cookie] = h
		}
	})

	return cookie, nil
}

func (s *Subscription) BindingModeChanges(h BindingModeChangeHandler) (Cookie, error) {
	//	s.clientmx.Lock()
	//	defer s.clientmx.Unlock()

	if err := s.checkClient(); err != nil {
		return EmptyCookie, err
	}

	cookie := Cookie(atomic.AddUint32(&s.currcookie, 1))

	doLocked(&s.bindingmodesmx, func() {
		if s.bindingmodes == nil {
			s.bindingmodes = map[Cookie]BindingModeChangeHandler{cookie: h}
			s.subscribeEvent(ModeEvent)
		} else {
			s.bindingmodes[cookie] = h
		}
	})

	return cookie, nil
}

// WindowChanges returns the channel that WindowChange events are yielded on.
func (s *Subscription) WindowChanges(h WindowChangeHandler) (Cookie, error) {
	//	s.clientmx.Lock()
	//	defer s.clientmx.Unlock()

	if err := s.checkClient(); err != nil {
		return EmptyCookie, err
	}

	cookie := Cookie(atomic.AddUint32(&s.currcookie, 1))

	doLocked(&s.windowsmx, func() {
		if s.windows == nil {
			s.windows = map[Cookie]WindowChangeHandler{cookie: h}
			s.subscribeEvent(WindowEvent)
		} else {
			s.windows[cookie] = h
		}
	})

	return cookie, nil
}

func (s *Subscription) BindingChanges(h BindingChangeHandler) (Cookie, error) {
	//	s.clientmx.Lock()
	//	defer s.clientmx.Unlock()

	if err := s.checkClient(); err != nil {
		return EmptyCookie, err
	}

	cookie := Cookie(atomic.AddUint32(&s.currcookie, 1))

	doLocked(&s.bindingsmx, func() {
		if s.bindings == nil {
			s.bindings = map[Cookie]BindingChangeHandler{cookie: h}
			s.subscribeEvent(BindingEvent)
		} else {
			s.bindings[cookie] = h
		}
	})

	return cookie, nil
}

// ShutdownChanges returns the channel that ShutdownChange events are yielded on.
func (s *Subscription) ShutdownChanges(h ShutdownChangeHandler) (Cookie, error) {
	//	s.clientmx.Lock()
	//	defer s.clientmx.Unlock()

	if err := s.checkClient(); err != nil {
		return EmptyCookie, err
	}

	cookie := Cookie(atomic.AddUint32(&s.currcookie, 1))

	doLocked(&s.shutdownsmx, func() {
		if s.shutdowns == nil {
			s.shutdowns = map[Cookie]ShutdownChangeHandler{cookie: h}
			s.subscribeEvent(ShutdownEvent)
		} else {
			s.shutdowns[cookie] = h
		}
	})

	return cookie, nil
}

func (s *Subscription) Ticks(h TickHandler) (Cookie, error) {
	//	s.clientmx.Lock()
	//	defer s.clientmx.Unlock()

	if err := s.checkClient(); err != nil {
		return EmptyCookie, nil
	}

	cookie := Cookie(atomic.AddUint32(&s.currcookie, 1))

	doLocked(&s.ticksmx, func() {
		if s.ticks == nil {
			s.ticks = map[Cookie]TickHandler{cookie: h}
			s.subscribeEvent(TickEvent)
		} else {
			s.ticks[cookie] = h
		}
	})

	return cookie, nil
}

func (s *Subscription) Run() {
	s.clientmx.Lock()
	defer s.clientmx.Unlock()
	for s.client != nil {
		var h header
		if err := binary.Read(s.client, binary.LittleEndian, &h); err != nil {
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

		switch EventPayloadType(h.PayloadType) {
		case WorkspaceEvent:
			if err := s.handleWorkspace(buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("run s.handleWorkspace: %s", err)})
			}
			break
		case ModeEvent:
			if err := s.handleBindingMode(buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("run s.handleBindingMode: %s", err)})
			}
			break
		case WindowEvent:
			if err := s.handleWindow(buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("run s.handleWindow: %s", err)})
			}
			break
		case BindingEvent:
			if err := s.handleBinding(buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("run s.handleBinding: %s", err)})
			}
			break
		case ShutdownEvent:
			if err := s.handleShutdown(buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("run s.handleShutdown: %s", err)})
			}
			break
		case TickEvent:
			if err := s.handleTick(buf); err != nil {
				s.sendError(&MonitoringError{
					fmt.Errorf("run s.handleTick: %s", err)})
			}
			break
		default:
			s.sendError(&MonitoringError{
				errors.New("Unknown type")})
		}
	}
}

func (s *Subscription) subscribeEvent(event EventPayloadType) {
	//	s.clientmx.Lock()
	//	defer s.clientmx.Unlock()
	res, err := s.client.Subscribe(event)
	if err != nil {
		s.sendError(&MonitoringError{fmt.Errorf("subscribeEvent s.client.Subscribe: %s", err)})
	}
	if !res.Success {
		s.sendError(&MonitoringError{errors.New("sway error: could not subscribe to event")})
	}
}

func (s *Subscription) checkClient() error {
	if s.client == nil {
		return errors.New("Cannot add or remove handlers on a closed subscription")
	}

	return nil
}

func (s *Subscription) sendError(err error) {
	for _, e := range s.errors {
		go func(ch chan<- error) {
			ch <- err
		}(e)
	}
}

func doLocked(m *sync.Mutex, action func()) {
	m.Lock()
	defer m.Unlock()
	action()
}
