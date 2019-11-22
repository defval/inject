package di

import (
	"reflect"

	"github.com/defval/inject/v2/di/internal/reflection"
)

// createConstructor
func createConstructor(name string, ctor interface{}) *constructorProvider {
	if ctor == nil {
		panicf("The constructor must be a function like `func(dep1, dep2...) (result, optionalError)`, got `%s`", "nil")
	}

	if !reflection.IsFunc(ctor) {
		panicf("The constructor must be a function like `func(dep1, dep2...) (result, optionalError)`, got `%s`", reflect.ValueOf(ctor).Type())
	}

	fn := reflection.InspectFunction(ctor)

	if fn.NumOut() == 0 {
		panicf("The constructor `%s` has no results", fn.Name)
	}

	if fn.NumOut() > 2 {
		panicf("The constructor `%s` has many results", fn.Name)
	}

	if fn.NumOut() == 2 && !reflection.IsError(fn.Out(1)) {
		panicf("The second result of constructor `%s` must be error, got %s", fn.Name, fn.Out(1))
	}

	return &constructorProvider{
		name: name,
		ctor: fn,
	}
}

// constructorProvider
type constructorProvider struct {
	name string
	ctor *reflection.Func
}

// identity returns constructor result type identity.
func (c constructorProvider) identity() identity {
	return identity{
		name: c.name,
		typ:  c.ctor.Out(0),
	}
}

// parameters
func (c constructorProvider) parameters() parameterList {
	var parameters parameterList

	for i := 0; i < c.ctor.NumIn(); i++ {
		p := parameter{
			identity: identity{
				typ: c.ctor.In(i),
			},
			optional: false,
		}

		parameters = append(parameters, p)
	}

	return parameters
}

// Provide
func (c constructorProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	out := c.ctor.Call(parameters)

	if len(out) == 1 || out[1].IsNil() {
		return out[0], nil
	}

	return out[0], out[1].Interface().(error)
}
