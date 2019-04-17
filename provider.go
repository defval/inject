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

	return &providerWrapper{
		providerType:  providerTypeFunc,
		providerValue: pvalue,
		resultType:    resultType,
	}, nil
}

// structProvider
func newStructProvider(provider interface{}) (*providerWrapper, error) {
	var ptype = reflect.TypeOf(provider)
	var pvalue = reflect.ValueOf(provider)

	// todo: add validation

	return &providerWrapper{
		providerType:  providerTypeStruct,
		providerValue: pvalue,
		resultType:    ptype,
	}, nil
}

// providerWrapper
type providerWrapper struct {
	providerType  providerType
	providerValue reflect.Value
	resultType    reflect.Type
}

// args
func (w *providerWrapper) args() (args []key) {
	switch w.providerType {
	case providerTypeFunc:
		for i := 0; i < w.providerValue.Type().NumIn(); i++ {
			args = append(args, key{typ: w.providerValue.Type().In(i)})
		}
	case providerTypeStruct:
		for i := 0; i < w.resultType.Elem().NumField(); i++ {
			var field = w.resultType.Elem().Field(i)

			name, exists := field.Tag.Lookup("inject")

			if !exists {
				continue
			}

			args = append(args, key{typ: field.Type, name: name})
		}
	}

	return args
}
