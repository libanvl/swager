package reply

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
