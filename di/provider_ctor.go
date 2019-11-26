package di

import (
	"reflect"

	"github.com/defval/inject/v2/di/internal/reflection"
)

// createConstructor
func createConstructor(name string, ctor interface{}) *constructorProvider {
	if ctor == nil {
		panicf("The constructor must be a function like `func(dep1, dep2...) (result, cleanup, error)`, got `%s`", "nil")
	}

	if !reflection.IsFunc(ctor) {
		panicf("The constructor must be a function like `func(dep1, dep2...) (result, cleanup, error)`, got `%s`", reflect.ValueOf(ctor).Type())
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

// resultKey returns constructor result type resultKey.
func (c constructorProvider) resultKey() key {
	return key{
		name: c.name,
		typ:  c.ctor.Out(0),
	}
}

// parameters returns constructor parameters
func (c constructorProvider) parameters() parameterList {
	var list parameterList

	for i := 0; i < c.ctor.NumIn(); i++ {
		ptype := c.ctor.In(i)

		p := parameter{
			key:      key{typ: ptype},
			optional: false,
			embed:    isEmbedParameter(ptype),
		}

		list = append(list, p)
	}

	return list
}

// Provide
func (c constructorProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	out := c.ctor.Call(parameters)

	if len(out) == 1 || out[1].IsNil() {
		return out[0], nil
	}

	return out[0], out[1].Interface().(error)
}
