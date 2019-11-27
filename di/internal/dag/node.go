package dag

// Node represents a graph node.
type Node = interface{}

type nodeList struct {
	nodes []Node
	set   map[Node]bool
}

func newNodeList() *nodeList {
	return &nodeList{
		nodes: make([]Node, 0),
		set:   make(map[Node]bool),
	}
}

func (l *nodeList) Copy() *nodeList {
	nodes := make([]Node, len(l.nodes))
	copy(nodes, l.nodes)

	set := make(map[Node]bool, len(nodes))
	for _, node := range nodes {
		set[node] = true
	}

	return &nodeList{
		nodes: nodes,
		set:   set,
	}
}

func (l *nodeList) Nodes() []Node {
	return l.nodes
}

func (l *nodeList) Count() int {
	return len(l.nodes)
}

func (l *nodeList) Exists(node Node) bool {
	_, ok := l.set[node]
	return ok
}

func (l *nodeList) Add(nodes ...Node) {
	for _, node := range nodes {
		if l.Exists(node) {
			continue
		}

		l.nodes = append(l.nodes, node)
		l.set[node] = true
	}
}

func (l *nodeList) Remove(nodes ...Node) {
	for i := len(l.nodes) - 1; i >= 0; i-- {
		for j, node := range nodes {
			if l.nodes[i] == node {
				copy(l.nodes[i:], l.nodes[i+1:])
				l.nodes[len(l.nodes)-1] = nil
				l.nodes = l.nodes[:len(l.nodes)-1]

				delete(l.set, node)

				copy(nodes[j:], nodes[j+1:])
				nodes[len(nodes)-1] = nil
				nodes = nodes[:len(nodes)-1]

				break
			}
		}
	}
}
