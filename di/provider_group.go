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
		params: parameterList{},
	}
}

// interfaceGroup
type interfaceGroup struct {
	result identity
	params parameterList
}

// Add
func (i *interfaceGroup) Add(identity identity) {
	i.params = append(i.params, parameter{
		identity: identity,
		optional: false,
	})
}

// identity
func (i interfaceGroup) identity() identity {
	return i.result
}

// parameters
func (i interfaceGroup) parameters() parameterList {
	return i.params
}

// Provide
func (i interfaceGroup) provide(parameters ...reflect.Value) (reflect.Value, error) {
	group := reflect.New(i.result.typ).Elem()
	return reflect.Append(group, parameters...), nil
}
