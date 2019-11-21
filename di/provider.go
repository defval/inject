package di

import "reflect"

// dependencyProvider
type dependencyProvider interface {
	Result() providerKey
	Parameters() ParameterList
	Provide(parameters ...reflect.Value) (reflect.Value, error)
}
