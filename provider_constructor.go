package inject

import (
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

// createConstructorProvider creates constructor provider.
func createConstructorProvider(value reflect.Value) (_ *constructorProvider, err error) {
	fn := runtime.FuncForPC(value.Pointer())

	if value.Type().NumOut() == 0 {
		return nil, errors.Errorf("%s must have at least one return value", fn.Name())
	}

	if value.Type().NumOut() > 2 {
		return nil, errors.Errorf("%s: constructor may have maximum two return values", fn.Name())
	}

	if value.Type().NumOut() == 2 && !value.Type().Out(1).Implements(errorInterface) {
		return nil, errors.Errorf("%s: second argument of constructor must be error, got %s", fn.Name(), value.Type().Out(1))
	}

	return &constructorProvider{
		cfn: value,
	}, nil
}

type constructorProvider struct {
	cfn reflect.Value // constructor function
}

func (w *constructorProvider) build(arguments []reflect.Value) (_ reflect.Value, err error) {
	result := w.cfn.Call(arguments)

	if len(result) == 1 || result[1].IsNil() {
		return result[0], nil
	}

	return result[0], errors.WithStack(result[1].Interface().(error))
}

func (w *constructorProvider) args() []key {
	pt := w.cfn.Type()

	var args []key
	for i := 0; i < pt.NumIn(); i++ {
		args = append(args, key{typ: pt.In(i)})
	}

	return args
}

func (w *constructorProvider) rtype() reflect.Type {
	return w.cfn.Type().Out(0)
}
