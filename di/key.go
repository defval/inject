package di

import (
	"fmt"
	"reflect"
)

// providerKey is a key of represented instance in di.
type providerKey struct {
	Name string
	Type reflect.Type
}

// String represent key as string.
func (k providerKey) String() string {
	if k.Name == "" {
		return fmt.Sprintf("%s", k.Type)
	}

	return fmt.Sprintf("%s[%s]", k.Type, k.Name)
}

// Extract extracts instance by key from container into target.
func (k providerKey) Extract(c *Container, target interface{}) error {
	value, err := k.Load(c)
	if err != nil {
		return err
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetValue.Set(value)

	return nil
}

// Load loads instance by key from container.
func (k providerKey) Load(c *Container) (reflect.Value, error) {
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
