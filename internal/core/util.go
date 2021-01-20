package core

import (
  "log"

  "go.i3wm.org/i3/v4"
)

func FindParent(nid i3.NodeID) *i3.Node {
  tree, err := i3.GetTree()
  if err != nil {
    log.Panic("Failed getting tree")
  }

  return tree.Root.FindChild(isParentOf(nid))
}

func isParentOf(nid i3.NodeID) func (n *i3.Node) bool {
  return func (nn *i3.Node) bool {
    for _, nnn := range(nn.Nodes) {
      if nnn.ID == nid {
        return true
      }
    }

    return false
  }
}
