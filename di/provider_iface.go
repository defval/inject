package di

import (
	"fmt"
	"reflect"

	"github.com/defval/inject/v2/di/internal/reflection"
)

// createInterfaceProvider
func createInterfaceProvider(provider dependencyProvider, as interface{}) *interfaceProvider {
	iface := reflection.InspectInterfacePtr(as)

	if !provider.identity().typ.Implements(iface.Type) {
		panicf("%s not implement %s", provider.identity(), iface.Type)
	}

	// replace type to interface
	result := identity{
		name: provider.identity().name,
		typ:  iface.Type,
	}

	return &interfaceProvider{
		result:   result,
		provider: provider,
	}
}

// interfaceProvider
type interfaceProvider struct {
	result   identity
	provider dependencyProvider
}

func (i *interfaceProvider) identity() identity {
	return i.result
}

func (i *interfaceProvider) parameters() parameterList {
	return append(parameterList{}, parameter{
		identity: i.provider.identity(),
		optional: false,
	})
}

func (i *interfaceProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	return parameters[0], nil
}

func (i *interfaceProvider) Multiple() *multipleInterfaceProvider {
	return &multipleInterfaceProvider{result: i.result}
}

// multipleInterfaceProvider
type multipleInterfaceProvider struct {
	result identity
}

func (m *multipleInterfaceProvider) identity() identity {
	return m.result
}

func (m *multipleInterfaceProvider) parameters() parameterList {
	return parameterList{}
}

func (m *multipleInterfaceProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, fmt.Errorf("%s have sereral implementations", m.result.typ)
}
