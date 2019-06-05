package inject

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/provider"
)

const (
	visitMarkUnmarked = iota
	visitMarkTemporary
	visitMarkPermanent
)

func createDefinition(po *providerOptions) (def *definition, err error) {
	wrapper, err := createProvider(po)

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

		if !wrapper.ResultType().Implements(ifaceTypeElem) {
			return nil, errors.Errorf("%s not implement %s interface", wrapper.ResultType(), ifaceTypeElem)
		}

		implements = append(implements, ifaceTypeElem)
	}

	return &definition{
		Key: provider.Key{
			Type: wrapper.ResultType(),
			Name: po.name,
		},
		Provider:   wrapper,
		Implements: implements,
	}, nil
}

type definition struct {
	Key        provider.Key
	Provider   provider.Provider
	Implements []reflect.Type
	In         []provider.Key
	Out        []provider.Key

	instance reflect.Value
	visited  int
}

func (d *definition) Create(args []reflect.Value) (instance reflect.Value, err error) {
	if d.instance.IsValid() {
		return d.instance, nil
	}

	instance, err = d.Provider.Provide(args)
	if err != nil {
		return instance, errors.WithStack(err)
	}

	if instance.Kind() == reflect.Ptr && instance.IsNil() {
		return instance, errors.New("nil provided")
	}

	d.instance = instance

	return d.instance, nil
}
