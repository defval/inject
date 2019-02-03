package injector

import "reflect"

// ID
func interfaceID(i interface{}) reflect.Type {
	return reflect.TypeOf(i).Elem()
}
