package ipc

import "encoding/json"

func (s *Subscription) handleWindow(buf []byte) error {
	wc := new(WindowChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	for _, h := range s.windows {
		go func(handler WindowChangeHandler) {
			handler.WindowChange(*wc)
		}(h)
	}

	return nil
}

func (s *Subscription) handleWorkspace(buf []byte) error {
	wc := new(WorkspaceChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	for _, h := range s.workspaces {
		go func(handler WorkspaceChangeHandler) {
			handler.WorkspaceChange(*wc)
		}(h)
	}

	return nil
}

func (s *Subscription) handleShutdown(buf []byte) error {
	sc := new(ShutdownChange)
	if err := json.Unmarshal(buf, sc); err != nil {
		return err
	}

	for _, h := range s.shutdowns {
		go func(handler ShutdownChangeHandler) {
			handler.ShutdownChange(*sc)
		}(h)
	}

	return nil
}

func (s *Subscription) handleBindingMode(buf []byte) error {
	bmc := new(BindingModeChange)
	if err := json.Unmarshal(buf, bmc); err != nil {
		return err
	}

	for _, h := range s.bindingmodes {
		go func(handler BindingModeChangeHandler) {
			handler.BindingModeChange(*bmc)
		}(h)
	}

	return nil
}

func (s *Subscription) handleBinding(buf []byte) error {
	bc := new(BindingChange)
	if err := json.Unmarshal(buf, bc); err != nil {
		return err
	}

	for _, h := range s.bindings {
		go func(handler BindingChangeHandler) {
			handler.BindingChange(*bc)
		}(h)
	}

	return nil
}

func (s *Subscription) handleTick(buf []byte) error {
	t := new(Tick)
	if err := json.Unmarshal(buf, t); err != nil {
		return err
	}

	for _, h := range s.ticks {
		go func(handler TickHandler) {
			handler.Tick(*t)
		}(h)
	}

	return nil
}
