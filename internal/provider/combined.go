package provider

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/graph"
)

var combinedProviderInterface = reflect.Indirect(reflect.ValueOf(new(CombinedProviderInterface))).Type()

// CombinedProviderInterface todo
type CombinedProviderInterface interface {
	IsInjectProvider()
}

// IsCombinedProvider todo
func IsCombinedProvider(rawProvider interface{}) bool {
	return reflect.ValueOf(rawProvider).Type().Implements(combinedProviderInterface)
}

// NewCombinedProvider todo
func NewCombinedProvider(rawProvider interface{}, tag string, exported bool) (_ *CombinedProvider, err error) {
	value := reflect.ValueOf(rawProvider)

	if _, exists := value.Type().MethodByName("Provide"); !exists {
		return nil, errors.Errorf("combined provider must have Provide() method")
	}

	ctorMethod := value.MethodByName("Provide")

	objectProvider, _ := NewObjectProvider(rawProvider, tag, exported) // todo: cases?

	ctorProvider, err := NewConstructorProvider(ctorMethod.Interface())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &CombinedProvider{objectProvider, ctorProvider}, nil
}

// CombinedProvider todo
type CombinedProvider struct {
	*ObjectProvider
	*ConstructorProvider
}

// Provide todo
func (p *CombinedProvider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	_, _ = p.ObjectProvider.Provide(arguments)
	return p.ConstructorProvider.Provide([]reflect.Value{})
}

// ResultType todo
func (p *CombinedProvider) ResultType() reflect.Type {
	return p.ConstructorProvider.ResultType()
}

// Arguments todo
func (p *CombinedProvider) Arguments() (args []graph.Key) {
	return p.ObjectProvider.Arguments()
}
