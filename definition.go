package inject

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type key struct {
	// type of provided value
	typ reflect.Type

	// optional name
	name string
}

// String
func (k key) String() string {
	return fmt.Sprintf("%s", k.typ)
}

// createDefinition
func createDefinition(po *providerOptions) (def *definition, err error) {
	var wrapper *providerWrapper
	if wrapper, err = wrapProvider(po.provider); err != nil {
		return nil, errors.WithStack(err)
	}

	var implements []key
	for _, iface := range po.implements {
		ifaceType := reflect.TypeOf(iface)

		if ifaceType.Kind() != reflect.Ptr || ifaceType.Elem().Kind() != reflect.Interface {
			return nil, errors.Errorf("argument for As() must be pointer to interface type, got %s", ifaceType)
		}

		ifaceTypeElem := ifaceType.Elem()

		if !wrapper.result.Implements(ifaceTypeElem) {
			return nil, errors.Errorf("%s not implement %s interface", wrapper.result, ifaceTypeElem)
		}

		implements = append(implements, key{typ: ifaceTypeElem})
	}

	return &definition{
		key: key{
			typ:  wrapper.result,
			name: po.name,
		},
		provider:   wrapper,
		implements: implements,
	}, nil
}

// definition
type definition struct {
	key        key
	provider   *providerWrapper
	implements []key

	in  []*definition
	out []*definition

	value   reflect.Value
	visited int
}

// String
func (d *definition) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s", d.key))

	if len(d.implements) > 0 {
		builder.WriteString(" as ")
		for i, key := range d.implements {
			builder.WriteString(fmt.Sprintf("%s", key.typ))

			if i != len(d.implements)-1 {
				builder.WriteString(", ")
			}
		}
	}

	if d.key.name != "" {
		builder.WriteString(fmt.Sprintf(" with name `%s`", d.key.name))
	}

	return builder.String()
}

func (d *definition) instance() (_ reflect.Value, err error) {
	if d.value.IsValid() {
		return d.value, nil
	}

	var values []reflect.Value
	for _, in := range d.in {
		var value reflect.Value
		if value, err = in.instance(); err != nil {
			return value, errors.Wrapf(err, "%s", in.key)
		}
		values = append(values, value)
	}

	switch d.provider.wrapperType {
	case providerTypeFunc:
		var result = d.provider.value.Call(values)

		d.value = result[0]

		if len(result) == 2 {
			if result[1].IsNil() {
				return d.value, nil
			}

			return d.value, errors.WithStack(result[1].Interface().(error))
		}
	case providerTypeStruct:
		var skip = 0
		for i := 0; i < d.provider.value.Elem().Type().NumField(); i++ {
			var fieldType = d.provider.value.Elem().Type().Field(i)
			var fieldValue = d.provider.value.Elem().Field(i)

			_, exists := fieldType.Tag.Lookup("inject")

			if !exists {
				skip++
				continue
			}

			fieldValue.Set(values[i-skip])
		}

		d.value = d.provider.value
	}

	if d.value.Kind() == reflect.Ptr && d.value.IsNil() {
		return d.value, errors.Errorf("nil provided")
	}

	return d.value, nil
}

func (d *definition) visit() (err error) {
	if d.visited == visitMarkPermanent {
		return
	}

	if d.visited == visitMarkTemporary {
		return fmt.Errorf("%s", d.provider.result)
	}

	d.visited = visitMarkTemporary

	for _, out := range d.out {
		if err = out.visit(); err != nil {
			return errors.Wrapf(err, "%s", d.provider.result)
		}
	}

	d.visited = visitMarkPermanent

	return nil
}
