package inject

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/defval/inject/internal/graph"
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

		prov, err := createProvider(po)
		if err != nil {
			return errors.WithStack(err)
		}

		node := graph.NewProviderNode(po.name, prov)

		if err = c.storage.Add(node); err != nil {
			return errors.WithStack(err)
		}

		// create group and interface alias nodes
		for _, iface := range po.implements {
			ifaceNode, err := graph.NewInterfaceNode(po.name, node, iface)

			if err != nil {
				return errors.Wrapf(err, "could not create interface alias")
			}

			if err = c.storage.Add(ifaceNode); err != nil {
				return errors.WithStack(err)
			}

			groupNode, err := c.storage.GroupNode(iface)
			if err != nil {
				return errors.WithStack(err)
			}

			if err = groupNode.Add(node); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

func (c *Container) applyReplacers() (err error) {
	for _, po := range c.replacers {
		if po.provider == nil {
			return errors.New("could not provide nil")
		}

		prov, err := createProvider(po)
		if err != nil {
			return errors.WithStack(err)
		}

		node := graph.NewProviderNode(po.name, prov)

		if err = c.storage.Replace(node); err != nil {
			return errors.WithStack(err)
		}

		// create group and interface alias nodes
		for _, iface := range po.implements {
			ifaceNode, err := graph.NewInterfaceNode(po.name, node, iface)

			if err != nil {
				return errors.Wrapf(err, "could not create interface alias")
			}

			if err = c.storage.Replace(ifaceNode); err != nil {
				return errors.WithStack(err)
			}

			groupNode, err := c.storage.GroupNode(iface)
			if err != nil {
				return errors.WithStack(err)
			}

			if err = groupNode.Replace(node); err != nil {
				return errors.WithStack(err)
			}
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
