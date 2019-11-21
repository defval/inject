package inject

// OPTIONS

// Option configures container. See inject.Provide(), inject.Bundle(), inject.Replace().
type Option interface{ apply(*Container) }

// Provide returns container option that explains how to create an instance of a type inside a container.
//
// The first argument is the provider. The provider can be constructor function, a pointer to a structure (or just
// structure) or everything else. There are some differences between these providers.
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
// Optionally, you can return a initializing error.
//
//   func NewServer(mux *http.ServeMux) (*http.Server, err error) {
//     if time.Now().Day = 1 {
//       return nil, errors.New("the server is down on the first day of a month")
//     }
//     return &http.Server{
//       Handler: mux,
//     }
//   }
//
// Other function signatures will cause error.
//
// For advanced providing use inject.dependencyProvider.
//
//   type AdminServerProvider struct {
//     inject.dependencyProvider
//
//     AdminMux http.Handler `inject:"admin"` // use named definition
//   }
//
//   func (p *AdminServerProvider) Provide() *http.Server {
//     return &http.Server{
//       Handler: p.AdminMux,
//     }
//   }
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

// Bundle group together container options.
//
//   accountBundle := inject.Bundle(
//     inject.Provide(NewAccountController),
//     inject.Provide(NewAccountRepository),
//   )
//
//   authBundle := inject.Bundle(
//     inject.Provide(NewAuthController),
//     inject.Provide(NewAuthRepository),
//   )
//
//   container, _ := New(
//     accountBundle,
//     authBundle,
//   )
func Bundle(options ...Option) Option {
	return option(func(container *Container) {
		for _, opt := range options {
			opt.apply(container)
		}
	})
}

// ProvideOption modifies default provide behavior. See inject.WithName(), inject.As(), inject.Exported().
type ProvideOption interface{ apply(*providerOptions) }

// WithName sets string identifier for provided value.
//
//   inject.Provide(&http.Server{}, inject.WithName("first"))
//   inject.Provide(&http.Server{}, inject.WithName("second"))
//
//   container.Extract(&server, inject.Name("second"))
func WithName(name string) ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.name = name
	})
}

// As specifies interfaces that implement provider instance. Provide with As() automatically checks that instance interfaces
// interface and creates slice group with it.
//
//   Provide(&http.ServerMux{}, inject.As(new(http.Handler)))
//
//   var handler http.Handler
//   container.Extract(&handler) // extract as interface
//
//   var handlers []http.Handler
//   container.Extract(&handlers) // extract group
func As(ifaces ...interface{}) ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.interfaces = append(provider.interfaces, ifaces...)

	})
}

// Prototype
func Prototype() ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.prototype = true
	})
}

// ExtractOption modifies default extract behavior. See inject.Name().
type ExtractOption interface{ apply(*extractOptions) }

// EXTRACT OPTIONS.

// Name specify definition name.
func Name(name string) ExtractOption {
	return extractOption(func(eo *extractOptions) {
		eo.name = name
	})
}

type option func(container *Container)

func (o option) apply(container *Container) { o(container) }

type provideOption func(provider *providerOptions)

func (o provideOption) apply(provider *providerOptions) { o(provider) }

type extractOption func(eo *extractOptions)

func (o extractOption) apply(eo *extractOptions) { o(eo) }
