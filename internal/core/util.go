package core

import (
	"github.com/libanvl/swager/pkg/ipc"
)

func FindParent(root *ipc.Node, childid int) *ipc.Node {
	return root.FindChild(isParentOf(childid))
}

func isParentOf(childid int) func(n *ipc.Node) bool {
	return func(nn *ipc.Node) bool {
		for _, nnn := range nn.Nodes {
			if nnn.ID == childid {
				return true
			}
		}

		return false
	}
}
