package di

import (
	"fmt"
	"reflect"

	"github.com/defval/inject/di/internal/reflection"
)

// createInterfaceProvider
func createInterfaceProvider(provider dependencyProvider, as interface{}) *interfaceProvider {
	iface := reflection.InspectInterfacePtr(as)

	if !provider.Result().Type.Implements(iface.Type) {
		panicf("%s not implement %s", provider.Result(), iface.Type)
	}

	// replace type to interface
	result := providerKey{
		Name: provider.Result().Name,
		Type: iface.Type,
	}

	return &interfaceProvider{
		result:   result,
		provider: provider,
	}
}

// interfaceProvider
type interfaceProvider struct {
	result   providerKey
	provider dependencyProvider
}

func (i *interfaceProvider) Result() providerKey {
	return i.result
}

func (i *interfaceProvider) Parameters() ParameterList {
	return append(ParameterList{}, i.provider.Result())
}

func (i *interfaceProvider) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	return parameters[0], nil
}

func (i *interfaceProvider) Multiple() *multipleInterfaceProvider {
	return &multipleInterfaceProvider{result: i.result}
}

// multipleInterfaceProvider
type multipleInterfaceProvider struct {
	result providerKey
}

func (m *multipleInterfaceProvider) Result() providerKey {
	return m.result
}

func (m *multipleInterfaceProvider) Parameters() ParameterList {
	return ParameterList{}
}

func (m *multipleInterfaceProvider) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, fmt.Errorf("%s have sereral implementations", m.result.Type)
}
