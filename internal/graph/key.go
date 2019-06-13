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
