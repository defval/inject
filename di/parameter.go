package di

import (
	"fmt"
	"reflect"
)

// ParameterList
type ParameterList []providerKey

// Register
func (l ParameterList) Register(container *Container, dependant providerKey) {
	for _, key := range l {
		if !container.graph.NodeExists(key) {
			panicf("%s: dependency %s not exists in container", dependant, key)
		}

		container.graph.AddEdge(key, dependant)
	}
}

// Load loads parameter values from container.
func (l ParameterList) Load(c *Container) ([]reflect.Value, error) {
	var values []reflect.Value
	for _, key := range l {
		value, err := key.Load(c)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", key, err)
		}

		values = append(values, value)
	}

	return values, nil
}
