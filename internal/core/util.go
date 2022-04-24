package core

import "github.com/libanvl/swager/ipc"

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

func First[T any](ts []*T, walk func(*T) []*T, pred func(*T) bool) *T {
	for _, t := range ts {
		if pred(t) {
			return t
		}

		if tt := First(walk(t), walk, pred); tt != nil {
			return tt
		}
	}

	return nil
}
