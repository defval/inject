package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

// wrapProvider
func wrapProvider(po *providerOptions) (wrapper providerWrapper, err error) {
	pv := reflect.ValueOf(po.provider)
	pt := pv.Type()

	switch true {
	case pt.Kind() == reflect.Func:
		if err = checkFunctionProvider(pt); err != nil {
			return nil, errors.WithStack(err)
		}

		return &funcProviderWrapper{
			value: pv,
		}, nil
	case pt.Kind() == reflect.Ptr && pt.Elem().Kind() == reflect.Struct:
		return &structPointerProviderWrapper{
			value:                pv,
			injectExportedFields: po.injectExportedFields,
		}, nil
	default:
		return &defaultProviderWrapper{
			value: pv,
		}, nil
	}
}

// check function provider
func checkFunctionProvider(pt reflect.Type) (err error) {
	// check function result types
	if pt.NumOut() <= 0 || pt.NumOut() > 2 {
		return errors.WithStack(errIncorrectFunctionProviderSignature)
	}

	if pt.NumOut() == 2 && !pt.Out(1).Implements(errorInterface) {
		return errors.WithStack(errIncorrectFunctionProviderSignature)
	}

	return nil
}

// providerWrapper encapsulates creating an instance
type providerWrapper interface {
	create(arguments []reflect.Value) (_ reflect.Value, err error)
	args() []key
	rtype() reflect.Type
}

type defaultProviderWrapper struct {
	value reflect.Value
}

func (w *defaultProviderWrapper) create(arguments []reflect.Value) (_ reflect.Value, err error) {
	return w.value, nil
}

func (w *defaultProviderWrapper) args() []key {
	return nil
}

func (w *defaultProviderWrapper) rtype() reflect.Type {
	return w.value.Type()
}

type structPointerProviderWrapper struct {
	value                reflect.Value
	injectExportedFields bool
}

func (w *structPointerProviderWrapper) create(arguments []reflect.Value) (_ reflect.Value, err error) {
	pe := w.value.Elem()

	skip := 0
	for i := 0; i < pe.Type().NumField(); i++ {
		var structField = pe.Type().Field(i)
		var fieldValue = pe.Field(i)

		_, exists := structField.Tag.Lookup("inject")

		if fieldValue.CanSet() && (exists || w.injectExportedFields) {
			fieldValue.Set(arguments[i-skip])
		} else {
			skip++
		}
	}

	return w.value, nil
}

func (w *structPointerProviderWrapper) args() []key {
	pv := w.value

	var args []key
	for i := 0; i < pv.Type().Elem().NumField(); i++ {
		structField := pv.Type().Elem().Field(i)
		fieldValue := pv.Elem().Field(i)

		name, exists := structField.Tag.Lookup("inject")

		if fieldValue.CanSet() && (exists || w.injectExportedFields) {
			args = append(args, key{typ: structField.Type, name: name})
		}
	}

	return args
}

func (w *structPointerProviderWrapper) rtype() reflect.Type {
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

	return result[0], errors.WithStack(result[1].Interface().(error))
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
