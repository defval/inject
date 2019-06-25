package provider

import (
	"reflect"
	"runtime"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/graph"
)

// IsConstructor todo
func IsConstructor(rawProvider interface{}) bool {
	return reflect.ValueOf(rawProvider).Kind() == reflect.Func
}

// NewConstructorProvider todo
func NewConstructorProvider(rawProvider interface{}) (*ConstructorProvider, error) {
	value := reflect.ValueOf(rawProvider)

	if value.Type().NumOut() == 0 {
		return nil, errors.Errorf("%s constructor function must have at least one return value", runtime.FuncForPC(value.Pointer()).Name())
	}

	if value.Type().NumOut() > 2 {
		return nil, errors.Errorf("%s constructor function must have maximum two return values", runtime.FuncForPC(value.Pointer()).Name())
	}

	if value.Type().NumOut() == 2 && !value.Type().Out(1).Implements(errorInterface) {
		return nil, errors.Errorf("%s second argument of constructor must be error, got %s", runtime.FuncForPC(value.Pointer()).Name(), value.Type().Out(1))
	}

	p := &ConstructorProvider{
		ctor: reflect.ValueOf(rawProvider),
	}

	return p, nil
}

// ConstructorProvider todo
type ConstructorProvider struct {
	ctor reflect.Value
}

// Provide todo
func (p *ConstructorProvider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	result := p.ctor.Call(arguments)

	if len(result) == 1 || result[1].IsNil() {
		return result[0], nil
	}

	return result[0], errors.WithStack(result[1].Interface().(error))
}

// ResultType todo
func (p *ConstructorProvider) ResultType() reflect.Type {
	return p.ctor.Type().Out(0)
}

// Arguments todo
func (p *ConstructorProvider) Arguments() (args []graph.Key) {
	pt := p.ctor.Type()

	for i := 0; i < pt.NumIn(); i++ {
		args = append(args, graph.Key{Type: pt.In(i)})
	}

	return args
}

var errorInterface = reflect.TypeOf(new(error)).Elem()
