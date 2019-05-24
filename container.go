package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

// New creates new container with provided options.
// Fore more information about container options see `Option` type.
func New(options ...Option) (_ *Container, err error) {
	var c = &Container{
		storage: &storage{
			keys:        make([]key, 0, 8),
			definitions: make(map[key]*definition, 8),
			ifaces:      make(map[reflect.Type][]*definition, 8),
		},
	}

	// apply options.
	for _, opt := range options {
		opt.apply(c)
	}

	if err = c.compile(); err != nil {
		return nil, errors.Wrapf(err, "could not compile container")
	}

	return c, nil
}

// Container is a dependency injection container.
type Container struct {
	providers []*providerOptions
	replacers []*providerOptions
	storage   *storage
}

// Populate populates given target pointer with type instance provided in container.
func (c *Container) Populate(target interface{}, options ...PopulateOption) (err error) {
	rv := reflect.ValueOf(target)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("populate target must be a not nil pointer")
	}

	rv = rv.Elem()

	var po = &populateOptions{
		target: rv,
	}

	for _, opt := range options {
		opt.apply(po)
	}

	k := key{
		typ:  po.target.Type(),
		name: po.name,
	}

	newValue, err := c.storage.Value(k)
	if err != nil {
		return errors.WithStack(err)
	}

	rv.Set(newValue)

	return nil
}

// compile.
func (c *Container) compile() (err error) {
	if err = c.registerProviders(); err != nil {
		return errors.WithStack(err)
	}

	if err = c.applyReplacers(); err != nil {
		return errors.WithStack(err)
	}

	return c.storage.Compile()
}

func (c *Container) registerProviders() (err error) {
	// register providers
	for _, po := range c.providers {
		if po.provider == nil {
			return errors.New("could not provide nil")
		}

		var def *definition
		if def, err = createDefinition(po); err != nil {
			return errors.Wrapf(err, "provide failed")
		}

		if err = c.storage.Add(def); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *Container) applyReplacers() (err error) {
	for _, po := range c.replacers {
		if po.provider == nil {
			return errors.New("replace provider could not be nil")
		}

		var def *definition
		if def, err = createDefinition(po); err != nil {
			return errors.Wrapf(err, "provide failed")
		}

		if err = c.storage.Replace(def); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

var (
	// errIncorrectFunctionProviderSignature.
	errIncorrectFunctionProviderSignature = errors.New("constructor must be a function with value and optional error as result")
	// errorInterface type for error interface implementation checking
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

// providerOptions.
type providerOptions struct {
	provider             interface{}
	name                 string
	implements           []interface{}
	injectExportedFields bool
}

// populateOptions
type populateOptions struct {
	target reflect.Value
	name   string
}
