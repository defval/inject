package inject

import (
	"fmt"
	"reflect"
)

// nodeStorage is ordered node map
type nodeStorage struct {
	keys  []reflect.Type
	nodes map[reflect.Type]*node
}

// add
func (s *nodeStorage) add(n *node) (err error) {
	if existingNode, ok := s.nodes[n.resultType]; ok {
		if existingNode.nodeType != nodeTypeGroup {
			return fmt.Errorf("%s already injected", n.resultType)
		}

		existingNode.args = append(existingNode.args, n.args...)

		return nil
	}

	s.nodes[n.resultType] = n
	s.keys = append(s.keys, n.resultType)

	return nil
}

func (s *nodeStorage) get(typ reflect.Type) (n *node, found bool) {
	n, found = s.nodes[typ]
	return n, found
}

func (s *nodeStorage) all() []*node {
	nodes := make([]*node, len(s.nodes))
	for i, key := range s.keys {
		nodes[i] = s.nodes[key]
	}
	return nodes
}
