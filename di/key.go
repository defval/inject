package di

import (
	"fmt"
	"reflect"
)

// identity is a key of represented instance in di.
type identity struct {
	name string
	typ  reflect.Type
}

// String represent key as string.
func (k identity) String() string {
	if k.name == "" {
		return fmt.Sprintf("%s", k.typ)
	}

	return fmt.Sprintf("%s[%s]", k.typ, k.name)
}

// Extract extracts instance by key from container into target.
func (k identity) Extract(c *Container, target interface{}) error {
	value, err := k.Load(c)
	if err != nil {
		return err
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetValue.Set(value)

	return nil
}

// Load loads instance by key from container.
func (k identity) Load(c *Container) (reflect.Value, error) {
	if !c.graph.NodeExists(k) {
		return reflect.Value{}, fmt.Errorf("type `%s` not exists in container", k)
	}

	provider := c.providers[k]

	values, err := provider.Parameters().Load(c)
	if err != nil {
		return reflect.Value{}, err
	}

	return provider.Provide(values...)
}
