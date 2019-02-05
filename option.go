package injector

// Provide ...
func Provide(providers ...interface{}) Option {
	return option(func(c *Injector) {
		c.providers = append(c.providers, providers...)
	})
}

// Bind ...
func Bind(bindings ...interface{}) Option {
	return option(func(c *Injector) {
		c.binders = append(c.binders, bindings)
	})
}

// Bundle ...
func Bundle(options ...Option) Option {
	return bundleOptions(options)
}

// Option ...
type Option interface {
	apply(c *Injector)
}

// option
type option func(c *Injector)

func (o option) apply(c *Injector) {
	o(c)
}

// bundle options
type bundleOptions []Option

func (o bundleOptions) apply(c *Injector) {
	for _, opt := range o {
		opt.apply(c)
	}
}
