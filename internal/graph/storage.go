package graph

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// NewStorage creates new storage.
func NewStorage() *Storage {
	return &Storage{
		keys:  make([]Key, 0),
		nodes: make(map[Key]Node),
	}
}

// Storage is a type graph storage implementation.
type Storage struct {
	keys  []Key
	nodes map[Key]Node
}

// All returns all nodes.
func (s *Storage) All() (nodes []Node) {
	for _, k := range s.keys {
		nodes = append(nodes, s.nodes[k])
	}

	return nodes
}

// Add adds node to storage.
func (s *Storage) Add(node Node) (err error) {
	if n, ok := s.nodes[node.Key()]; ok {
		if ifaceNode, ok := n.(*InterfaceNode); ok {
			ifaceNode.multiple = true

			return nil
		}

		return errors.Errorf("%s: use named definition if you have several instances of the same type", node.Key())
	}

	s.keys = append(s.keys, node.Key())
	s.nodes[node.Key()] = node

	return nil
}

// Replace node in a storage.
func (s *Storage) Replace(node Node) (err error) {
	_, isProviderNode := node.(*ProviderNode)

	if _, ok := s.nodes[node.Key()]; !ok && !isProviderNode {
		return errors.Errorf("type %s not provided", node.Key())
	}

	s.nodes[node.Key()] = node

	return nil
}

// GroupNode returns or creates group node by interface.
func (s *Storage) GroupNode(iface *InterfaceNode) (_ *GroupNode) {
	groupNode := NewGroupNode(iface)

	if _, exists := s.nodes[groupNode.Key()]; !exists {
		s.keys = append(s.keys, groupNode.Key())
		s.nodes[groupNode.Key()] = groupNode
	}

	return s.nodes[groupNode.Key()].(*GroupNode)
}

// Extract extracts node instance into target.
func (s *Storage) Extract(name string, target reflect.Value) (err error) {
	k := Key{
		Type: target.Type(),
		Name: name,
	}

	node, exists := s.nodes[k]

	if !exists {
		return errors.Errorf("type %s not provided", k)
	}

	return node.Extract(target)
}

// Compile compiles the container.
func (s *Storage) Compile() (err error) {
	// link provide nodes
	for _, node := range s.nodes {
		if provideNode, ok := node.(*ProviderNode); ok {
			for _, k := range provideNode.Arguments() {
				argumentNode, exists := s.nodes[k]
				if !exists {
					return errors.Errorf("type %s not provided", k)
				}

				argumentNode.Of(provideNode.Key())
				provideNode.in = append(provideNode.in, argumentNode)
			}
		}
	}

	return s.detectCycles()
}

func (s *Storage) detectCycles() (err error) {
	visited := make(map[Key]visitStatus)

	for _, k := range s.keys {
		if err = s.visit(visited, s.nodes[k]); err != nil {
			return errors.Wrapf(err, "cycle detected")
		}
	}

	return nil
}

func (s *Storage) visit(visited map[Key]visitStatus, node Node) (err error) {
	if visited[node.Key()] == visitMarkPermanent {
		return
	}

	if visited[node.Key()] == visitMarkTemporary {
		return fmt.Errorf("%s", node.Key())
	}

	visited[node.Key()] = visitMarkTemporary

	for _, inKey := range node.Arguments() {
		if err = s.visit(visited, s.nodes[inKey]); err != nil {
			return errors.Wrapf(err, "%s", node.Key())
		}
	}

	visited[node.Key()] = visitMarkPermanent

	return nil
}

type visitStatus int

const (
	visitMarkTemporary visitStatus = iota + 1
	visitMarkPermanent
)
