package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

// New creates new container with provided options.
// Fore more information about container options see `Option` type.
func New(options ...Option) (_ *Container, err error) {
	var c = &Container{
		storage: make(map[string]*storage),
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
	storage   map[string]*storage
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

	newValue, err := c.getStorage(po.namespace).Value(k)
	if err != nil {
		return errors.WithStack(err)
	}

	rv.Set(newValue)

	return nil
}

func (c *Container) getStorage(namespace string) *storage {
	if _, ok := c.storage[namespace]; ok {
		return c.storage[namespace]
	}

	c.storage[namespace] = &storage{
		keys:        []key{},
		definitions: map[key]*definition{},
		ifaces:      map[reflect.Type][]*definition{},
	}

	return c.storage[namespace]
}

// compile.
func (c *Container) compile() (err error) {
	if err = c.registerProviders(); err != nil {
		return errors.WithStack(err)
	}

	if err = c.applyReplacers(); err != nil {
		return errors.WithStack(err)
	}

	for _, storage := range c.storage {
		if err = storage.Compile(); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
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

		if err = c.getStorage(po.namespace).Add(def); err != nil {
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

		if err = c.getStorage(po.namespace).Replace(def); err != nil {
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
	name                 string
	provider             interface{}
	namespace            string
	implements           []interface{}
	injectExportedFields bool
}

// populateOptions
type populateOptions struct {
	name      string
	namespace string
	target    reflect.Value
}
