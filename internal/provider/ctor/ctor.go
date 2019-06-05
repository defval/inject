package ctor

import (
	"reflect"
	"runtime"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/provider"
)

var errorInterface = reflect.TypeOf(new(error)).Elem()

// IsValid
func IsValid(provider interface{}) (err error) {
	if provider == nil {
		return errors.Errorf("nil provider")
	}

	value := reflect.ValueOf(provider)

	if value.Kind() != reflect.Func {
		return errors.Errorf("constructor provider must be a function, got %s", value.Type())
	}

	if value.Type().NumOut() == 0 {
		return errors.Errorf("%s must have at least one return value", runtime.FuncForPC(value.Pointer()).Name())
	}

	if value.Type().NumOut() > 2 {
		return errors.Errorf("%s: constructor function must have maximum two return values", runtime.FuncForPC(value.Pointer()).Name())
	}

	if value.Type().NumOut() == 2 && !value.Type().Out(1).Implements(errorInterface) {
		return errors.Errorf("%s: second argument of constructor must be error, got %s", runtime.FuncForPC(value.Pointer()).Name(), value.Type().Out(1))
	}

	return nil
}

// Option
type Option func(provider *Provider)

// New
func New(rawProvider interface{}, options ...Option) (*Provider, error) {
	if err := IsValid(rawProvider); err != nil {
		return nil, err
	}

	p := &Provider{
		ctor: reflect.ValueOf(rawProvider),
	}

	for _, opt := range options {
		opt(p)
	}

	return p, nil
}

// Provider
type Provider struct {
	ctor reflect.Value
}

func (p *Provider) String() string {
	return runtime.FuncForPC(p.ctor.Pointer()).Name()
}

func (p *Provider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	result := p.ctor.Call(arguments)

	if len(result) == 1 || result[1].IsNil() {
		return result[0], nil
	}

	return result[0], errors.WithStack(result[1].Interface().(error))
}

func (p *Provider) ResultType() reflect.Type {
	return p.ctor.Type().Out(0)
}

func (p *Provider) Arguments() (args []provider.Key) {
	pt := p.ctor.Type()

	for i := 0; i < pt.NumIn(); i++ {
		args = append(args, provider.Key{Type: pt.In(i)})
	}

	return args
}
