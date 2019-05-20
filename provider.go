package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

// providerWrapper
type providerWrapper interface {
	create(arguments []reflect.Value) (_ reflect.Value, err error)
	args() []key
	rtype() reflect.Type
}

type structProviderWrapper struct {
	value reflect.Value
}

func (w *structProviderWrapper) create(arguments []reflect.Value) (_ reflect.Value, err error) {
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

func (w *structProviderWrapper) args() []key {
	pt := w.value.Type()

	var args []key
	for i := 0; i < pt.Elem().NumField(); i++ {
		var field = pt.Elem().Field(i)

		name, exists := field.Tag.Lookup("inject")

		if !exists {
			continue
		}

		args = append(args, key{typ: field.Type, name: name})
	}

	return args
}

func (w *structProviderWrapper) rtype() reflect.Type {
	return w.value.Type()
}

type funcProviderWrapper struct {
	value     reflect.Value
	result    reflect.Type
	arguments []key
}

func (w *funcProviderWrapper) create(arguments []reflect.Value) (_ reflect.Value, err error) {
	var result = w.value.Call(arguments)

	if len(result) == 1 || result[1].IsNil() {
		return result[0], nil
	}

	if len(result) == 2 {
		return result[0], errors.WithStack(result[1].Interface().(error))
	}

	panic("incorrect constructor function")
}

func (w *funcProviderWrapper) args() []key {
	pt := w.value.Type()

	var args []key
	for i := 0; i < pt.NumIn(); i++ {
		args = append(args, key{typ: pt.In(i)})
	}

	return args
}

func (w *funcProviderWrapper) rtype() reflect.Type {
	return w.value.Type().Out(0)
}

// wrapProvider
func wrapProvider(provider interface{}) (wrapper providerWrapper, err error) {
	pv := reflect.ValueOf(provider)
	pt := pv.Type()

	if pt.Kind() == reflect.Func {
		if err = checkFunctionProvider(pt); err != nil {
			return nil, errors.WithStack(err)
		}

		return &funcProviderWrapper{
			value: pv,
		}, nil
	}

	if pt.Kind() == reflect.Ptr && pt.Elem().Kind() == reflect.Struct {
		return &structProviderWrapper{
			value: pv,
		}, nil
	}

	return nil, errors.WithStack(errIncorrectProviderType)
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
