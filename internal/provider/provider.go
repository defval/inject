package provider

import (
	"fmt"
	"reflect"
	"strings"
)

// Type is provider type
type Type int

const (
	Constructor Type = iota + 1
	Combined
	Object
	Direct
)

// Key
type Key struct {
	Type reflect.Type
	Name string
}

// deprecated
func (k Key) IsGroup() bool {
	return k.Type.Kind() == reflect.Slice && k.Type.Elem().Kind() == reflect.Interface
}

// Value creates new value of key type.
func (k Key) Value() reflect.Value {
	return reflect.New(k.Type).Elem()
}

func (k Key) String() string {
	return fmt.Sprintf("%s", k.Type) // todo: add name
}

// Provider
type Provider interface {
	fmt.Stringer
	Provide(arguments []reflect.Value) (reflect.Value, error)
	ResultType() reflect.Type
	Arguments() (args []Key)
}

// DetectType
func DetectType(rawProvider interface{}) Type {
	value := reflect.ValueOf(rawProvider)

	if value.Kind() == reflect.Func {
		return Constructor
	}
	_, provideMethodExists := value.Type().MethodByName("Provide")
	if strings.HasSuffix(value.Type().String(), "Provider") && provideMethodExists {
		return Combined
	}

	isStruct := value.Kind() == reflect.Struct
	isStructPtr := value.Kind() == reflect.Ptr && reflect.Indirect(value).Kind() == reflect.Struct

	if isStruct || isStructPtr {
		return Object
	}

	return Direct
}
