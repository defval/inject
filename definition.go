package inject

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// key.
type key struct {
	// type of provided value
	typ reflect.Type

	// optional name
	name string
}

// String.
func (k key) String() string {
	return fmt.Sprintf("%s", k.typ) // todo: add name
}

// createDefinition.
func createDefinition(po *providerOptions) (def *definition, err error) {
	wrapper, err := wrapProvider(po)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var implements []key
	for _, iface := range po.implements {
		ifaceType := reflect.TypeOf(iface)

		if ifaceType.Kind() != reflect.Ptr || ifaceType.Elem().Kind() != reflect.Interface {
			return nil, errors.Errorf("argument for As() must be pointer to interface type, got %s", ifaceType)
		}

		ifaceTypeElem := ifaceType.Elem()

		if !wrapper.rtype().Implements(ifaceTypeElem) {
			return nil, errors.Errorf("%s not implement %s interface", wrapper.rtype(), ifaceTypeElem)
		}

		implements = append(implements, key{typ: ifaceTypeElem, name: po.name})
	}

	return &definition{
		key: key{
			typ:  wrapper.rtype(),
			name: po.name,
		},
		provider:   wrapper,
		implements: implements,
	}, nil
}

// definition.
type definition struct {
	key        key
	provider   providerWrapper
	implements []key

	in  []*definition
	out []*definition

	instance reflect.Value
	visited  int
}

// String.
func (d *definition) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s", d.key))

	if len(d.implements) > 0 {
		builder.WriteString(" as ")
		for i, k := range d.implements {
			builder.WriteString(fmt.Sprintf("%s", k))

			if i != len(d.implements)-1 {
				builder.WriteString(", ")
			}
		}
	}

	return builder.String()
}

// load.
func (d *definition) load() (instance reflect.Value, err error) {
	if d.instance.IsValid() {
		return d.instance, nil
	}

	var arguments []reflect.Value
	for _, arg := range d.in {
		instance, err := arg.load()

		if err != nil {
			return reflect.Value{}, errors.Wrapf(err, "%s", arg)
		}
		arguments = append(arguments, instance)
	}

	if instance, err = d.provider.create(arguments); err != nil {
		return reflect.Value{}, errors.WithStack(err)
	}

	if instance.Kind() == reflect.Ptr && instance.IsNil() {
		return reflect.Value{}, errors.New("nil provided")
	}

	d.instance = instance

	return instance, nil
}

// visit.
func (d *definition) visit() (err error) {
	if d.visited == visitMarkPermanent {
		return
	}

	if d.visited == visitMarkTemporary {
		return fmt.Errorf("%s", d.key)
	}

	d.visited = visitMarkTemporary

	for _, out := range d.out {
		if err = out.visit(); err != nil {
			return errors.Wrapf(err, "%s", d.key)
		}
	}

	d.visited = visitMarkPermanent

	return nil
}
