package di

import "reflect"

// dependencyProvider
type dependencyProvider interface {
	Result() providerKey
	Parameters() parameterList
	Provide(parameters ...reflect.Value) (reflect.Value, error)
}
