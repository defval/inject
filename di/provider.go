package di

import "reflect"

// implementation
type provider interface {
	resultKey() key
	parameters() parameterList
	provide(parameters ...reflect.Value) (reflect.Value, error)
}
