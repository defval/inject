package di

import (
	"reflect"
	"strings"
)

// createStructProvider
func createStructProvider(name string, provider Provider) *structProvider {
	result := provider.Identity()
	value := reflect.ValueOf(provider)

	return &structProvider{
		id: identity{
			name: name,
			typ:  result.typ,
		},
		provider: reflect.Indirect(value),
		typ:      reflect.Indirect(value).Type(),
		ctor:     createConstructor(name, value.MethodByName("Provide").Interface()),
	}
}

type structProvider struct {
	id       identity
	provider reflect.Value
	typ      reflect.Type
	ctor     *constructorProvider
}

func (s *structProvider) identity() identity {
	return s.id
}

func (s *structProvider) parameters() parameterList {
	var pl parameterList

	for i := 0; i < s.typ.NumField(); i++ {
		name, optional, isDependency := s.inspectField(i)
		if !isDependency {
			continue
		}

		pl = append(pl, parameter{
			identity: identity{
				name: name,
				typ:  s.provider.Field(i).Type(),
			},
			optional: optional,
		})
	}

	return pl
}

func (s *structProvider) inspectField(num int) (name string, optional bool, isDependency bool) {
	value, exists := s.typ.Field(num).Tag.Lookup("di")
	if !exists {
		return "", false, false
	}

	name, optional = s.parseTag(value)

	return name, optional, true
}

func (s *structProvider) parseTag(tag string) (name string, optional bool) {
	options := strings.Split(tag, ",")
	if len(options) == 0 {
		return "", false
	}

	if len(options) == 1 && options[0] == "optional" {
		return "", true
	}

	if len(options) == 1 {
		return options[0], false
	}

	if len(options) == 2 && options[1] == "optional" {
		return options[0], true
	}

	panic("incorrect di tag")
}

func (s *structProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	return s.ctor.provide()
}
