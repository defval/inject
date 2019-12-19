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
func (p parameter) ResolveProvider(c *Container) (provider, bool) {
	for _, pt := range providerLookupSequence {
		k := key{
			name: p.name,
			res:  p.res,
			typ:  pt,
		}
		provider, exists := c.provider(k)
		if !exists {
			continue
		}
		return provider, true
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
	values, err := pl.ResolveValues(c)
	if err != nil {
		return reflect.Value{}, err
	}
	value, err := provider.Provide(values...)
	if err != nil {
		return value, ErrParameterProvideFailed{k: provider.Key(), err: err}
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
