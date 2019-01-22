package ding

import (
	"fmt"
	"sync"
)

// New creates new container
func New(options ...Option) *Container {
	var container = &Container{}

	for _, opt := range options {
		opt.apply(container)
	}

	return container
}

// Container ...
type Container struct {
	graph graph
}

// Populate
func (c *Container) Error() error {
	return nil
}

// node a single node that composes the tree
type node struct {
	constructor interface{}
}

func (n *node) String() string {
	return fmt.Sprintf("%v", n.constructor)
}

// graph the Items graph
type graph struct {
	nodes []*node
	edges map[node][]*node
	lock  sync.RWMutex
}

// AddNode adds a node to the graph
func (g *graph) AddNode(n *node) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *graph) AddEdge(n1, n2 *node) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[node][]*node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *graph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}
