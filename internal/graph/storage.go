package graph

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/provider"
)

// NewStorage
func NewStorage() *Storage {
	return &Storage{
		keys:  make([]provider.Key, 0),
		nodes: make(map[provider.Key]Node),
	}
}

// Storage
type Storage struct {
	keys  []provider.Key
	nodes map[provider.Key]Node
}

// Check
func (s *Storage) Add(node Node) (err error) {
	if _, ok := s.nodes[node.Key()]; ok {
		return errors.Errorf("%s: use named definition if you have several instances of the same type", node.Key())
	}

	s.keys = append(s.keys, node.Key())
	s.nodes[node.Key()] = node

	return nil
}

func (s *Storage) Replace(node Node) (err error) {
	_, isProviderNode := node.(*ProviderNode)

	if _, ok := s.nodes[node.Key()]; !ok && !isProviderNode {
		return errors.Errorf("type %s not provided", node.Key())
	}

	s.nodes[node.Key()] = node

	return nil
}

// GroupNode
func (s *Storage) GroupNode(iface interface{}) (_ *GroupNode, err error) {
	groupNode, err := NewGroupNode(iface)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if _, exists := s.nodes[groupNode.Key()]; !exists {
		s.nodes[groupNode.Key()] = groupNode
	}

	return s.nodes[groupNode.Key()].(*GroupNode), nil
}

// Get
func (s *Storage) Extract(name string, value reflect.Value) (err error) {
	k := provider.Key{
		Type: value.Type(),
		Name: name,
	}

	node, exists := s.nodes[k]

	if !exists {
		return errors.Errorf("type %s not provided", k)
	}

	return node.Extract(value)
}

// Compile
func (s *Storage) Compile() (err error) {
	// link provide nodes
	for _, node := range s.nodes {
		if provideNode, ok := node.(*ProviderNode); ok {
			for _, k := range provideNode.Arguments() {
				argumentNode, exists := s.nodes[k]
				if !exists {
					return errors.Errorf("type %s not provided", k)
				}

				provideNode.in = append(provideNode.in, argumentNode)
			}
			continue
		}
	}

	return s.detectCycles()
}

func (s *Storage) detectCycles() (err error) {
	visited := make(map[provider.Key]visitStatus)

	for _, k := range s.keys {
		if err = s.visit(visited, s.nodes[k]); err != nil {
			return errors.Wrapf(err, "cycle detected")
		}
	}

	return nil
}

func (s *Storage) visit(visited map[provider.Key]visitStatus, node Node) (err error) {
	if visited[node.Key()] == visitMarkPermanent {
		return
	}

	if visited[node.Key()] == visitMarkTemporary {
		return fmt.Errorf("%s", node.Key())
	}

	visited[node.Key()] = visitMarkTemporary

	switch concreteNode := node.(type) {
	case *ProviderNode:
		for _, in := range concreteNode.in {
			if err = s.visit(visited, in); err != nil {
				return errors.Wrapf(err, "%s", concreteNode.Key())
			}
		}
	case *InterfaceNode:
		for _, in := range concreteNode.node.in {
			if err = s.visit(visited, in); err != nil {
				return errors.Wrapf(err, "%s", concreteNode.Key())
			}
		}
	case *GroupNode:
		for _, in := range concreteNode.in {
			if err = s.visit(visited, in); err != nil {
				return errors.Wrapf(err, "%s", concreteNode.Key())
			}
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
