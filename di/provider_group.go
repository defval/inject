package di

import (
	"reflect"
)

// createInterfaceGroup creates new group from provided key.
func createInterfaceGroup(key identity) *interfaceGroup {
	return &interfaceGroup{
		result: identity{
			typ: reflect.SliceOf(key.typ),
		},
		keys: []identity{},
	}
}

// interfaceGroup
type interfaceGroup struct {
	result identity
	keys   []identity
}

// Add
func (i *interfaceGroup) Add(key identity) {
	i.keys = append(i.keys, key)
}

// Identity
func (i interfaceGroup) Identity() identity {
	return i.result
}

// Parameters
func (i interfaceGroup) Parameters() parameterList {
	return i.keys
}

// Provide
func (i interfaceGroup) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	group := reflect.New(i.result.typ).Elem()
	return reflect.Append(group, parameters...), nil
}
