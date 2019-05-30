package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

// New creates new container with provided options.
func New(options ...Option) (_ *Container, err error) {
	var c = &Container{
		storage: &storage{
			keys:        []key{},
			definitions: map[key]*definition{},
			ifaces:      map[reflect.Type][]*definition{},
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

// Extract populates given target pointer with type instance provided in container.
//
//   var server *http.Server
//   if err = container.Extract(&server); err != nil {
//     // extract failed
//   }
//
//   server.ListenAndServe()
//
// If a target type does not exist in a container or instance type building failed, Extract() returns an error.
// Use ExtractOption for modifying the behavior of this function.
func (c *Container) Extract(target interface{}, options ...ExtractOption) (err error) {
	targetValue := reflect.ValueOf(target)

	// target value type needs to be a pointer
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return errors.New("extract target must be a not nil pointer")
	}

	targetValue = targetValue.Elem()

	var po = &extractOptions{
		target: targetValue,
	}

	// apply extract options
	for _, opt := range options {
		opt.apply(po)
	}

	// create a key to find in a storage
	k := key{
		typ:  po.target.Type(),
		name: po.name,
	}

	newValue, err := c.storage.Value(k)
	if err != nil {
		return errors.WithStack(err)
	}

	targetValue.Set(newValue)

	return nil
}

func (c *Container) compile() (err error) {
	if err = c.registerProviders(); err != nil {
		return errors.WithStack(err)
	}

	if err = c.applyReplacers(); err != nil {
		return errors.WithStack(err)
	}

	if err = c.storage.Compile(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Container) registerProviders() (err error) {
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
	errIncorrectFunctionProviderSignature = errors.New("constructor must be a function with value and optional error as result")
	errorInterface                        = reflect.TypeOf((*error)(nil)).Elem()
)

type providerOptions struct {
	name                 string
	provider             interface{}
	implements           []interface{}
	injectExportedFields bool
}

type extractOptions struct {
	name   string
	target reflect.Value
}
