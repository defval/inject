package di

import (
	"reflect"
)

// createInterfaceGroup creates new group from provided key.
func createInterfaceGroup(key providerKey) *interfaceGroup {
	return &interfaceGroup{
		result: providerKey{
			Type: reflect.SliceOf(key.Type),
		},
		keys: []providerKey{},
	}
}

// interfaceGroup
type interfaceGroup struct {
	result providerKey
	keys   []providerKey
}

// Add
func (i *interfaceGroup) Add(key providerKey) {
	i.keys = append(i.keys, key)
}

// Result
func (i interfaceGroup) Result() providerKey {
	return i.result
}

// Parameters
func (i interfaceGroup) Parameters() parameterList {
	return i.keys
}

// Provide
func (i interfaceGroup) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	group := reflect.New(i.result.Type).Elem()
	return reflect.Append(group, parameters...), nil
}
