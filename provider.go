package inject

import (
	"reflect"
)

// providerWrapper
type providerWrapper interface {
	build(arguments []reflect.Value) (_ reflect.Value, err error)
	args() []key
	rtype() reflect.Type
}

// createProvider creates provider
func createProvider(po *providerOptions) (wrapper providerWrapper, err error) {
	value := reflect.ValueOf(po.provider)

	if value.Kind() == reflect.Func {
		return createConstructorProvider(value)
	}

	if isStructPtr(value) || isStruct(value) {
		return createObjectProvider(value, po.includeExported)
	}

	return &defaultProviderWrapper{
		value: value,
	}, nil
}

// defaultProviderWrapper
type defaultProviderWrapper struct {
	value reflect.Value
}

func (w *defaultProviderWrapper) build(arguments []reflect.Value) (_ reflect.Value, err error) {
	return w.value, nil
}

func (w *defaultProviderWrapper) args() []key {
	return nil
}

func (w *defaultProviderWrapper) rtype() reflect.Type {
	return w.value.Type()
}
