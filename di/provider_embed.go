package di

import (
	"reflect"
	"strings"
)

// createStructProvider
func createEmbedProvider(p parameter) *embedParamProvider {
	result := p.resultKey()

	var embedType reflect.Type
	if result.typ.Kind() == reflect.Ptr {
		embedType = result.typ.Elem()
	} else {
		embedType = result.typ
	}

	return &embedParamProvider{
		key:        result,
		embedType:  embedType,
		embedValue: reflect.New(embedType).Elem(),
	}
}

type embedParamProvider struct {
	key        key
	embedType  reflect.Type
	embedValue reflect.Value
}

func (s *embedParamProvider) resultKey() key {
	return s.key
}

func (s *embedParamProvider) parameters() parameterList {
	var pl parameterList

	for i := 0; i < s.embedType.NumField(); i++ {
		name, optional, isDependency := s.inspectField(i)
		if !isDependency {
			continue
		}

		// parameter field
		pField := s.embedType.Field(i)

		pl = append(pl, parameter{
			key: key{
				name: name,
				typ:  pField.Type,
			},
			optional: optional,
			embed:    isEmbedParameter(pField.Type),
		})
	}

	return pl
}

func (s *embedParamProvider) inspectField(num int) (name string, optional bool, isDependency bool) {
	tag, tagExists := s.embedType.Field(num).Tag.Lookup("di")
	canSet := s.embedValue.Field(num).CanSet()
	if !tagExists || !canSet {
		return "", false, false
	}

	name, optional = s.parseTag(tag)

	return name, optional, true
}

func (s *embedParamProvider) parseTag(tag string) (name string, optional bool) {
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

func (s *embedParamProvider) provide(parameters ...reflect.Value) (reflect.Value, error) {
	for i, offset := 0, 0; i < s.embedType.NumField(); i++ {
		_, _, isDependency := s.inspectField(i)
		if !isDependency {
			offset++
			continue
		}

		s.embedValue.Field(i).Set(parameters[i-offset])
	}

	return s.embedValue, nil
}
