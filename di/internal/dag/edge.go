package dag

type directedEdgeList struct {
	outgoingEdges map[Node]*nodeList
	incomingEdges map[Node]*nodeList
}

func newDirectedEdgeList() *directedEdgeList {
	return &directedEdgeList{
		outgoingEdges: make(map[Node]*nodeList),
		incomingEdges: make(map[Node]*nodeList),
	}
}

func (l *directedEdgeList) Copy() *directedEdgeList {
	outgoingEdges := make(map[Node]*nodeList, len(l.outgoingEdges))
	for node, edges := range l.outgoingEdges {
		outgoingEdges[node] = edges.Copy()
	}

	incomingEdges := make(map[Node]*nodeList, len(l.incomingEdges))
	for node, edges := range l.incomingEdges {
		incomingEdges[node] = edges.Copy()
	}

	return &directedEdgeList{
		outgoingEdges: outgoingEdges,
		incomingEdges: incomingEdges,
	}
}

func (l *directedEdgeList) Count() int {
	return len(l.outgoingEdges)
}

func (l *directedEdgeList) HasOutgoingEdges(node Node) bool {
	_, ok := l.outgoingEdges[node]
	return ok
}

func (l *directedEdgeList) OutgoingEdgeCount(node Node) int {
	if list := l.outgoingNodeList(node, false); list != nil {
		return list.Count()
	}
	return 0
}

func (l *directedEdgeList) outgoingNodeList(node Node, create bool) *nodeList {
	if list, ok := l.outgoingEdges[node]; ok {
		return list
	}
	if create {
		list := newNodeList()
		l.outgoingEdges[node] = list
		return list
	}
	return nil
}

func (l *directedEdgeList) OutgoingEdges(node Node) []Node {
	if list := l.outgoingNodeList(node, false); list != nil {
		return list.Nodes()
	}
	return nil
}

func (l *directedEdgeList) HasIncomingEdges(node Node) bool {
	_, ok := l.incomingEdges[node]
	return ok
}

func (l *directedEdgeList) IncomingEdgeCount(node Node) int {
	if list := l.incomingNodeList(node, false); list != nil {
		return list.Count()
	}
	return 0
}

func (l *directedEdgeList) incomingNodeList(node Node, create bool) *nodeList {
	if list, ok := l.incomingEdges[node]; ok {
		return list
	}
	if create {
		list := newNodeList()
		l.incomingEdges[node] = list
		return list
	}
	return nil
}

func (l *directedEdgeList) IncomingEdges(node Node) []Node {
	if list := l.incomingNodeList(node, false); list != nil {
		return list.Nodes()
	}
	return nil
}

func (l *directedEdgeList) Add(from Node, to Node) {
	l.outgoingNodeList(from, true).Add(to)
	l.incomingNodeList(to, true).Add(from)
}

func (l *directedEdgeList) Remove(from Node, to Node) {
	if list := l.outgoingNodeList(from, false); list != nil {
		list.Remove(to)

		if list.Count() == 0 {
			delete(l.outgoingEdges, from)
		}
	}
	if list := l.incomingNodeList(to, false); list != nil {
		list.Remove(from)

		if list.Count() == 0 {
			delete(l.incomingEdges, to)
		}
	}
}

func (l *directedEdgeList) Exists(from Node, to Node) bool {
	if list := l.outgoingNodeList(from, false); list != nil {
		return list.Exists(to)
	}
	return false
}
