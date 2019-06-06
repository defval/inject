package graph

import (
	"fmt"
	"reflect"
)

// Key
type Key struct {
	Type reflect.Type
	Name string
}

// Value creates new value of key type.
func (k Key) Value() reflect.Value {
	return reflect.New(k.Type).Elem()
}

func (k Key) String() string {
	if k.Name == "" {
		return fmt.Sprintf("%s", k.Type)
	}

	return fmt.Sprintf("%s[%s]", k.Type, k.Name)
}

// ObjectProvider
type InstanceProvider interface {
	Provide(arguments []reflect.Value) (reflect.Value, error)
	ResultType() reflect.Type
	Arguments() (args []Key)
}

// Node
type Node interface {
	Key() Key
	Extract(target reflect.Value) (err error)
}
