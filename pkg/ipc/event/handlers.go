package event

import "encoding/json"

func (s *subscription) handleWindowEvent(buf []byte) error {
	wc := new(WindowChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	go func(ch chan<- *WindowChange) {
		ch <- wc
	}(s.window)

	return nil
}

func (s *subscription) handleWorkspaceEvent(buf []byte) error {
	wc := new(WorkspaceChange)
	if err := json.Unmarshal(buf, wc); err != nil {
		return err
	}

	go func(ch chan<- *WorkspaceChange) {
		ch <- wc
	}(s.workspace)

	return nil
}

func (s *subscription) handleShutdownEvent(buf []byte) error {
	sc := new(ShutdownChange)
	if err := json.Unmarshal(buf, sc); err != nil {
		return err
	}

	go func(ch chan<- *ShutdownChange) {
		ch <- sc
	}(s.shutdown)

	return nil
}
