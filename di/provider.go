package di

import "reflect"

// Provider
type Provider interface {
	Identity() Identity
	Provide() (interface{}, error)
}

// dependencyProvider
type dependencyProvider interface {
	identity() identity
	parameters() parameterList
	provide(parameters ...reflect.Value) (reflect.Value, error)
}
