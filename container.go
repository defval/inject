package inject

import (
	"github.com/defval/inject/v2/di"
)

// New creates a new container with provided options.
func New(options ...Option) *Container {
	var c = &Container{
		container: di.New(),
	}

	// apply options.
	for _, opt := range options {
		opt.apply(c)
	}

	c.compile()

	return c
}

// Container is a dependency injection container.
type Container struct {
	providers []*providerOptions
	container *di.Container
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
	var po = &extractOptions{
		target: target,
	}

	// apply extract options
	for _, opt := range options {
		opt.apply(po)
	}

	return c.container.Extract(di.ExtractParams{
		Name:   po.name,
		Target: target,
	})
}

// Cleanup
func (c *Container) Cleanup() {
	c.container.Cleanup()
}

func (c *Container) compile() {
	for _, po := range c.providers {
		c.container.Provide(di.ProvideParams{
			Name:        po.name,
			Provider:    po.provider,
			Interfaces:  po.interfaces,
			IsPrototype: po.prototype,
		})
	}

	c.container.Compile()

	return
}
