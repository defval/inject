package inject

import "github.com/defval/inject/v2/di"

// OPTIONS

// Option configures container. See inject.Provide(), inject.Bundle(), inject.Replace().
type Option interface {
	apply(*Container)
}

// Provide returns container option that explains how to create an instance of a type inside a container.
//
// The first argument is the constructor function. A constructor is a function that creates an instance of the required
// type. It can take an unlimited number of arguments needed to create an instance - the first returned value.
//
//   func NewServer(mux *http.ServeMux) *http.Server {
//     return &http.Server{
//       Handle: mux,
//     }
//   }
//
// Optionally, you can return a cleanup function and initializing error.
//
//   func NewServer(mux *http.ServeMux) (*http.Server, cleanup func(), err error) {
//     if time.Now().Day = 1 {
//       return nil, nil, errors.New("the server is down on the first day of a month")
//     }
//
//     server := &http.Server{
//       Handler: mux,
//     }
//
//     cleanup := func() {
//       _ = server.Close()
//     }
//
//     return server, cleanup, nil
//   }
//
// Other function signatures will cause error.
func Provide(provider interface{}, options ...ProvideOption) Option {
	return option(func(container *Container) {
		var po = di.ProvideParams{
			Provider:   provider,
			Parameters: map[string]interface{}{},
		}

		for _, opt := range options {
			opt.apply(&po)
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

// ProvideOption modifies default provide behavior. See inject.WithName(), inject.As(), inject.Prototype().
type ProvideOption interface {
	apply(params *di.ProvideParams)
}

// WithName sets string identifier for provided value.
//
//   inject.Provide(&http.Server{}, inject.WithName("first"))
//   inject.Provide(&http.Server{}, inject.WithName("second"))
//
//   container.Extract(&server, inject.Name("second"))
func WithName(name string) ProvideOption {
	return provideOption(func(provider *di.ProvideParams) {
		provider.Name = name
	})
}

// As specifies interfaces that implement provider instance. Provide with As() automatically checks that constructor
// result implements interface and creates slice group with it.
//
//   Provide(&http.ServerMux{}, inject.As(new(http.Handler)))
//
//   var handler http.Handler
//   container.Extract(&handler) // extract as interface
//
//   var handlers []http.Handler
//   container.Extract(&handlers) // extract group
func As(ifaces ...interface{}) ProvideOption {
	return provideOption(func(provider *di.ProvideParams) {
		provider.Interfaces = append(provider.Interfaces, ifaces...)

	})
}

// Prototype modifies Provide() behavior. By default, each type resolves as a singleton. This option sets that
// each type resolving creates a new instance of the type.
//
//   Provide(&http.Server{], inject.Prototype())
//
//   var server1 *http.Server
//   var server2 *http.Server
//   container.Extract(&server1, &server2)
func Prototype() ProvideOption {
	return provideOption(func(provider *di.ProvideParams) {
		provider.IsPrototype = true
	})
}

// ParameterBag is a provider parameter bag. It stores a construction parameters. It is a alternative way to
// configure type.
//
//   inject.Provide(NewServer, inject.ParameterBag{
//     "addr": ":8080",
//   })
//
//   NewServer(pb inject.ParameterBag) *http.Server {
//     return &http.Server{
//       Addr: pb.RequireString("addr"),
//     }
//   }
type ParameterBag map[string]interface{}

func (p ParameterBag) apply(provider *di.ProvideParams) {
	for k, v := range p {
		provider.Parameters[k] = v
	}
}

// ExtractOption modifies default extract behavior. See inject.Name().
type ExtractOption interface {
	apply(params *di.ExtractParams)
}

// EXTRACT OPTIONS.

// Name specify definition name.
func Name(name string) ExtractOption {
	return extractOption(func(eo *di.ExtractParams) {
		eo.Name = name
	})
}

type option func(container *Container)

func (o option) apply(container *Container) { o(container) }

type provideOption func(provider *di.ProvideParams)

func (o provideOption) apply(provider *di.ProvideParams) { o(provider) }

type extractOption func(eo *di.ExtractParams)

func (o extractOption) apply(eo *di.ExtractParams) { o(eo) }

type extractOptions struct {
	name   string
	target interface{}
}
