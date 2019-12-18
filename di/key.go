package di

import (
	"fmt"
	"reflect"

	"github.com/emicklei/dot"
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

// IsPrimary
func (k key) IsPrimary() bool {
	if k.typ.Kind() == reflect.Slice {
		return false
	}
	if k.typ.Kind() == reflect.Interface {
		return false
	}
	return true
}

// Package
func (k key) SubGraph() string {
	var pkg string
	switch k.typ.Kind() {
	case reflect.Slice, reflect.Ptr:
		pkg = k.typ.Elem().PkgPath()
	default:
		pkg = k.typ.PkgPath()
	}

	return pkg
}

// Visualize
func (k key) Visualize(node *dot.Node) {
	node.Label(k.String())
	node.Attr("fontname", "COURIER")
	node.Attr("style", "filled")
	node.Attr("fontcolor", "white")
	if k.typ.Kind() == reflect.Slice {
		node.Attr("shape", "doubleoctagon")
		node.Attr("color", "#E54B4B")
		return
	}
	if k.typ.Kind() == reflect.Interface {
		node.Attr("color", "#2589BD")
		return
	}
	node.Attr("color", "#46494C")
	node.Box()
}

func (k key) resolve(c *Container) (reflect.Value, error) {
	provider, exists := c.provider(k)
	if !exists {
		return reflect.Value{}, ErrProviderNotFound{k: k}
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
		return fmt.Errorf("%s: %s", k, err)
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetValue.Set(value)

	return nil
}
