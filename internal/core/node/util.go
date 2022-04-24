package node

import "github.com/libanvl/swager/ipc"

const NodeNameScratch = "__i3_scratch"

type NodePredicate func(*ipc.Node) bool

func First(root *ipc.Node, predicate NodePredicate) *ipc.Node {
	if root == nil {
		return nil
	}

	if predicate(root) {
		return root
	}

	for _, n := range root.Nodes {
		if nn := First(n, predicate); nn != nil {
			return nn
		}
	}

	for _, n := range root.FloatingNodes {
		if nn := First(n, predicate); nn != nil {
			return nn
		}
	}

	return nil
}

func Count(root *ipc.Node, pred NodePredicate) int {
	count := 0

	if root == nil {
		return count
	}

	if pred(root) {
		count++
	}

	for _, n := range root.Nodes {
		count += Count(n, pred)
	}

	for _, n := range root.FloatingNodes {
		count += Count(n, pred)
	}

	return count
}

func IsScratchpad(n *ipc.Node) bool {
	return MatchName(NodeNameScratch)(n)
}

func IsLeaf(n *ipc.Node) bool {
	return len(n.Nodes) == 0
}

func IsFocused(n *ipc.Node) bool {
	return n.Focused
}

func MatchName(name string) NodePredicate {
	return func(n *ipc.Node) bool {
		return n.Name == name
	}
}

func MatchType(t ipc.NodeType) NodePredicate {
	return func(n *ipc.Node) bool {
		return n.Type == t
	}
}

func MatchAnd(left NodePredicate, right NodePredicate) NodePredicate {
	return func(n *ipc.Node) bool {
		return left(n) && right(n)
	}
}

func MatchNot(pred NodePredicate) NodePredicate {
	return func(n *ipc.Node) bool {
		return !pred(n)
	}
}

func MatchParentOf(child *ipc.Node) NodePredicate {
	return func(n *ipc.Node) bool {
		for _, nn := range n.Nodes {
			if nn.ID == child.ID {
				return true
			}
		}

		return false
	}
}
