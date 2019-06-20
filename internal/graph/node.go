package graph

import (
	"reflect"

	"github.com/emicklei/dot"
)

// Node
type Node interface {
	Arguments
	Key() Key
	DotNode(graph *dot.Graph) dot.Node
	Extract(target reflect.Value) (err error)
	Out() []Key
	Of(k Key)
}

// Arguments
type Arguments interface {
	Arguments() (args []Key)
}

// InstanceProvider
type InstanceProvider interface {
	Arguments
	Provide(arguments []reflect.Value) (reflect.Value, error)
	ResultType() reflect.Type
}
