package dag

import "fmt"

// NewNode
func NewNode(id string, value interface{}) *Node {
	return &Node{
		id:          id,
		value:       value,
		ancestors:   []*Node{},
		descendants: []*Node{},
	}
}

// Node
type Node struct {
	id          string
	value       interface{}
	ancestors   []*Node
	descendants []*Node
}

// ID
func (n *Node) ID() string {
	return n.id
}

// Value
func (n *Node) Value() interface{} {
	return n.value
}

// Ancestors
func (n *Node) Ancestors() []*Node {
	return n.ancestors
}

// Descendants
func (n *Node) Descendants() []*Node {
	return n.descendants
}

// ConnectWith
func (n *Node) ConnectWith(descendant *Node) (err error) {
	for _, ancestor := range n.ancestors {
		if ancestor.id == descendant.id {
			return fmt.Errorf("could not add %s ancestor node to %s: already exists", descendant.id, n.id)
		}
	}

	descendant.ancestors = append(descendant.ancestors, n)
	n.descendants = append(n.descendants, descendant)

	return nil
}
