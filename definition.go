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

// definition
type definition struct {
	key        key
	provider   *providerWrapper
	implements []key

	in []*definition

	value reflect.Value
}

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
	var values []reflect.Value
	for _, in := range d.in {
		var value reflect.Value
		if value, err = in.instance(); err != nil {
			return value, errors.Wrapf(err, "%s", in.key)
		}
		values = append(values, value)
	}

	var result = d.provider.providerValue.Call(values)
	d.value = result[0]

	if len(result) == 2 {
		if result[1].IsNil() {
			return d.value, nil
		}

		return d.value, errors.WithStack(result[1].Interface().(error))
	}

	return d.value, nil
}
