package inject

import (
	"reflect"
)

func isStructPtr(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr && reflect.Indirect(value).Kind() == reflect.Struct
}

func isStruct(value reflect.Value) bool {
	return value.Kind() == reflect.Struct
}

// createObjectProvider creates object provider.
func createObjectProvider(value reflect.Value, includeExported bool) (_ *objectProvider, err error) {
	if isStructPtr(value) {
		return &objectProvider{
			id:              value.Type().String(),
			value:           value,
			isPointer:       true,
			includeExported: includeExported,
		}, nil
	}

	return &objectProvider{
		id:              value.Type().String(),
		value:           reflect.New(value.Type()),
		includeExported: includeExported,
	}, nil
}

type objectProvider struct {
	id              string // for debug
	value           reflect.Value
	isPointer       bool
	includeExported bool
}

func (w *objectProvider) build(arguments []reflect.Value) (_ reflect.Value, err error) {
	value := w.value.Elem()

	skip := 0
	for i := 0; i < value.NumField(); i++ {
		_, _, injectable := w.isFieldInjectable(i)
		if injectable {
			value.Field(i).Set(arguments[i-skip])
		} else {
			skip++
		}
	}

	if w.isPointer {
		return w.value, nil
	}

	return reflect.Indirect(w.value), nil
}

func (w *objectProvider) args() []key {
	value := w.value.Elem()

	var args []key
	for i := 0; i < value.NumField(); i++ {
		typ, name, injectable := w.isFieldInjectable(i)
		if !injectable {
			continue
		}

		args = append(args, key{typ: typ, name: name})
	}

	return args
}

func (w *objectProvider) rtype() reflect.Type {
	if !w.isPointer {
		return reflect.Indirect(w.value).Type()
	}

	return w.value.Type()
}

func (w *objectProvider) isFieldInjectable(fieldNum int) (typ reflect.Type, name string, injectable bool) {
	value := w.value.Elem()

	name, exists := value.Type().Field(fieldNum).Tag.Lookup("inject")
	return value.Type().Field(fieldNum).Type, name, value.Field(fieldNum).CanSet() && (exists || w.includeExported)
}
