package graph

import (
	"reflect"

	"github.com/emicklei/dot"
	"github.com/pkg/errors"
)

// NewProviderNode
func NewProviderNode(name string, p InstanceProvider) (_ *ProviderNode) {
	node := &ProviderNode{
		key: Key{
			Type: p.ResultType(),
			Name: name,
		},
		InstanceProvider: p,
	}

	return node
}

// ProviderNode
type ProviderNode struct {
	InstanceProvider

	in       []Node
	key      Key
	instance reflect.Value
}

func (n *ProviderNode) Key() Key {
	return n.key
}

func (n *ProviderNode) DotNode(graph *dot.Graph) dot.Node {
	node := graph.Node(n.Key().String())
	node.Label(n.Key().String())
	node.Attr("color", "limegreen")
	node.Attr("fontname", "Helvetica")
	node.Attr("fontcolor", "white")
	node.Attr("style", "filled")
	node.Box()
	return node
}

func (n *ProviderNode) Extract(target reflect.Value) (err error) {
	if n.instance.IsValid() {
		target.Set(n.instance)
		return nil
	}

	var arguments []reflect.Value
	for _, argumentNode := range n.in {
		argumentTarget := reflect.New(argumentNode.Key().Type).Elem()

		if err = argumentNode.Extract(argumentTarget); err != nil {
			return errors.WithStack(err)
		}

		arguments = append(arguments, argumentTarget)
	}

	value, err := n.Provide(arguments)

	if err != nil {
		return errors.Wrapf(err, "%s", n.key)
	}

	if value.Kind() == reflect.Ptr && value.IsNil() {
		return errors.Errorf("%s: nil provided", n.Key())
	}

	n.instance = value
	target.Set(n.instance)

	return nil
}
