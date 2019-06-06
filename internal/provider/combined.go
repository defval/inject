package provider

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/graph"
)

var combinedProviderInterface = reflect.Indirect(reflect.ValueOf(new(CombinedProviderInterface))).Type()

// CombinedProviderInterface
type CombinedProviderInterface interface {
	IsInjectProvider()
}

// IsCombinedProvider
func IsCombinedProvider(rawProvider interface{}) bool {
	return reflect.ValueOf(rawProvider).Type().Implements(combinedProviderInterface)
}

// NewConstructorProvider
func NewCombinedProvider(rawProvider interface{}, tag string, exported bool) (_ *CombinedProvider, err error) {
	value := reflect.ValueOf(rawProvider)

	if _, exists := value.Type().MethodByName("Provide"); !exists {
		return nil, errors.Errorf("combined provider must have Provide() method")
	}

	ctorMethod := value.MethodByName("Provide")

	objectProvider, err := NewObjectProvider(rawProvider, tag, exported)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ctorProvider, err := NewConstructorProvider(ctorMethod.Interface())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &CombinedProvider{objectProvider, ctorProvider}, nil
}

// ObjectProvider
type CombinedProvider struct {
	*ObjectProvider
	*ConstructorProvider
}

func (p *CombinedProvider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	_, _ = p.ObjectProvider.Provide(arguments)
	return p.ConstructorProvider.Provide([]reflect.Value{})
}

func (p *CombinedProvider) ResultType() reflect.Type {
	return p.ConstructorProvider.ResultType()
}

func (p *CombinedProvider) Arguments() (args []graph.Key) {
	return p.ObjectProvider.Arguments()
}
