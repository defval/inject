package di

import (
	"fmt"
	"reflect"
)

// createParameterBugProvider
func createParameterBugProvider(key key, parameters ParameterBag) provider {
	return createConstructor(key.String(), func() ParameterBag { return parameters })
}

// Parameters
type ParameterBag map[string]interface{}

// String
func (b ParameterBag) RequireString(key string) string {
	value, ok := b[key].(string)
	if !ok {
		panic(fmt.Sprintf("value for string key `%s` not found", key))
	}

	return value
}

var parameterBagType = reflect.TypeOf(ParameterBag{})
