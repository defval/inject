package ding

// Provide ...
func Provide(constructors ...interface{}) Option {
	return option(func(c *Container) {

	})
}

// Populate ...
func Populate(targets ...interface{}) Option {
	return option(func(c *Container) {

	})
}

// Bind ...
func Bind(pseudonims ...interface{}) Option {
	return option(func(c *Container) {})
}

// Bundle ...
func Bundle(options ...Option) Option {
	return bundleOptions(options)
}

// Option ...
type Option interface {
	apply(c *Container)
}

// option
type option func(c *Container)

func (o option) apply(c *Container) {
	o(c)
}

// bundle options
type bundleOptions []Option

func (o bundleOptions) apply(c *Container) {
	for _, opt := range o {
		opt.apply(c)
	}
}
