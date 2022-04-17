package core

import "github.com/libanvl/swager/pkg/ipc"

func Focused(ws []ipc.Workspace) *ipc.Workspace {
	for _, w := range ws {
		if w.Focused {
			return &w
		}
	}

	return nil
}

func Accept[T ~string](s T, allowed ...T) bool {
	for _, t := range allowed {
		if s == t {
			return true
		}
	}

	return false
}

func Deny[T ~string](s T, denied ...T) bool {
	for _, t := range denied {
		if s == t {
			return true
		}
	}

	return false
}
