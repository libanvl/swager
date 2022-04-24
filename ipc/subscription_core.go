package ipc

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type mapSyncPair[E EventArgs] struct {
	handlers map[Cookie]func(E)
	mx       sync.Mutex
}

func register[E EventArgs](s *Subscription, msp *mapSyncPair[E], ept EventPayloadType, h func(E)) (Cookie, error) {
	if err := s.ensureClient(); err != nil {
		return EmptyCookie, err
	}

	cookie := Cookie(atomic.AddUint32(&s.currcookie, 1))

	doLocked(&msp.mx, func() {
		if msp.handlers == nil {
			msp.handlers = map[Cookie]func(E){cookie: h}
			s.subscribeEvent(ept)
		} else {
			msp.handlers[cookie] = h
		}
	})

	return cookie, nil
}

func handle[E EventArgs](handlers map[Cookie]func(E), buf []byte) error {
	args := new(E)
	if err := json.Unmarshal(buf, args); err != nil {
		return err
	}

	for _, h := range handlers {
		go func(handler func(E)) {
			handler(*args)
		}(h)
	}

	return nil
}

func doLocked(m *sync.Mutex, action func()) {
	m.Lock()
	defer m.Unlock()
	action()
}

func (s *Subscription) subscribeEvent(event EventPayloadType) {
	if s.client != nil {
		res, err := s.client.Subscribe(event)
		if err != nil {
			s.sendError(&MonitoringError{fmt.Errorf("subscribeEvent s.client.Subscribe: %s", err)})
		}
		if !res.Success {
			s.sendError(&MonitoringError{errors.New("sway error: could not subscribe to event")})
		}
	}
}

func (s *Subscription) ensureClient() error {
	if s.client == nil {
		return errors.New("Cannot add handlers on a closed subscription")
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

func (msp *mapSyncPair[E]) reset() {
	doLocked(&msp.mx, func() {
		msp.handlers = nil
	})
}
