package di

import (
	"fmt"
	"reflect"
)

// Identity
type Identity struct {
	typ reflect.Type
}

// IdentityOf
func IdentityOf(v interface{}) Identity {
	return Identity{
		typ: reflect.TypeOf(v),
	}
}

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
		return reflect.Value{}, errTypeNotProvided{identity: k}
	}

	provider := c.providers[k]

	values, err := provider.parameters().Load(c)
	if err != nil {
		return reflect.Value{}, err
	}

	return provider.provide(values...)
}

// errTypeNotProvided
type errTypeNotProvided struct {
	identity identity
}

func (e errTypeNotProvided) Error() string {
	return fmt.Sprintf("type `%s` not exists in container", e.identity)
}
