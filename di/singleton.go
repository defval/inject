package di

import (
	"reflect"
)

// asSingleton creates a singleton wrapper.
func asSingleton(provider provider) *singletonWrapper {
	return &singletonWrapper{provider: provider}
}

// singletonWrapper is a embedParamProvider wrapper. Stores provided value for prevent reinitialization.
type singletonWrapper struct {
	provider               // source provider
	value    reflect.Value // value cache
}

// Provide
func (s *singletonWrapper) provide(parameters ...reflect.Value) (reflect.Value, error) {
	if s.value.IsValid() {
		return s.value, nil
	}

	value, err := s.provider.provide(parameters...)
	if err != nil {
		return reflect.Value{}, err
	}

	s.value = value

	return value, nil
}
