package injector

// Provide ...
func Provide(providers ...interface{}) Option {
	return option(func(injector *Injector) {
		injector.providers = append(injector.providers, providers...)
	})
}

// Bind ...
func Bind(bindings ...interface{}) Option {
	return option(func(injector *Injector) {
		injector.bindings = append(injector.bindings, bindings)
	})
}

// Group
func Group(of interface{}, members ...interface{}) Option {
	return option(func(injector *Injector) {
		injector.groups = append(injector.groups, &group{
			of:      of,
			members: members,
		})
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
