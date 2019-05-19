package inject

// OPTIONS

// Option modifies container.
type Option interface{ apply(*Container) }

// ProvideOption
type ProvideOption interface{ apply(*providerOptions) }

// ApplyOption.
type ApplyOption interface{ apply(*modifierOptions) }

// PopulateOption.
type PopulateOption interface{ apply(*populateOptions) }

// CONTAINER OPTIONS.

// Provide provide dependency with options.
func Provide(provider interface{}, options ...ProvideOption) Option {
	return option(func(container *Container) {
		var po = &providerOptions{
			provider: provider,
		}

		for _, opt := range options {
			opt.apply(po)
		}

		container.providers = append(container.providers, po)
	})
}

// Apply.
func Apply(modifier interface{}, options ...ApplyOption) Option {
	return option(func(container *Container) {
		var mo = &modifierOptions{
			modifier: modifier,
		}

		for _, opt := range options {
			opt.apply(mo)
		}

		container.modifiers = append(container.modifiers, mo)
	})
}

// Package
func Package(options ...Option) Option {
	return option(func(container *Container) {
		for _, opt := range options {
			opt.apply(container)
		}
	})
}

// SetLogger.
func SetLogger(logger Logger) Option {
	return option(func(container *Container) {
		container.logger = logger
	})
}

// NopLogger.
func NopLogger() Option {
	return option(func(container *Container) {
		container.logger = &nopLogger{}
	})
}

// PROVIDE OPTIONS.

// Name
func Name(name string) ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.name = name
	})
}

// As
func As(ifaces ...interface{}) ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.implements = append(provider.implements, ifaces...)

	})
}

// POPULATE OPTIONS.

func PopulateName(name string) PopulateOption {
	return populateOption(func(populate *populateOptions) {
		populate.name = name
	})
}

// option internal
type option func(container *Container)

func (o option) apply(container *Container) { o(container) }

// provide option internal
type provideOption func(provider *providerOptions)

func (o provideOption) apply(provider *providerOptions) { o(provider) }

// apply option internal
type applyOption func(modifier *modifierOptions)

func (o applyOption) apply(modifier *modifierOptions) { o(modifier) }

// populate option internal
type populateOption func(populate *populateOptions)

func (o populateOption) apply(populate *populateOptions) { o(populate) }
