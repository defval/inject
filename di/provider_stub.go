package di

import (
	"fmt"
	"reflect"
)

// providerStub
type providerStub struct {
	msg string
	res key
}

// newProviderStub
func newProviderStub(k key, msg string) *providerStub {
	return &providerStub{res: k, msg: msg}
}

func (m *providerStub) Key() key {
	return m.res
}

func (m *providerStub) ParameterList() parameterList {
	return parameterList{}
}

func (m *providerStub) Provide(_ ...reflect.Value) (reflect.Value, error) {
	return reflect.Value{}, fmt.Errorf(m.msg)
}
