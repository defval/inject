package graph

import (
	"reflect"

	"github.com/defval/inject/internal/provider"
)

// Node
type Node interface {
	Key() provider.Key
	Extract(target reflect.Value) (err error)
}
