package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

// createCombinedProvider
func createCombinedProvider(value reflect.Value) (wrapper providerWrapper, err error) {
	objectProvider, _ := createObjectProvider(value, true)

	method := value.MethodByName("Provide")
	constructorProvider, err := createConstructorProvider(method)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &combinedProvider{objectProvider, constructorProvider}, nil
}

type combinedProvider struct {
	*objectProvider
	*constructorProvider
}

func (w *combinedProvider) build(arguments []reflect.Value) (value reflect.Value, err error) {
	_, _ = w.objectProvider.build(arguments)
	return w.constructorProvider.build([]reflect.Value{})
}

func (w *combinedProvider) args() []key {
	return w.objectProvider.args()
}

func (w *combinedProvider) rtype() reflect.Type {
	return w.constructorProvider.rtype()
}
