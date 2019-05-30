package inject

// OPTIONS

// Option configures container. See inject.Provide(), inject.Bundle(), inject.Replace().
type Option interface{ apply(*Container) }

// ProvideOption modifies default provide behavior. See inject.WithName(), inject.As(), inject.Exported().
type ProvideOption interface{ apply(*providerOptions) }

// ExtractOption modifies default extract behavior. See inject.Name().
type ExtractOption interface{ apply(*extractOptions) }

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
//     if time.Now().Day = 1 {
//       return nil, errors.New("the server is down on the first day of a month")
// 	   }
//     return &http.Server{
//       Handler: mux,
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

// Replace replaces a already provided definition to another one.
//
// You may replace definition with one result type.
//
//   inject.New(
//     inject.Provide(&http.Server{Addr: ":80"}),
//     inject.Replace(&http.Server{Addr: ":8080"}),
//   )
//
// Or with one implemented interface.
//
//   inject.New(
//     inject.Provide(&http.ServeMux{}, inject.As(new(http.Handler))),
//     inject.Replace(&mux.AnotherMux{}, inject.As(new(http.Handler))),
//   )
func Replace(provider interface{}, options ...ProvideOption) Option {
	return option(func(container *Container) {
		var po = &providerOptions{
			provider: provider,
		}

		for _, opt := range options {
			opt.apply(po)
		}

		container.replacers = append(container.replacers, po)
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
//   container, _ := New(accountBundle, authBundle)
func Bundle(options ...Option) Option {
	return option(func(container *Container) {
		for _, opt := range options {
			opt.apply(container)
		}
	})
}

// PROVIDE OPTIONS.

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

// As specifies interface that implement provider instance.
func As(ifaces ...interface{}) ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.implements = append(provider.implements, ifaces...)

	})
}

// Exported indicates that all public fields of the structure should be injected.
//
//   type AccountController struct {
//     Accounts AccountRepository // will be injected without tag 'inject'
//   }
//
//   inject.Provide(NewAccountRepository, inject.As(new(AccountRepository)))
//   inject.Provide(&AccountController{}, inject.Exported())
func Exported() ProvideOption {
	return provideOption(func(provider *providerOptions) {
		provider.injectExportedFields = true
	})
}

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
