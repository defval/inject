package di

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/defval/inject/v2/di/internal/reflection"
)

type ctorType int

const (
	ctorUnknown      ctorType = iota // unknown ctor signature
	ctorSimple                       // (deps) (result)
	ctorError                        // (deps) (result, error)
	ctorCleanup                      // (deps) (result, cleanup)
	ctorCleanupError                 // (deps) (result, cleanup, error)
)

// determineCtorType
func determineCtorType(fn *reflection.Func) ctorType {
	if fn.NumOut() == 1 {
		return ctorSimple
	}

	if fn.NumOut() == 2 {
		if reflection.IsError(fn.Out(1)) {
			return ctorError
		}

		if reflection.IsCleanup(fn.Out(1)) {
			return ctorCleanup
		}
	}

	if fn.NumOut() == 3 && reflection.IsCleanup(fn.Out(1)) && reflection.IsError(fn.Out(2)) {
		return ctorCleanupError
	}

	panic(fmt.Sprintf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", fn.Name))
}

// createConstructor
func createConstructor(name string, ctor interface{}) *constructorProvider {
	if ctor == nil {
		panicf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", "nil")
	}

	if !reflection.IsFunc(ctor) {
		panicf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", reflect.ValueOf(ctor).Type())
	}

	fn := reflection.InspectFunction(ctor)
	ctorType := determineCtorType(fn)

	return &constructorProvider{
		name:     name,
		ctor:     fn,
		ctorType: ctorType,
	}
}

// constructorProvider
type constructorProvider struct {
	name     string
	ctor     *reflection.Func
	ctorType ctorType
	clean    *reflection.Func
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
func (c *constructorProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	out := c.ctor.Call(parameters)

	switch c.ctorType {
	case ctorSimple:
		return out[0], nil
	case ctorError:
		instance := out[0]
		err := out[1]

		if err.IsNil() {
			return instance, nil
		}

		return instance, err.Interface().(error)
	case ctorCleanup:
		c.saveCleanup(out[1])
		return out[0], nil
	case ctorCleanupError:
		instance := out[0]
		cleanup := out[1]
		err := out[2]

		c.saveCleanup(cleanup)

		if err.IsNil() {
			return instance, nil
		}

		return instance, err.Interface().(error)
	}

	return reflect.Value{}, errors.New("you found a bug, please create new issue for " +
		"this: https://github.com/defval/inject/issues/new")
}

func (c *constructorProvider) saveCleanup(value reflect.Value) {
	c.clean = reflection.InspectFunction(value.Interface())
}

func (c *constructorProvider) cleanup() {
	if c.clean != nil && c.clean.IsValid() {
		c.clean.Call([]reflect.Value{})
	}
}
