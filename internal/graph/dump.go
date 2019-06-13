package graph

// Graph
type Graph struct {
	Nodes []Key
	Edges []Edge
}

// Edge
type Edge struct {
	From Key
	To   Key
}
