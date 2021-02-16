package comm

import "github.com/libanvl/swager/pkg/ipc"

type subdemux struct {
	sub        *ipc.Subscription
	workspaces []chan *ipc.WorkspaceChange
	windows    []chan *ipc.WindowChange
	bindings   []chan *ipc.BindingChange
	shutdowns  []chan *ipc.ShutdownChange
	ticks      []chan *ipc.Tick
}

func SubDemux(sub *ipc.Subscription) *subdemux {
	return &subdemux{sub: sub}
}

func (s *subdemux) WorkspaceChanges() <-chan *ipc.WorkspaceChange {
	if s.workspaces == nil {
		s.workspaces = make([]chan *ipc.WorkspaceChange, 0)
		ch := s.sub.WorkspaceChanges()
		go func(ch <-chan *ipc.WorkspaceChange) {
			for val := range ch {
				for _, wc := range s.workspaces {
					wc <- val
				}
			}
		}(ch)
	}

	ch := make(chan *ipc.WorkspaceChange)
	s.workspaces = append(s.workspaces, ch)
	return ch
}

func (s *subdemux) WindowChanges() <-chan *ipc.WindowChange {
	if s.windows == nil {
		s.windows = make([]chan *ipc.WindowChange, 0)
		ch := s.sub.WindowChanges()
		go func(ch <-chan *ipc.WindowChange) {
			for val := range ch {
				for _, wc := range s.windows {
					wc <- val
				}
			}
		}(ch)
	}

	ch := make(chan *ipc.WindowChange)
	s.windows = append(s.windows, ch)
	return ch
}

func (s *subdemux) BindingChanges() <-chan *ipc.BindingChange {
	if s.bindings == nil {
		s.bindings = make([]chan *ipc.BindingChange, 0)
		ch := s.sub.BindingChanges()
		go func(ch <-chan *ipc.BindingChange) {
			for val := range ch {
				for _, wc := range s.bindings {
					wc <- val
				}
			}
		}(ch)
	}

	ch := make(chan *ipc.BindingChange)
	s.bindings = append(s.bindings, ch)
	return ch
}

func (s *subdemux) ShutdownChanges() <-chan *ipc.ShutdownChange {
	if s.shutdowns == nil {
		s.shutdowns = make([]chan *ipc.ShutdownChange, 0)
		ch := s.sub.ShutdownChanges()
		go func(ch <-chan *ipc.ShutdownChange) {
			for val := range ch {
				for _, wc := range s.shutdowns {
					wc <- val
				}
			}
		}(ch)
	}

	ch := make(chan *ipc.ShutdownChange)
	s.shutdowns = append(s.shutdowns, ch)
	return ch
}

func (s *subdemux) Ticks() <-chan *ipc.Tick {
	if s.ticks == nil {
		s.ticks = make([]chan *ipc.Tick, 0)
		ch := s.sub.Ticks()
		go func(ch <-chan *ipc.Tick) {
			for val := range ch {
				for _, wc := range s.ticks {
					wc <- val
				}
			}
		}(ch)
	}

	ch := make(chan *ipc.Tick)
	s.ticks = append(s.ticks, ch)
	return ch
}
