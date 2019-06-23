<img width="312" src="https://github.com/defval/inject/raw/master/logo.png">[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Dependency%20injection%20container%20for%20Golang&url=https://github.com/defval/inject&hashtags=golang,go,di,dependency-injection)

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://godoc.org/github.com/defval/inject)
![Release](https://img.shields.io/github/tag/defval/inject.svg?label=release&color=24B898&logo=github&style=for-the-badge)
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)
![Contributors](https://img.shields.io/github/contributors/defval/inject.svg?style=for-the-badge)


Dependency injection container allows you to inject dependencies
into constructors or structures without the need to have specified
each argument manually.

This container implementation inspired by [google/wire](https://github.com/google/wire),
[uber-go/fx](https://github.com/uber-go/fx) and [uber-go/dig](https://github.com/uber-go/dig).

See [godoc](https://godoc.org/github.com/defval/inject) for feel the difference.

## Contents

- [Installing](#installing)
- [Type injection](#type-injection)
- [Groups](#groups)
- [Bundles](#bundles)
- [Replace](#replace)
- [Named definitions](#named-definitions)
- [Visualize](#visualize-graphviz)

## Installing

```shell
go get -u github.com/defval/inject
```

## Type injection

Define constructors:

```go
// NewHTTPHandler is a http mux constructor.
func NewHTTPServeMux() *http.ServeMux {
	return &http.ServeMux{}
}

// NewHTTPServer is a http server constructor, handler will be injected 
// by container. If environment variable `STATUS == "stoped"` extract
// server cause error.
func NewHTTPServer(handler *net.ServeMux) (*http.Server, error) {
	if os.Getenv("STATUS") == "stopped" {
		return nil, errors.New("server stoped")
	}
	
	return &http.Server{
		Handler: handler,
	}, nil
}
```

Build container and extract values:

```go
// build container
container, err := inject.New(
    inject.Provide(NewHTTPServeMux), // provide mux
    inject.Provide(NewHTTPServer), // provide server
)

// don't forget to handle errors © golang

// define variable for *http.Server
var server *http.Server

// extract into this variable
container.Extract(&server)

// use it!
server.ListenAndServe()
```

## Groups

When you have two or more implementations of same interface:

```go
// NewUserController
func NewUserController() *UserController {
	return &UserController{}
}

// NewPostController
func NewPostController() *PostController {
	return &PostController()
}

// Controller
type Controller interface {
	RegisterRoutes()
}
```

Group it!

```go
// IController is a java style interface alias =D
// inject.As(new(Controller)) looks worse in readme.
var IController = new(Controller)

container, err := inject.New(
	inject.Provide(NewUserController, inject.As(IController)),
	inject.Provide(NewPostController, inject.As(IController)),
)

var controllers []Controller
// extract all controllers
container.Extract(&controllers)

// and do something!!!
for _, ctrl := range controllers {
	ctrl.RegisterRoutes()
}
```

## Return structs, accept interfaces!

Bind implementations as interfaces:

```go
// NewHandler is a http mux constructor. Returns concrete
// implementation - *http.ServeMux.
func NewServeMux() *http.ServeMux {
	return &http.ServeMux{}
}

// NewServer is a http server constructor. Needs handler for 
// working.
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
	}
}
```

Provide concrete implementation as interface:

```go
var IHandler = new(http.Handler)

container, err := inject.New(
    inject.Provide(NewServeMux, inject.As(IHandler)),
    inject.Provide(NewServer),
)

var handler http.Handler
container.Extract(&handler) // *http.ServeMux will be extracted

var server *http.Server
container.Extract(&server) // server.Handler is *http.ServeMux
```

## Bundles

```go
// ProcessingBundle responsible for processing
var ProcessingBundle = inject.Bundle(
    inject.Provide(processing.NewDispatcher),
    inject.Provide(processing.NewProvider),
    inject.Provide(processing.NewProxy, inject.As(IProxy)),
)

// BillingBundle responsible for billing
var BillingBundle = inject.Bundle(
    inject.Provide(billing.NewInteractor),
    inject.Provide(billing.NewInvoiceRepository, inject.As(new(InvoiceRepository)))
)
```

And test each one separately.

```go
func TestProcessingBundle(t *testing.T) {
    bundle, err := inject.New(
        ProcessingBundle,
        inject.Replace(processing.NewDevProxy, inject.As(IProxy)),
    )
    
    var dispatcher *processing.Dispatcher
    container.Extract(&dispatcher)
    
    dispatcher.Dispatch(ctx context.Context, thing)
}
```

## Replace

```go
var options []inject.Options

if os.Getenv("ENV") == "dev" {
    options = append(options, inject.Replace(billing.NewInvoiceRepositoryMock), inject.As(new(InvoiceRepository)))
}

container, err := inject.New(options...)
```

## Named definitions

```go
container, err := inject.New{
	inject.Provide(NewDefaultServer, inject.WithName("default")),
	inject.Provide(NewAdminServer, inject.WithName("admin")),
}

var defaultServer *http.Server
var adminServer *http.Server

container.Extract(&defaultServer, inject.Name("default"))
container.Extract(&adminServer, inject.Name("admin"))
```

Or with struct provider:

```go
// Application
type Application struct {
    Server *http.Server `inject:"default"`
    AdminServer *http.Server `inject:"admin"`
}
```

```go
container, err := inject.New(
    inject.Provide(NewDefaultServer, inject.WithName("default")), 
    inject.Provide(NewAdminServer, inject.WithName("admin")),
    inject.Provide(&Application)
)
```

If you don't like tags as much as I do, then look to
`inject.Exported()` provide option.

## Use combined provider

For advanced providing use combined provider. It's both - struct and constructor providers.

```go
// ServerProvider
type ServerProvider struct {
	Mux *http.Server `inject:"dude_mux"`
}

// Provide is a container predefined constructor function for *http.Server.
func (p *ServerProvider) Provide() *http.Server {
	return &http.Server{
		Handler: p.Mux,
	}
}
```

## Visualize ([Graphviz](https://www.graphviz.org/))

Write visualization into `io.Writer`. Check out result on <a href="https://dreampuf.github.io/GraphvizOnline" target="_blank">graphviz online tool!</a>

```go
    // visualization data target
    buffer := &bytes.Buffer{}
    
    // write container visualization
    container.WriteTo(buffer)
```

This is visualization of container example.

<img src="https://github.com/defval/inject/raw/master/graph.png">

