package di

import (
	"fmt"
	"reflect"
)

type parameter struct {
	identity identity
	optional bool
}

// parameterList
type parameterList []parameter

// Register
func (pl parameterList) Register(container *Container, dependant identity) {
	for _, param := range pl {
		if !container.graph.NodeExists(param.identity) {
			panicf("%s: dependency %s not exists in container", dependant, param.identity)
		}

		container.graph.AddEdge(param.identity, dependant)
	}
}

// Load loads parameter values from container.
func (pl parameterList) Load(c *Container) ([]reflect.Value, error) {
	var values []reflect.Value
	for _, parameter := range pl {
		value, err := parameter.identity.Load(c)
		if err != nil {
			// if type not provided and parameter is optional append zero value to values
			if _, ok := err.(errTypeNotProvided); ok && parameter.optional {
				values = append(values, reflect.Value{})
				continue
			}

			return nil, fmt.Errorf("%s: %s", parameter.identity, err)
		}

		values = append(values, value)
	}

	return values, nil
}
