package graph

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
	"github.com/pkg/errors"
)

// NewStorage
func NewStorage() *Storage {
	return &Storage{
		keys:  make([]Key, 0),
		nodes: make(map[Key]Node),
	}
}

// Storage
type Storage struct {
	keys  []Key
	nodes map[Key]Node
}

// Check
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
		s.keys = append(s.keys, groupNode.Key())
		s.nodes[groupNode.Key()] = groupNode
	}

	return s.nodes[groupNode.Key()].(*GroupNode), nil
}

// Get
func (s *Storage) Extract(name string, value reflect.Value) (err error) {
	k := Key{
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

				argumentNode.Of(provideNode.Key())
				provideNode.in = append(provideNode.in, argumentNode)
			}
		}
	}

	return s.detectCycles()
}

// Graph
func (s *Storage) Graph() *dot.Graph {

	root := dot.NewGraph(dot.Directed)

	for _, k := range s.keys {
		switch s.nodes[k].(type) {
		case *GroupNode, *InterfaceNode:
			if len(s.nodes[k].Out()) == 0 {
				continue
			}
		}

		var pkg string
		switch k.Type.Kind() {
		case reflect.Slice, reflect.Ptr:
			pkg = k.Type.Elem().PkgPath()
		default:
			pkg = k.Type.PkgPath()
		}

		subGraph := root.Subgraph(pkg, dot.ClusterOption{})
		subGraph.Attr("style", "rounded")
		subGraph.Attr("bgcolor", "#E8E8E8")
		subGraph.Attr("color", "lightgrey")
		subGraph.Attr("fontname", "COURIER")
		subGraph.Attr("fontcolor", "#46494C")
		graphNode := s.nodes[k].DotNode(subGraph)

		for _, in := range s.nodes[k].Arguments() {
			var argPkg string
			switch in.Type.Kind() {
			case reflect.Slice, reflect.Ptr:
				argPkg = in.Type.Elem().PkgPath()
			default:
				argPkg = in.Type.PkgPath()
			}

			subGraph := root.Subgraph(argPkg, dot.ClusterOption{})
			subGraph.Attr("style", "rounded")
			subGraph.Attr("color", "lightgrey")
			subGraph.Attr("bgcolor", "#E8E8E8")
			subGraph.Attr("fontname", "COURIER")
			subGraph.Attr("fontcolor", "#46494C")
			root.Edge(s.nodes[in].DotNode(subGraph), graphNode).Attr("color", "#949494")
		}
	}

	return root
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
