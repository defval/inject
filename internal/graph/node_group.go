package graph

import (
	"reflect"

	"github.com/emicklei/dot"
	"github.com/pkg/errors"
)

// NewGroupNode
func NewGroupNode(iface interface{}) (_ *GroupNode, err error) {
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

	return &GroupNode{
		key: Key{
			Type: reflect.SliceOf(typ),
		},
		in: make([]*ProviderNode, 0),
	}, nil
}

// GroupNode
type GroupNode struct {
	WithOut
	key Key

	in []*ProviderNode

	node *dot.Node
}

// Key
func (n *GroupNode) Key() Key {
	return n.key
}

// Arguments
func (n *GroupNode) Arguments() (args []Key) {
	for _, in := range n.in {
		args = append(args, in.Key())
	}

	return args
}

func (n *GroupNode) ArgumentNodes() (args []Node) {
	for _, in := range n.in {
		args = append(args, in)
	}

	return args
}

// Check
func (n *GroupNode) Add(node *ProviderNode) (err error) {
	if !node.ResultType().Implements(n.key.Type.Elem()) {
		return errors.Errorf("type %s not implement %s interface", node.ResultType(), n.key.Type.Elem())
	}

	n.in = append(n.in, node)

	return nil
}

func (n *GroupNode) Replace(node *ProviderNode) (err error) {
	for i, in := range n.in {
		if node.Key() != in.Key() {
			continue
		}

		n.in[i] = node
	}

	return nil
}

// Extract
func (n *GroupNode) Extract(target reflect.Value) (err error) {
	if target.Kind() != reflect.Slice || target.Type().Elem().Kind() != reflect.Interface {
		return errors.Errorf("target value for extracting must be a slice of interfaces, got %s", target.Kind())
	}

	var members []reflect.Value
	for _, node := range n.in {
		memberTarget := reflect.New(n.key.Type.Elem())

		if err = node.Extract(memberTarget.Elem()); err != nil {
			return errors.WithStack(err)
		}

		members = append(members, memberTarget.Elem())
	}

	target.Set(reflect.Append(target, members...))

	return nil
}
