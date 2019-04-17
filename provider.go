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

// funcProvider
func newFuncProvider(provider interface{}) (*providerWrapper, error) {
	var ptype = reflect.TypeOf(provider)
	var pvalue = reflect.ValueOf(provider)

	if ptype.NumOut() <= 0 || ptype.NumOut() > 2 {
		return nil, errors.WithStack(ErrIncorrectProviderType)
	}

	if ptype.NumOut() == 2 && ptype.Out(1).Implements(errorInterface) == false {
		return nil, errors.WithStack(ErrIncorrectProviderType)
	}

	var resultType = pvalue.Type().Out(0) // todo

	var args []key
	for i := 0; i < ptype.NumIn(); i++ {
		args = append(args, key{typ: ptype.In(i)})
	}

	return &providerWrapper{
		providerType:  providerTypeFunc,
		args:          args,
		providerValue: pvalue,
		resultType:    resultType,
	}, nil
}

// structProvider
func newStructProvider(provider interface{}) (*providerWrapper, error) {
	var ptype = reflect.TypeOf(provider)
	var pvalue = reflect.ValueOf(provider)

	var args []key
	for i := 0; i < ptype.Elem().NumField(); i++ {
		var field = ptype.Elem().Field(i)

		name, exists := field.Tag.Lookup("inject")

		if !exists {
			continue
		}

		args = append(args, key{typ: field.Type, name: name})
	}

	return &providerWrapper{
		providerType:  providerTypeStruct,
		args:          args,
		providerValue: pvalue,
		resultType:    ptype,
	}, nil
}

// providerWrapper
type providerWrapper struct {
	providerType  providerType
	args          []key
	providerValue reflect.Value
	resultType    reflect.Type
}
