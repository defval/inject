package graph

import (
	"reflect"

	"github.com/pkg/errors"
)

// NewInterfaceNode
func NewInterfaceNode(name string, node *ProviderNode, iface interface{}) (_ *InterfaceNode, err error) {
	if iface == nil {
		return nil, errors.Errorf("nil interface") // todo: improve message
	}

	typ := reflect.TypeOf(iface)

	if typ.Kind() != reflect.Ptr {
		return nil, errors.Errorf("interface type must be a pointer to interface")
	}

	typ = typ.Elem()

	if typ.Kind() != reflect.Interface {
		return nil, errors.Errorf("only interface supported") // todo: improve message
	}

	return &InterfaceNode{
		key: Key{
			Type: typ,
			Name: name,
		},
		node: node,
	}, nil
}

// InterfaceNode
type InterfaceNode struct {
	WithOut
	key      Key
	node     *ProviderNode
	multiple bool
}

func (n *InterfaceNode) Key() Key {
	return n.key
}

func (n *InterfaceNode) Arguments() (args []Key) {
	return append(args, n.node.Key())
}

func (n *InterfaceNode) ArgumentNodes() (args []Node) {
	return append(args, n.node)
}

func (n *InterfaceNode) Extract(target reflect.Value) (err error) {
	if n.multiple {
		return errors.Errorf("could not extract %s: you have several instances of this interface type, use WithName() to identify it", n.Key())
	}

	if !target.Type().Implements(n.key.Type) {
		return errors.Errorf("%s not implement %s", target.Type(), n.key.Type)
	}

	return n.node.Extract(target)
}
