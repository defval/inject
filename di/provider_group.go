package di

import (
	"reflect"
)

// newGroupProvider creates new group from provided resultKey.
func newGroupProvider(k key) *interfaceGroup {
	ifaceKey := key{
		res: reflect.SliceOf(k.res),
		typ: ptGroup,
	}

	return &interfaceGroup{
		result: ifaceKey,
		pl:     parameterList{},
	}
}

// interfaceGroup
type interfaceGroup struct {
	result key
	pl     parameterList
}

// Add
func (i *interfaceGroup) Add(k key) {
	i.pl = append(i.pl, parameter{
		name:     k.name,
		res:      k.res,
		optional: false,
		embed:    false,
	})
}

// resultKey
func (i interfaceGroup) Key() key {
	return i.result
}

// parameters
func (i interfaceGroup) ParameterList() parameterList {
	return i.pl
}

// Provide
func (i interfaceGroup) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	group := reflect.New(i.result.res).Elem()
	return reflect.Append(group, parameters...), nil
}
