package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

type providerType int

const (
	providerTypeFunc providerType = iota
	providerTypeStruct
)

// providerWrapper
type providerWrapper struct {
	wrapperType providerType
	value       reflect.Value
	arguments   []key
	result      reflect.Type
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

	return nil, errors.WithStack(ErrIncorrectProviderType)
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
		wrapperType: providerTypeFunc,
		arguments:   args,
		value:       pv,
		result:      result,
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
		wrapperType: providerTypeStruct,
		value:       pv,
		arguments:   args,
		result:      pt,
	}, nil
}

// check function provider
func checkFunctionProvider(pt reflect.Type) (err error) {
	// check function result types
	if pt.NumOut() <= 0 || pt.NumOut() > 2 {
		return errors.WithStack(ErrIncorrectProviderType)
	}

	if pt.NumOut() == 2 && !pt.Out(1).Implements(errorInterface) {
		return errors.WithStack(ErrIncorrectProviderType)
	}

	return nil
}
