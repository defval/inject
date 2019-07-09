package graph

import (
	"reflect"

	"github.com/pkg/errors"
)

// newInterfaceNode creates new interface node.
func newInterfaceNode(name string, node *ProviderNode, iface interface{}) (_ *InterfaceNode, err error) {
	if iface == nil {
		return nil, errors.Errorf("nil interface") // todo: improve message
	}

	typ := reflect.TypeOf(iface)

	if typ.Kind() != reflect.Ptr {
		return nil, errors.Errorf("As() argument must be a pointer to interface, like new(http.Handler), got %s", typ.Kind())
	}

	typ = typ.Elem()

	if typ.Kind() != reflect.Interface {
		return nil, errors.Errorf("As() argument must be a pointer to interface, like new(http.Handler), got %s", typ.Kind()) // todo: improve message
	}

	if !node.ResultType().Implements(typ) {
		return nil, errors.Errorf("%s interface not implemented", typ.String())
	}

	return &InterfaceNode{
		key: Key{
			Type: typ,
			Name: name,
		},
		node: node,
	}, nil
}

// InterfaceNode ...
type InterfaceNode struct {
	outTrait
	key      Key
	node     *ProviderNode
	multiple bool
}

// Key returns unique node identifier.
func (n *InterfaceNode) Key() Key {
	return n.key
}

// Arguments returns another node keys that included in this group node.
func (n *InterfaceNode) Arguments() (args []Key) {
	return append(args, n.node.Key())
}

// ArgumentNodes return another nodes that included in this group.
// todo: Arguments() and ArgumentNodes() is too similar
func (n *InterfaceNode) ArgumentNodes() (args []Node) {
	return append(args, n.node)
}

// Extract extracts node instance to target.
func (n *InterfaceNode) Extract(target reflect.Value) (err error) {
	if n.multiple {
		return errors.Errorf("could not extract %s: you have several instances of this interface type, use WithName() to identify it", n.Key())
	}

	return n.node.Extract(target)
}
