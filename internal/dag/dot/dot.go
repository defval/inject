package dot

import (
	"io"
	"text/template"

	"github.com/defval/inject/internal/dag"
)

// FromGraph
func FromGraph(graph *dag.Graph) *Graph {
	g := &Graph{
		ID: graph.ID(),
	}

	for _, node := range graph.Nodes() {
		g.Nodes = append(g.Nodes, &Node{
			ID: node.ID(),
		})

		for _, descendant := range node.Descendants() {
			g.Edges = append(g.Edges, &Edge{
				From: node.ID(),
				To:   descendant.ID(),
			})
		}
	}

	for _, sub := range graph.Subgraphs() {
		g.Subgraphs = append(g.Subgraphs, FromGraph(sub))
	}

	return g
}

// Graph
type Graph struct {
	ID        string
	Nodes     []*Node
	Edges     []*Edge
	Subgraphs []*Graph
}

// WriteTo
func (g *Graph) WriteTo(w io.Writer) (n int64, err error) {
	return 0, _mainGraphTemplate.Execute(w, g)
}

// Node
type Node struct {
	ID string
}

// Edge
type Edge struct {
	From string
	To   string
}

var _mainGraphTemplate = template.Must(template.New("DotGraph").Parse(`digraph {
{{- range .Subgraphs }}
	subgraph cluster_{{ .ID }} {
	{{- range .Nodes }}
		{{ .ID -}}
	{{ end }}
	
	{{- range .Edges }}
		{{ .From }} -> {{ .To }}
	{{ end -}}
	{{ end }}
	
	{{- range .Nodes }}
		{{ .ID -}}
	{{ end }}
	}

{{- range .Nodes }}
	{{ .ID -}}
{{ end }}

{{- range .Edges }}
	{{ .From }} -> {{ .To }}
{{ end -}}
`))
