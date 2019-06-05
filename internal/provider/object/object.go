package object

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/provider"
)

const defaultArgumentTag = "inject"

// IsValid
func IsValid(rawProvider interface{}) (err error) {
	if rawProvider == nil {
		return errors.Errorf("nil provider")
	}

	value := reflect.ValueOf(rawProvider)

	isStruct := value.Kind() == reflect.Struct
	isStructPtr := value.Kind() == reflect.Ptr && reflect.Indirect(value).Kind() == reflect.Struct

	if !isStruct && !isStructPtr {
		return errors.Errorf("object provider must be a struct or pointer to struct")
	}

	return nil
}

type Option func(provider *Provider)

// Tag
func Tag(tag string) Option {
	return func(provider *Provider) {
		provider.tag = tag
	}
}

// Exported
func Exported() Option {
	return func(provider *Provider) {
		provider.includeExported = true
	}
}

// New creates object provider.
func New(rawProvider interface{}, options ...Option) (*Provider, error) {
	if err := IsValid(rawProvider); err != nil {
		return nil, errors.WithStack(err)
	}

	value := reflect.ValueOf(rawProvider)

	p := &Provider{
		id:    value.Type().String(),
		value: value,
		tag:   defaultArgumentTag,
	}

	if value.Kind() == reflect.Struct {
		p.isValue = true
	}

	for _, opt := range options {
		opt(p)
	}

	return p, nil
}

// Provider
type Provider struct {
	id    string
	value reflect.Value

	// options
	tag             string
	includeExported bool

	// internal flag
	isValue bool
}

func (p *Provider) Provide(arguments []reflect.Value) (reflect.Value, error) {
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

func (p *Provider) ResultType() reflect.Type {
	return p.value.Type()
}

func (p *Provider) Arguments() (args []provider.Key) {
	var value = p.elem()

	for i := 0; i < value.NumField(); i++ {
		name, injectable := p.isFieldInjectable(i)
		if !injectable {
			continue
		}

		args = append(args, provider.Key{
			Type: value.Field(i).Type(),
			Name: name,
		})
	}

	return args
}

func (p *Provider) elem() (elem reflect.Value) {
	value := p.value

	if p.isValue {
		value = reflect.New(value.Type())
	}

	return value.Elem()
}

func (p *Provider) isFieldInjectable(fieldNum int) (name string, _ bool) {
	value := p.elem()

	name, exists := value.Type().Field(fieldNum).Tag.Lookup(p.tag)
	return name, value.Field(fieldNum).CanSet() && (exists || p.includeExported)
}
