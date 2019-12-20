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
	providers []di.ProvideParams
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
	var params = di.ExtractParams{
		Target: target,
	}
	// apply extract options
	for _, opt := range options {
		opt.apply(&params)
	}
	return c.container.Extract(params)
}

// Invoke invokes custom function. Dependencies of function will be resolved via container.
func (c *Container) Invoke(fn interface{}) error {
	var params = di.InvokeParams{
		Fn: fn,
	}
	return c.container.Invoke(params)
}

// Cleanup cleanup container.
func (c *Container) Cleanup() {
	c.container.Cleanup()
}

func (c *Container) compile() {
	for _, po := range c.providers {
		c.container.Provide(po)
	}
	c.container.Compile()
	return
}
