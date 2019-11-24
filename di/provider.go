package di

import "reflect"

// implementation
type provider interface {
	resultKey() key
	parameters() providerParameterList
	provide(parameters ...reflect.Value) (reflect.Value, error)
}
