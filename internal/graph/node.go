package graph

import (
	"reflect"
)

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

// Node
type Node interface {
	Arguments
	Key() Key
	Extract(target reflect.Value) (err error)
}
