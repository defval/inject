package graph

import (
	"github.com/defval/inject/internal/Scope"
	"github.com/pkg/errors"
	"reflect"
)

// NewProviderNode creates new provider node.
func NewProviderNode(name string, lifetime Scope.ProviderLifetime, p InstanceProvider) (_ *ProviderNode) {
	node := &ProviderNode{
		key: Key{
			Type: p.ResultType(),
			Name: name,
		},
		lifetime:         lifetime,
		InstanceProvider: p,
	}

	return node
}

// ProviderNode ...
type ProviderNode struct {
	outTrait
	InstanceProvider

	in       []Node
	key      Key
	lifetime Scope.ProviderLifetime
	instance reflect.Value
}

func (n *ProviderNode) Lifetime() Scope.ProviderLifetime {
	return n.lifetime
}

// Key returns unique node identifier.
func (n *ProviderNode) Key() Key {
	return n.key
}

// ArgumentNodes return another nodes that included in this group.
// todo: Arguments() and ArgumentNodes() is too similar
func (n *ProviderNode) ArgumentNodes() (args []Node) {
	for _, in := range n.in {
		args = append(args, in)
	}

	return args
}

// Extract extracts node into target.
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

	if n.lifetime == Scope.Singleton {
		n.instance = value
	}

	target.Set(value)

	return nil
}
