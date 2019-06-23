package dot

import (
	"reflect"

	"github.com/emicklei/dot"

	"github.com/defval/inject/internal/graph"
)

type Graph = dot.Graph

// NewGraph
func NewGraphFromStorage(storage *graph.Storage) *dot.Graph {
	root := dot.NewGraph(dot.Directed)

	root.Attr("splines", "ortho")

	for _, node := range storage.All() {
		switch node.(type) {
		case *graph.GroupNode, *graph.InterfaceNode:
			if len(node.Out()) == 0 {
				continue
			}
		}

		pkgGraph := addPkgSubgraph(root, node)
		graphNode := addDotNode(pkgGraph, node)

		for _, in := range node.ArgumentNodes() {
			pkgGraph := addPkgSubgraph(root, in)

			root.Edge(addDotNode(pkgGraph, in), graphNode).Attr("color", "#949494")
		}
	}

	return root
}

func addDotNode(root *dot.Graph, n graph.Node) dot.Node {
	result := root.Node(n.Key().String())
	result.Label(n.Key().String())

	result.Attr("fontname", "COURIER")
	result.Attr("style", "filled")
	result.Attr("fontcolor", "white")
	switch n.(type) {
	case *graph.ProviderNode:
		result.Attr("color", "#46494C")
		result.Box()
	case *graph.InterfaceNode:
		result.Attr("color", "#2589BD")
	case *graph.GroupNode:
		result.Attr("shape", "doubleoctagon")
		result.Attr("color", "#E54B4B")
	}

	return result
}

// addPkgSubgraph
func addPkgSubgraph(root *dot.Graph, node graph.Node) *dot.Graph {
	pkgGraph := root.Subgraph(packageString(node), dot.ClusterOption{})
	pkgGraph.Attr("label", "")
	pkgGraph.Attr("style", "rounded")
	pkgGraph.Attr("bgcolor", "#E8E8E8")
	pkgGraph.Attr("color", "lightgrey")
	pkgGraph.Attr("fontname", "COURIER")
	pkgGraph.Attr("fontcolor", "#46494C")

	return pkgGraph
}

func packageString(node graph.Node) string {
	var pkg string
	switch node.Key().Type.Kind() {
	case reflect.Slice, reflect.Ptr:
		pkg = node.Key().Type.Elem().PkgPath()
	default:
		pkg = node.Key().Type.PkgPath()
	}

	return pkg
}
