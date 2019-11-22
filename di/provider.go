package di

import "reflect"

// dependencyProvider
type dependencyProvider interface {
	Identity() identity
	Parameters() parameterList
	Provide(parameters ...reflect.Value) (reflect.Value, error)
}
