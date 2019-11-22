package di

import (
	"reflect"
)

// asSingleton creates a singleton wrapper.
func asSingleton(provider dependencyProvider) *singletonWrapper {
	return &singletonWrapper{dependencyProvider: provider}
}

// singletonWrapper is a structProvider wrapper. Stores provided value for prevent reinitialization.
type singletonWrapper struct {
	dependencyProvider               // source structProvider
	value              reflect.Value // value cache
}

// Provide
func (s *singletonWrapper) provide(parameters ...reflect.Value) (reflect.Value, error) {
	if s.value.IsValid() {
		return s.value, nil
	}

	value, err := s.dependencyProvider.provide(parameters...)
	if err != nil {
		return reflect.Value{}, err
	}

	s.value = value

	return value, nil
}
