package core

import (
	"github.com/libanvl/swager/pkg/ipc/reply"
)

func FindParent(root *reply.Node, childid int) *reply.Node {
	return root.FindChild(isParentOf(childid))
}

func isParentOf(childid int) func(n *reply.Node) bool {
	return func(nn *reply.Node) bool {
		for _, nnn := range nn.Nodes {
			if nnn.ID == childid {
				return true
			}
		}

		return false
	}
}
