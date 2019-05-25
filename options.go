package inject

// OPTIONS

// Option configures container.
type Option interface {
	Namespace(name string) Option

	namespace() string
	apply(*Container)
}

// ProvideOption modifies default provide behavior.
type ProvideOption interface{ apply(*providerOptions) }

// PopulateOption modifies default populate behavior.
type PopulateOption interface{ apply(*populateOptions) }

// CONTAINER OPTIONS.

// Provide provide dependency with options.
func Provide(provider interface{}, options ...ProvideOption) Option {
	var opt = &option{
		ns: &defaultNamespace,
	}

	opt.fn = func(container *Container) {
		var po = &providerOptions{
			namespace: *opt.ns,
			provider:  provider,
		}

		for _, opt := range options {
			opt.apply(po)
		}

		container.providers = append(container.providers, po)
	}

	return opt
}

// Replace replaces provided interface by new implementation.
func Replace(provider interface{}, options ...ProvideOption) Option {
	return option{
		fn: func(container *Container) {
			var po = &providerOptions{
				provider: provider,
			}

			for _, opt := range options {
				opt.apply(po)
			}

			container.replacers = append(container.replacers, po)
		},
	}
}

// Bundle group together container options.
func Bundle(options ...Option) Option {
	return option{
		fn: func(container *Container) {
			for _, opt := range options {
				opt.apply(container)
			}
		},
	}
}

// PROVIDE OPTIONS.

// WithName sets string identifier for provided value.
func WithName(name string) ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.name = name
	})
}

// As specifies interface.
func As(ifaces ...interface{}) ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.implements = append(provider.implements, ifaces...)

	})
}

// Exported option.
func Exported() ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.injectExportedFields = true
	})
}

// POPULATE OPTIONS.

// Name ...
func Name(name string) PopulateOption {
	return populateOption(func(populate *populateOptions) {
		populate.name = name
	})
}

// Namespace
func Namespace(name string) PopulateOption {
	return populateOption(func(populate *populateOptions) {
		populate.namespace = name
	})
}

// option internal
type option struct {
	ns *string
	fn func(c *Container)
}

func (o option) Namespace(name string) Option {
	o.ns = &name
	return o
}

func (o option) namespace() string {
	return *o.ns
}

func (o option) apply(container *Container) { o.fn(container) }

// provide option internal
type provideOption func(provider *providerOptions)

func (o provideOption) apply(provider *providerOptions) { o(provider) }

// populate option internal
type populateOption func(populate *populateOptions)

func (o populateOption) apply(populate *populateOptions) { o(populate) }
