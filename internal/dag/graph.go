package dag

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyExists = errors.New("node already exists")
	ErrNotExists     = errors.New("node not exists")
)

// New creates new directed acyclic graph.
func New() *Graph {
	return &Graph{
		nodeIDs:   []string{},
		nodes:     map[string]*Node{},
		subgraphs: []*Graph{},
	}
}

// Graph is a directed acyclic graph container.
type Graph struct {
	id      string
	nodeIDs []string
	nodes   map[string]*Node
	parent  *Graph

	subgraphs []*Graph
}

// ID
func (g *Graph) ID() string {
	return g.id
}

// Nodes returns all nodes in graph.
func (g *Graph) Nodes() []*Node {
	var nodes []*Node

	for _, id := range g.nodeIDs {
		nodes = append(nodes, g.nodes[id])
	}

	return nodes
}

// Subgraph returns subgraph with id, if not exists it will be created.
func (g *Graph) Subgraph(id string) *Graph {
	if g.parent != nil {
		panic("subsubgraphs not allowed")
	}

	for _, subgraph := range g.subgraphs {
		if subgraph.id == id {
			return subgraph
		}
	}

	sub := New()
	sub.id = id
	sub.parent = g
	g.subgraphs = append(g.subgraphs, sub)

	return sub
}

// Subgraphs returns subgraphs.
func (g *Graph) Subgraphs() []*Graph {
	return g.subgraphs
}

// Add adds node to graph.
func (g *Graph) Add(node *Node) error {
	if _, exists := g.nodes[node.id]; exists {
		return fmt.Errorf("%s: %w", node, ErrAlreadyExists)
	}

	g.nodeIDs = append(g.nodeIDs, node.id)
	g.nodes[node.id] = node

	return nil
}

// Replace replaces node in graph.
func (g *Graph) Replace(id string, value interface{}) error {
	old, exists := g.nodes[id]
	if !exists {
		return fmt.Errorf("%s: %w", id, ErrNotExists)
	}

	old.value = value

	return nil
}

// Node returns node by id.
func (g *Graph) Node(id string) (*Node, error) {
	node, exists := g.nodes[id]
	if !exists {
		return nil, fmt.Errorf("%s: %w", id, ErrNotExists)
	}

	return node, nil
}

// Exists.
func (g *Graph) Exists(id string) bool {
	_, exists := g.nodes[id]
	return exists
}
