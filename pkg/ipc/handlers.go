package ipc

import "encoding/json"

func (s *Subscription) handleWindow(buf []byte) error {
	wc := new(WindowChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	go func(ch chan<- *WindowChange) {
		ch <- wc
	}(s.window)

	return nil
}

func (s *Subscription) handleWorkspace(buf []byte) error {
	wc := new(WorkspaceChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	go func(ch chan<- *WorkspaceChange) {
		ch <- wc
	}(s.workspace)

	return nil
}

func (s *Subscription) handleShutdown(buf []byte) error {
	sc := new(ShutdownChange)
	if err := json.Unmarshal(buf, sc); err != nil {
		return err
	}

	go func(ch chan<- *ShutdownChange) {
		ch <- sc
	}(s.shutdown)

	return nil
}

func (s *Subscription) handleBindingMode(buf []byte) error {
	bmc := new(BindingModeChange)
	if err := json.Unmarshal(buf, bmc); err != nil {
		return err
	}

	go func(ch chan<- *BindingModeChange) {
		ch <- bmc
	}(s.bindingmode)

	return nil
}

func (s *Subscription) handleBinding(buf []byte) error {
	bc := new(BindingChange)
	if err := json.Unmarshal(buf, bc); err != nil {
		return err
	}

	go func(ch chan<- *BindingChange) {
		ch <- bc
	}(s.binding)

	return nil
}

func (s *Subscription) handleTick(buf []byte) error {
	t := new(Tick)
	if err := json.Unmarshal(buf, t); err != nil {
		return err
	}

	go func(ch chan<- *Tick) {
		ch <- t
	}(s.tick)

	return nil
}
