package graph

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/provider"
)

// NewProviderNode
func NewProviderNode(name string, p provider.Provider) (_ *ProviderNode) {
	node := &ProviderNode{
		key: provider.Key{
			Type: p.ResultType(),
			Name: name,
		},
		Provider: p,
	}

	return node
}

// ProviderNode
type ProviderNode struct {
	provider.Provider

	in       []Node
	key      provider.Key
	instance reflect.Value
}

func (n *ProviderNode) Key() provider.Key {
	return n.key
}

func (n *ProviderNode) Extract(target reflect.Value) (err error) {
	if n.instance.IsValid() {
		target.Set(n.instance)
		return nil
	}

	var arguments []reflect.Value
	for _, argumentNode := range n.in {
		argumentTarget := argumentNode.Key().Value()

		if err = argumentNode.Extract(argumentTarget); err != nil {
			return errors.WithStack(err)
		}

		arguments = append(arguments, argumentTarget)
	}

	value, err := n.Provider.Provide(arguments)

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
