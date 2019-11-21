package reflection

import "reflect"

var errorInterface = reflect.TypeOf(new(error)).Elem()

// IsError
func IsError(typ reflect.Type) bool {
	return typ.Implements(errorInterface)
}

// IsPtr
func IsPtr(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr
}
