package combined

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/provider"
	"github.com/defval/inject/internal/provider/ctor"
	"github.com/defval/inject/internal/provider/object"
)

// IsValid
func IsValid(rawProvider interface{}) (err error) {
	if err = object.IsValid(rawProvider); err != nil {
		return errors.WithStack(err)
	}

	value := reflect.ValueOf(rawProvider)

	if !strings.HasSuffix(value.Type().String(), "Provider") {
		return errors.Errorf("combined provider must have Provider suffix")
	}

	ctorMethod := value.MethodByName("Provide")

	if !ctorMethod.IsValid() {
		return errors.Errorf("combined provider must have Provide() method")
	}

	if err = ctor.IsValid(ctorMethod.Interface()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// New
func New(rawProvider interface{}, options ...object.Option) (*Provider, error) {
	if err := IsValid(rawProvider); err != nil {
		return nil, errors.WithStack(err)
	}

	value := reflect.ValueOf(rawProvider)
	ctorMethod := value.MethodByName("Provide")

	structPtrProvider, err := object.New(rawProvider, options...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ctorProvider, err := ctor.New(ctorMethod.Interface())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Provider{
		object: structPtrProvider,
		ctor:   ctorProvider,
	}, nil
}

// Provider
type Provider struct {
	object *object.Provider
	ctor   *ctor.Provider
}

func (p *Provider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	_, _ = p.object.Provide(arguments)
	return p.ctor.Provide([]reflect.Value{})
}

func (p *Provider) ResultType() reflect.Type {
	return p.ctor.ResultType()
}

func (p *Provider) Arguments() (args []provider.Key) {
	return p.object.Arguments()
}
