package inject

import (
	"io"
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/graph"
	"github.com/defval/inject/internal/graph/dot"
)

// New creates a new container with provided options.
func New(options ...Option) (_ *Container, err error) {
	var c = &Container{
		storage: graph.NewStorage(),
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
	storage   *graph.Storage
}

// WriteTo writes container entities as a graphviz dot nodes to writer, like
// https://raw.githubusercontent.com/defval/inject/master/graph.png.
func (c *Container) WriteTo(w io.Writer) (n int64, err error) {
	dot.NewGraphFromStorage(c.storage).Write(w)
	return n, err
}

// Extract populates given target pointer with type instance provided in the container.
//
//   var server *http.Server
//   if err = container.Extract(&server); err != nil {
//     // extract failed
//   }
//
// If the target type does not exist in a container or instance type building failed, Extract() returns an error.
// Use ExtractOption for modifying the behavior of this function.
func (c *Container) Extract(target interface{}, options ...ExtractOption) (err error) {
	targetValue := reflect.ValueOf(target)

	// target value type needs to be a pointer
	if targetValue.Kind() != reflect.Ptr {
		return errors.New("extract target must be a pointer")
	}

	targetValue = targetValue.Elem()

	var po = &extractOptions{
		target: targetValue,
	}

	// apply extract options
	for _, opt := range options {
		opt.apply(po)
	}

	if err = c.storage.Extract(po.name, targetValue); err != nil {
		return errors.WithStack(err)
	}

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

		prov, err := determineInstanceProvider(po)
		if err != nil {
			return errors.WithStack(err)
		}

		node := graph.NewProviderNode(po.name, prov)

		if err = c.storage.Add(node, po.implements...); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *Container) applyReplacers() (err error) {
	for _, po := range c.replacers {
		if po.provider == nil {
			return errors.New("could not provide nil")
		}

		prov, err := determineInstanceProvider(po)
		if err != nil {
			return errors.WithStack(err)
		}

		node := graph.NewProviderNode(po.name, prov)

		if err = c.storage.Replace(node, po.implements...); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

type providerOptions struct {
	name            string
	provider        interface{}
	implements      []interface{}
	includeExported bool
}

type extractOptions struct {
	name   string
	target reflect.Value
}
