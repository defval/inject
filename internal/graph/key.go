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

// String
func (k Key) String() string {
	if k.Name == "" {
		return fmt.Sprintf("%s", k.Type)
	}

	return fmt.Sprintf("%s[%s]", k.Type, k.Name)
}

// WithOut
type WithOut struct {
	keys []Key
}

// Of
func (o *WithOut) Of(k Key) {
	o.keys = append(o.keys, k)
}

// WithOut
func (o *WithOut) Out() []Key {
	return o.keys
}
