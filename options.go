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
	apply(options *modifierOptions)
}
type applyOption func(options *modifierOptions)

func (o applyOption) apply(options *modifierOptions) { o(options) }

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
type provideOption func(options *providerOptions)

func (o provideOption) apply(options *providerOptions) { o(options) }

// Name
func Name(name string) ProvideOption {
	return provideOption(func(options *providerOptions) {
		options.name = name
	})
}

// As
func As(ifaces ...interface{}) ProvideOption {
	return provideOption(func(options *providerOptions) {
		options.implements = append(options.implements, ifaces...)

	})
}
