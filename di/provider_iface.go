package di

import (
	"fmt"
	"reflect"

	"github.com/defval/inject/v2/di/internal/reflection"
)

// createInterfaceProvider
func createInterfaceProvider(provider provider, as interface{}) *interfaceProvider {
	iface := reflection.InspectInterfacePtr(as)

	if !provider.resultKey().typ.Implements(iface.Type) {
		panicf("%s not implement %s", provider.resultKey(), iface.Type)
	}

	return &interfaceProvider{
		result: key{
			name: provider.resultKey().name,
			typ:  iface.Type,
		},
		implementation: provider,
	}
}

// interfaceProvider
type interfaceProvider struct {
	result         key
	implementation provider
}

func (i *interfaceProvider) resultKey() key {
	return i.result
}

func (i *interfaceProvider) parameters() providerParameterList {
	pl := providerParameterList{
		providerKey: i.resultKey(),
	}

	pl.add(parameterRequired{i.implementation.resultKey()})

	return pl
}

func (i *interfaceProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	return parameters[0], nil
}

func (i *interfaceProvider) Multiple() *multipleInterfaceProvider {
	return &multipleInterfaceProvider{result: i.result}
}

// multipleInterfaceProvider
type multipleInterfaceProvider struct {
	result key
}

func (m *multipleInterfaceProvider) resultKey() key {
	return m.result
}

func (m *multipleInterfaceProvider) parameters() providerParameterList {
	return providerParameterList{}
}

func (m *multipleInterfaceProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, fmt.Errorf("%s have sereral implementations", m.result.typ)
}
