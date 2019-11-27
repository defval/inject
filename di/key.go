package di

import (
	"fmt"
	"reflect"
)

// resultKey is a id of represented instance in container
type key struct {
	name string
	typ  reflect.Type
}

// String represent resultKey as string.
func (k key) String() string {
	if k.name == "" {
		return fmt.Sprintf("%s", k.typ)
	}

	return fmt.Sprintf("%s[%s]", k.typ, k.name)
}

// resultKey
func (k key) resultKey() key {
	return k
}

func (k key) resolve(c *Container) (reflect.Value, error) {
	provider, exists := c.provider(k)
	if !exists {
		return reflect.Value{}, errProviderNotFound{k: k}
	}

	values, err := provider.parameters().resolve(c)
	if err != nil {
		return reflect.Value{}, err
	}

	return provider.provide(values...)
}

// Extract extracts instance by resultKey from container into target.
func (k key) extract(c *Container, target interface{}) error {
	value, err := k.resolve(c)
	if err != nil {
		return err
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetValue.Set(value)

	return nil
}

// errProviderNotFound
type errProviderNotFound struct {
	k key
}

func (e errProviderNotFound) Error() string {
	return fmt.Sprintf("type `%s` not exists in container", e.k)
}
