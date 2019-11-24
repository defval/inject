package di

import (
	"reflect"
)

// createInterfaceGroup creates new group from provided resultKey.
func createInterfaceGroup(k key) *interfaceGroup {
	ifaceKey := key{
		typ: reflect.SliceOf(k.typ),
	}

	return &interfaceGroup{
		result: ifaceKey,
		pl: providerParameterList{
			providerKey: ifaceKey,
		},
	}
}

// interfaceGroup
type interfaceGroup struct {
	result key
	pl     providerParameterList
}

// Add
func (i *interfaceGroup) Add(k key) {
	i.pl.add(parameterRequired{k})
}

// resultKey
func (i interfaceGroup) resultKey() key {
	return i.result
}

// parameters
func (i interfaceGroup) parameters() providerParameterList {
	return i.pl
}

// Provide
func (i interfaceGroup) provide(parameters ...reflect.Value) (reflect.Value, error) {
	group := reflect.New(i.result.typ).Elem()
	return reflect.Append(group, parameters...), nil
}
