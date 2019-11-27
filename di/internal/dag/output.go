package dag

import (
	"fmt"

	"github.com/tmc/dot"
)

// DOTGraph returns a textual representation of the graph in the DOT graph
// description language.
func (g *DirectedGraph) DOTGraph() string {
	graph := dot.NewGraph("G")
	graph.SetType(dot.DIGRAPH)

	itemsByNode := make(map[Node]*dot.Node)
	for _, node := range g.Nodes() {
		item := dot.NewNode(fmt.Sprintf("%v", node))
		itemsByNode[node] = item
		graph.AddNode(item)
	}

	for fromNode, fromItem := range itemsByNode {
		for _, toNode := range g.OutgoingEdges(fromNode) {
			if toItem, ok := itemsByNode[toNode]; ok {
				graph.AddEdge(dot.NewEdge(fromItem, toItem))
			}
		}
	}

	return graph.String()
}
