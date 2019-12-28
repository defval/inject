package di

import (
	"reflect"
)

// Parameter
type Parameter struct {
	internalParameter
}

// parameterRequired
type parameter struct {
	name     string
	res      reflect.Type
	optional bool
	embed    bool
}

func (p parameter) String() string {
	return key{name: p.name, res: p.res}.String()
}

// ResolveProvider resolves parameter provider
func (p parameter) ResolveProvider(c *Container) (internalProvider, bool) {
	for _, pt := range providerLookupSequence {
		k := key{
			name: p.name,
			res:  p.res,
			typ:  pt,
		}
		if !c.graph.Exists(k) {
			continue
		}
		node := c.graph.Get(k)
		return node.Value.(internalProvider), true
	}
	return nil, false
}

func (p parameter) ResolveValue(c *Container) (reflect.Value, error) {
	provider, exists := p.ResolveProvider(c)
	if !exists && p.optional {
		return reflect.New(p.res).Elem(), nil
	}
	if !exists {
		return reflect.Value{}, ErrParameterProviderNotFound{param: p}
	}
	pl := provider.ParameterList()
	values, err := pl.Resolve(c)
	if err != nil {
		return reflect.Value{}, err
	}
	value, cleanup, err := provider.Provide(values...)
	if err != nil {
		return value, ErrParameterProvideFailed{k: provider.Key(), err: err}
	}
	if cleanup != nil {
		c.cleanups = append(c.cleanups, cleanup)
	}
	return value, nil
}

// isEmbedParameter
func isEmbedParameter(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct && typ.Implements(parameterInterface)
}

// internalParameter
type internalParameter interface {
	isDependencyInjectionParameter()
}

// parameterInterface
var parameterInterface = reflect.TypeOf(new(internalParameter)).Elem()
