# Inject
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge&logo=go&logoColor=ffffff)](https://godoc.org/github.com/defval/inject)
![Release](https://img.shields.io/github/tag/defval/inject.svg?label=release&logo=github&style=for-the-badge)
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)


Dependency injection container allows you to inject dependencies
into constructors or structures without the need to have specified
each argument manually.

This container implementation inspired by [google/wire](https://github.com/google/wire),
[uber-go/fx](https://github.com/uber-go/fx) and [uber-go/dig](https://github.com/uber-go/dig).

See [godoc](https://godoc.org/github.com/defval/inject) for feel the difference.

## Installing

```shell
go get -u github.com/defval/inject
```

## Make dependency injection easy

Define constructors:

```go
// NewHTTPHandler is a http mux constructor.
func NewHTTPServeMux() *http.ServeMux {
	return &http.ServeMux{}
}

// NewHTTPServer is a http server constructor, handler will be injected 
// by container.
func NewHTTPServer(handler *net.ServeMux) *http.Server {
	return &http.Server{
		Handler: handler,
	}
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

## Group interfaces

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

### Why do you need this?
#### Keep it testable!

Add mocks:

```go
func NewHandlerMock() *HandlerMock {
	return &HandlerMock
}
```

And save ability for mock interface implementation.

```go
func TestServer(t *testing.T) {
	handlerMock := NewHandlerMock()
	
	server := NewServer(handlerMock)
	
	// test server with mock
}
```

#### Change you code behaviour in different environments

```go
var options []inject.Option

if os.Getenv("ENV") == "dev" {
	options = append(options, inject.Provide(NewHandlerMock, inject.As(IHandler)))
} else {
	options = append(options, inject.Provide(NewServeMux, inject.As(IHandler)))
}
```

## Group your code in bundles.

```go
// ProcessingBundle responsible for processing
var ProcessingBundle = inject.Bundle(
    inject.Provide(processing.NewDispatcher),
    inject.Provide(processing.NewProvider),
    inject.Provide(processing.NewProxy),
)

// BillingBundle responsible for billing
var BillingBundle = inject.Bundle(
    inject.Provide(billing.NewInteractor),
    inject.Provide(billing.NewInvoiceRepository, inject.As(new(InvoiceRepository)))
)
```

## Replace dependencies

```go
var options []inject.Options

if os.Getenv("ENV") == "dev" {
    options = append(options, inject.Replace(billing.NewInvoiceRepositoryMock), inject.As(new(InvoiceRepository)))
}

container, err := inject.New(options...)
```

## Use named definitions

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

## Visualize dependency graph [unreleased]

Container supports `fmt.Stringer` interface. The string is a graph
description via [graphviz dot language](https://www.graphviz.org/).

This is visualization of container example.

<img src="https://github.com/defval/inject/raw/master/graph.png">