package inject

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

const (
	visitMarkUnmarked = iota
	visitMarkTemporary
	visitMarkPermanent
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

// Value creates value of key type
func (k key) Value() reflect.Value {
	return reflect.New(k.typ).Elem()
}

// IsGroup checks that key may be a group
func (k key) IsGroup() bool {
	return k.typ.Kind() == reflect.Slice && k.typ.Elem().Kind() == reflect.Interface
}

// createDefinition.
func createDefinition(po *providerOptions) (def *definition, err error) {
	wrapper, err := wrapProvider(po)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var implements []reflect.Type
	for _, iface := range po.implements {
		ifaceType := reflect.TypeOf(iface)

		if ifaceType.Kind() != reflect.Ptr || ifaceType.Elem().Kind() != reflect.Interface {
			return nil, errors.Errorf("argument for As() must be pointer to interface type, got %s", ifaceType)
		}

		ifaceTypeElem := ifaceType.Elem()

		if !wrapper.rtype().Implements(ifaceTypeElem) {
			return nil, errors.Errorf("%s not implement %s interface", wrapper.rtype(), ifaceTypeElem)
		}

		implements = append(implements, ifaceTypeElem)
	}

	return &definition{
		Key: key{
			typ:  wrapper.rtype(),
			name: po.name,
		},
		Provider:   wrapper,
		Implements: implements,
	}, nil
}

// definition.
type definition struct {
	Key        key
	Provider   providerWrapper
	Implements []reflect.Type
	In         []key
	Out        []key

	instance reflect.Value
	visited  int
}

// value.
func (d *definition) Create(args []reflect.Value) (instance reflect.Value, err error) {
	if d.instance.IsValid() {
		return d.instance, nil
	}

	instance, err = d.Provider.create(args)
	if err != nil {
		return instance, errors.WithStack(err)
	}

	if instance.Kind() == reflect.Ptr && instance.IsNil() {
		return instance, errors.New("nil provided")
	}

	d.instance = instance

	return d.instance, nil
}
