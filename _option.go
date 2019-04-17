package inject

// Provide ...
func Provide(providers ...interface{}) Option {
	return option(func(injector *Injector) {
		injector.providers = append(injector.providers, providers...)
	})
}

// Bind ...
func Bind(iface interface{}, implementation interface{}) Option {
	return option(func(injector *Injector) {
		injector.bindings = append(injector.bindings, &bind{
			iface:          iface,
			implementation: implementation,
		})
	})
}

// Group
func Group(of interface{}, members ...interface{}) Option {
	return option(func(injector *Injector) {
		injector.groups = append(injector.groups, &group{
			iface:           of,
			implementations: members,
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
