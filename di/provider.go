package di

import "reflect"

// provider
type provider interface {
	resultKey() key
	parameters() parameterList
	provide(parameters ...reflect.Value) (reflect.Value, error)
}

// cleanup
type cleanup interface {
	cleanup()
}
