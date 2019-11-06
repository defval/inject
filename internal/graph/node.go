package graph

import (
	"github.com/defval/inject/internal/Lifetime"
	"reflect"
)

// Node ...
type Node interface {
	Arguments
	ArgumentNodes() []Node
	Key() Key
	Extract(target reflect.Value) (err error)
	Out() []Key
	Of(k Key)
	Lifetime() Lifetime.ProviderLifetime
}

// Arguments ...
type Arguments interface {
	Arguments() (args []Key)
}

// InstanceProvider ...
type InstanceProvider interface {
	Arguments
	Provide(arguments []reflect.Value) (reflect.Value, error)
	ResultType() reflect.Type
}
