package inject

// OPTIONS

// Option configures container. See inject.Provide(), inject.Bundle(), inject.Replace().
// todo: Namespace
type Option interface {
	Namespace(name string) Option

	namespace() string
	apply(*Container)
}

// ProvideOption modifies default provide behavior. See inject.WithName(), inject.As(), inject.Exported().
type ProvideOption interface{ apply(*providerOptions) }

// PopulateOption modifies default populate behavior. See inject.Name(), inject.Namespace().
type PopulateOption interface{ apply(*populateOptions) }

// Provide returns container option that explains to it how to create an instance of a type inside a container.
//
// The first argument is the provider. A provider can be a constructor function, a pointer to a structure
// (or just a structure) and everything else. There are some differences between these providers.
//
// A constructor function is a function that creates an instance of the required type. It can take an unlimited
// number of arguments needed to create an instance - the first returned value.
//
//   func NewServer(mux *http.ServeMux) *http.Server {
//     return &http.Server{
//       Handle: mux,
//     }
//   }
//
// Optionally, you can return an error to create an instance.
//
//   func NewServer(mux *http.ServeMux) (*http.Server, err error) {
// 	   if time.Now().Day = 1 {
// 			return nil, errors.New("the server is down on the first day of a month")
// 	   }
//     return &http.Server{
//       Handler: mux,
//     }
//   }
//
func Provide(provider interface{}, options ...ProvideOption) Option {
	var opt = &option{}

	opt.fn = func(container *Container) {
		var po = &providerOptions{
			provider:  provider,
			namespace: opt.namespace(),
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
	var opt = &option{}

	opt.fn = func(container *Container) {
		var po = &providerOptions{
			provider:  provider,
			namespace: opt.namespace(),
		}

		for _, opt := range options {
			opt.apply(po)
		}

		container.replacers = append(container.replacers, po)
	}

	return opt
}

// Bundle group together container options.
func Bundle(options ...Option) Option {
	var bundleOpt = &option{
		ns: "",
	}

	bundleOpt.fn = func(container *Container) {
		for _, opt := range options {
			opt.Namespace(bundleOpt.ns)
			opt.apply(container)
		}
	}

	return bundleOpt
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
	ns string
	fn func(c *Container)
}

func (o *option) Namespace(name string) Option {
	o.ns = name

	return o
}

func (o *option) namespace() string {
	return o.ns
}

func (o *option) apply(container *Container) { o.fn(container) }

type provideOption func(provider *providerOptions)

func (o provideOption) apply(provider *providerOptions) { o(provider) }

type populateOption func(populate *populateOptions)

func (o populateOption) apply(populate *populateOptions) { o(populate) }
