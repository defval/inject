package di

import (
	"fmt"
	"reflect"

	"github.com/defval/inject/di/internal/reflection"
)

// createInterfaceProvider
func createInterfaceProvider(provider dependencyProvider, as interface{}) *interfaceProvider {
	iface := reflection.InspectInterfacePtr(as)

	if !provider.Identity().typ.Implements(iface.Type) {
		panicf("%s not implement %s", provider.Identity(), iface.Type)
	}

	// replace type to interface
	result := identity{
		name: provider.Identity().name,
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

func (i *interfaceProvider) Identity() identity {
	return i.result
}

func (i *interfaceProvider) Parameters() parameterList {
	return append(parameterList{}, i.provider.Identity())
}

func (i *interfaceProvider) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	return parameters[0], nil
}

func (i *interfaceProvider) Multiple() *multipleInterfaceProvider {
	return &multipleInterfaceProvider{result: i.result}
}

// multipleInterfaceProvider
type multipleInterfaceProvider struct {
	result identity
}

func (m *multipleInterfaceProvider) Identity() identity {
	return m.result
}

func (m *multipleInterfaceProvider) Parameters() parameterList {
	return parameterList{}
}

func (m *multipleInterfaceProvider) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, fmt.Errorf("%s have sereral implementations", m.result.typ)
}
