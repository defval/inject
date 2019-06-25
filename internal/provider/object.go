package provider

import (
	"reflect"

	"github.com/defval/inject/internal/graph"
)

// IsObjectProvider todo
func IsObjectProvider(rawProvider interface{}) bool {
	value := reflect.ValueOf(rawProvider)

	isStruct := value.Kind() == reflect.Struct
	isStructPtr := value.Kind() == reflect.Ptr && reflect.Indirect(value).Kind() == reflect.Struct

	return isStruct || isStructPtr
}

// NewObjectProvider creates object provider.
func NewObjectProvider(rawProvider interface{}, tag string, includeExported bool) (*ObjectProvider, error) {
	value := reflect.ValueOf(rawProvider)

	p := &ObjectProvider{
		id:              value.Type().String(),
		value:           value,
		tag:             tag,
		includeExported: includeExported,
	}

	if value.Kind() == reflect.Struct {
		p.isValue = true
	}

	return p, nil
}

// ObjectProvider todo
type ObjectProvider struct {
	id    string
	value reflect.Value

	// options
	tag             string
	includeExported bool

	// internal flag
	isValue bool
}

// Provide todo
func (p *ObjectProvider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	elem := p.elem()

	skip := 0
	for i := 0; i < elem.NumField(); i++ {
		_, injectable := p.isFieldInjectable(i)
		if injectable {
			elem.Field(i).Set(arguments[i-skip])
		} else {
			skip++
		}
	}

	if p.isValue {
		return elem, nil
	}

	return p.value, nil
}

// ResultType todo
func (p *ObjectProvider) ResultType() reflect.Type {
	return p.value.Type()
}

// Arguments todo
func (p *ObjectProvider) Arguments() (args []graph.Key) {
	var value = p.elem()

	for i := 0; i < value.NumField(); i++ {
		name, injectable := p.isFieldInjectable(i)
		if !injectable {
			continue
		}

		args = append(args, graph.Key{
			Type: value.Field(i).Type(),
			Name: name,
		})
	}

	return args
}

func (p *ObjectProvider) elem() (elem reflect.Value) {
	value := p.value

	if p.isValue {
		value = reflect.New(value.Type())
	}

	return value.Elem()
}

func (p *ObjectProvider) isFieldInjectable(fieldNum int) (name string, _ bool) {
	value := p.elem()

	if value.Type().Field(fieldNum).Type.String() == "inject.Provider" {
		return "", false
	}

	name, exists := value.Type().Field(fieldNum).Tag.Lookup(p.tag)
	return name, value.Field(fieldNum).CanSet() && (exists || p.includeExported)
}
