package di

import (
	"fmt"
	"reflect"
)

// Parameter
type Parameter struct {
	internalParameter
}

// parameterRequired
type parameter struct {
	key
	optional bool
	embed    bool
}

func (p parameter) resolve(c *Container) (reflect.Value, error) {
	value, err := p.key.resolve(c)
	if _, notFound := err.(ErrProviderNotFound); notFound && p.optional {
		// create empty instance of type
		return reflect.New(p.typ).Elem(), nil
	}

	if err != nil {
		return reflect.Value{}, err
	}

	return value, nil
}

// parameterList
type parameterList []parameter

// resolve loads all parameters presented in parameter list.
func (pl parameterList) resolve(c *Container) ([]reflect.Value, error) {
	var values []reflect.Value
	for _, p := range pl {
		pvalue, err := p.resolve(c)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", p.resultKey(), err)
		}

		values = append(values, pvalue)
	}

	return values, nil
}

// internalParameter
type internalParameter interface{}

var parameterInterface = reflect.TypeOf(new(internalParameter)).Elem()

// isEmbedParameter
func isEmbedParameter(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct && typ.Implements(parameterInterface)
}
