package dag

import (
	"fmt"

	"github.com/emicklei/dot"
)

// NodeVisualizer
type NodeVisualizer interface {
	Visualize(node *dot.Node)
	SubGraph() string
	IsPrimary() bool
}

// DOTGraph returns a textual representation of the graph in the DOT graph
// description language.
func (g *DirectedGraph) DOTGraph() *dot.Graph {
	root := dot.NewGraph(dot.Directed)
	root.Attr("splines", "ortho")

	subgraphs := make(map[string]*dot.Graph)
	itemsByNode := make(map[Node]dot.Node)
	for _, node := range g.Nodes() {
		visualizer := node.(NodeVisualizer)

		if !g.HasOutgoingEdges(node) && !visualizer.IsPrimary() {
			continue
		}

		name := fmt.Sprintf("%s", node)
		subgraph, ok := subgraphs[visualizer.SubGraph()]
		if !ok {
			subgraph = root.Subgraph(visualizer.SubGraph(), dot.ClusterOption{})
			subgraphs[visualizer.SubGraph()] = subgraph
			applySubGraphStyle(subgraph)
		}
		item := subgraph.Node(name)
		visualizer.Visualize(&item)
		itemsByNode[node] = item

	}

	for fromNode, fromItem := range itemsByNode {
		for _, toNode := range g.OutgoingEdges(fromNode) {
			if toItem, ok := itemsByNode[toNode]; ok {
				root.Edge(fromItem, toItem).Attr("color", "#949494")
			}
		}
	}

	return root
}

func applySubGraphStyle(graph *dot.Graph) {
	graph.Attr("label", "")
	graph.Attr("style", "rounded")
	graph.Attr("bgcolor", "#E8E8E8")
	graph.Attr("color", "lightgrey")
	graph.Attr("fontname", "COURIER")
	graph.Attr("fontcolor", "#46494C")
}
