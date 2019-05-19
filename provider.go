package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

// providerType.
type providerType int

const (
	providerTypeFunc providerType = iota
	providerTypeStruct
)

// providerWrapper
type providerWrapper struct {
	typ       providerType
	value     reflect.Value
	result    reflect.Type
	arguments []key
}

// instance
func (w *providerWrapper) instance(arguments []reflect.Value) (_ reflect.Value, err error) {
	switch w.typ {
	case providerTypeFunc:
		var result = w.value.Call(arguments)

		if len(result) == 1 || result[1].IsNil() {
			return result[0], nil
		}

		if len(result) == 2 {
			return result[0], errors.WithStack(result[1].Interface().(error))
		}

		panic("incorrect constructor function")
	case providerTypeStruct:
		pe := w.value.Elem()

		skip := 0
		for i := 0; i < pe.Type().NumField(); i++ {
			var fieldType = pe.Type().Field(i)
			var fieldValue = pe.Field(i)

			_, exists := fieldType.Tag.Lookup("inject")

			if !exists {
				skip++
				continue
			}

			fieldValue.Set(arguments[i-skip])
		}

		if w.value.Kind() == reflect.Ptr && w.value.IsNil() {
			return w.value, errors.Errorf("nil provided")
		}

		return w.value, nil
	}

	panic("unknown provider type")

}

// wrapProvider
func wrapProvider(provider interface{}) (wrapper *providerWrapper, err error) {
	var pt = reflect.TypeOf(provider)

	if pt.Kind() == reflect.Func {
		return wrapFunction(provider)
	}

	if pt.Kind() == reflect.Ptr && pt.Elem().Kind() == reflect.Struct {
		return wrapStruct(provider)
	}

	return nil, errors.WithStack(errIncorrectProviderType)
}

// wrapFunction
func wrapFunction(provider interface{}) (_ *providerWrapper, err error) {
	// provider value
	var pv = reflect.ValueOf(provider)

	// provider type
	var pt = pv.Type()

	if err = checkFunctionProvider(pt); err != nil {
		return nil, err
	}

	var result = pv.Type().Out(0)

	var args []key
	for i := 0; i < pt.NumIn(); i++ {
		args = append(args, key{typ: pt.In(i)})
	}

	return &providerWrapper{
		typ:       providerTypeFunc,
		arguments: args,
		value:     pv,
		result:    result,
	}, nil
}

// structProvider
func wrapStruct(provider interface{}) (*providerWrapper, error) {
	// provider value
	var pv = reflect.ValueOf(provider)

	// provider type
	var pt = pv.Type()

	var args []key
	for i := 0; i < pt.Elem().NumField(); i++ {
		var field = pt.Elem().Field(i)

		name, exists := field.Tag.Lookup("inject")

		if !exists {
			continue
		}

		args = append(args, key{typ: field.Type, name: name})
	}

	return &providerWrapper{
		typ:       providerTypeStruct,
		value:     pv,
		arguments: args,
		result:    pt,
	}, nil
}

// check function provider
func checkFunctionProvider(pt reflect.Type) (err error) {
	// check function result types
	if pt.NumOut() <= 0 || pt.NumOut() > 2 {
		return errors.WithStack(errIncorrectProviderType)
	}

	if pt.NumOut() == 2 && !pt.Out(1).Implements(errorInterface) {
		return errors.WithStack(errIncorrectProviderType)
	}

	return nil
}
