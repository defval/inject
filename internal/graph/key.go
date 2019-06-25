package graph

import (
	"fmt"
	"reflect"
)

// Key unique identifier of node graph.
type Key struct {
	Type reflect.Type
	Name string
}

// String is a string representation of key.
func (k Key) String() string {
	if k.Name == "" {
		return fmt.Sprintf("%s", k.Type)
	}

	return fmt.Sprintf("%s[%s]", k.Type, k.Name)
}

// outTrait
type outTrait struct {
	keys []Key
}

// Of
func (o *outTrait) Of(k Key) {
	o.keys = append(o.keys, k)
}

// outTrait
func (o *outTrait) Out() []Key {
	return o.keys
}
