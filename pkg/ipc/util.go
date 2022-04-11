package ipc

// FindChild uses a breadth-first search. Nodes are searched
// before FloatingNodes. Returns the first
// child Node matching the predicate, or nil if none found.
func (n *Node) FindChild(predicate func(*Node) bool) *Node {
	if predicate(n) {
		return n
	}

	for _, c := range n.Nodes {
		if con := c.FindChild(predicate); con != nil {
			return con
		}
	}

	for _, c := range n.FloatingNodes {
		if con := c.FindChild(predicate); con != nil {
			return con
		}
	}

	return nil
}

const nodeNameScratch = "__i3_scratch"

// IsScratch returns a value indicating whether
// the Node is the scratchpad.
func (n Node) IsScratchpad() bool {
	return n.Name == nodeNameScratch
}
