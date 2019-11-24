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

// register registers resultKey as dependency of providerKey in the container.
func (k key) register(c *Container, dependant key) {
	if !c.exists(k) {
		panicf("%s: dependency %s not exists in container", dependant, k)
	}

	c.registerDependency(k, dependant)
}

func (k key) load(c *Container) (reflect.Value, error) {
	provider, err := c.provider(k)
	if err != nil {
		return reflect.Value{}, err
	}

	values, err := provider.parameters().load(c)
	if err != nil {
		return reflect.Value{}, err
	}

	return provider.provide(values...)
}

// Extract extracts instance by resultKey from container into target.
func (k key) extractInto(c *Container, target interface{}) error {
	value, err := k.load(c)
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
