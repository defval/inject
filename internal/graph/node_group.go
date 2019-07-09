package graph

import (
	"reflect"

	"github.com/emicklei/dot"
	"github.com/pkg/errors"
)

// newGroupNode creates new group node.
func newGroupNode(iface *InterfaceNode) *GroupNode {
	return &GroupNode{
		key: Key{
			Type: reflect.SliceOf(iface.Key().Type),
		},
		in: make([]*ProviderNode, 0),
	}
}

// GroupNode is a group node.
type GroupNode struct {
	outTrait
	key Key

	in []*ProviderNode

	node *dot.Node
}

// Key returns node unique identifier.
func (n *GroupNode) Key() Key {
	return n.key
}

// Arguments returns another node keys that included in this group node.
func (n *GroupNode) Arguments() (args []Key) {
	for _, in := range n.in {
		args = append(args, in.Key())
	}

	return args
}

// ArgumentNodes return another nodes that included in this group.
// todo: Arguments() and ArgumentNodes() is too similar
func (n *GroupNode) ArgumentNodes() (args []Node) {
	for _, in := range n.in {
		args = append(args, in)
	}

	return args
}

// Add adds provider node to group.
func (n *GroupNode) Add(node *ProviderNode) (err error) {
	n.in = append(n.in, node)

	return nil
}

// Replace replaces provider node in group.
func (n *GroupNode) Replace(node *ProviderNode) (err error) {
	for i, in := range n.in {
		if node.Key() != in.Key() {
			continue
		}

		n.in[i] = node
	}

	return nil
}

// Extract extracts group instance into target.
func (n *GroupNode) Extract(target reflect.Value) (err error) {
	// todo: test case
	// if target.Kind() != reflect.Slice || target.Type().Elem().Kind() != reflect.Interface {
	// 	return errors.Errorf("target value for extracting must be a slice of interfaces, got %s", target.Kind())
	// }

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
