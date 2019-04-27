package inject

// Option modify container.
type Option interface{ apply(container *Container) }
type option func(container *Container)

func (o option) apply(container *Container) { o(container) }

// Apply apply function on
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

// ApplyOption
type ApplyOption interface {
	apply(modifier *modifierOptions)
}
type applyOption func(modifier *modifierOptions)

func (o applyOption) apply(modifier *modifierOptions) { o(modifier) }

// Package
func Package(options ...Option) Option {
	return option(func(container *Container) {
		for _, opt := range options {
			opt.apply(container)
		}
	})
}

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

// ProvideOption
type ProvideOption interface {
	apply(options *providerOptions)
}
type provideOption func(provider *providerOptions)

func (o provideOption) apply(provider *providerOptions) { o(provider) }

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
