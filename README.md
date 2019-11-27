<img width="312"
src="https://github.com/defval/inject/raw/master/logo.png">[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Dependency%20injection%20container%20for%20Golang&url=https://github.com/defval/inject&hashtags=golang,go,di,dependency-injection)

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://godoc.org/github.com/defval/inject)
![Release](https://img.shields.io/github/tag/defval/inject.svg?label=release&color=24B898&logo=github&style=for-the-badge)
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)
![Contributors](https://img.shields.io/github/contributors/defval/inject.svg?style=for-the-badge)

## How will dependency injection help me?

Dependency injection is one form of the broader technique of inversion
of control. It is used to increase modularity of the program and make it
extensible.

## Contents

- [Installing](#installing)
- [Documentation](#documentation)
  - [Providing](#providing)
  - [Extraction](#extraction)
  - [Interfaces and groups](#interfaces-and-groups)
- [Advanced features](#advanced-features)
  - [Named definitions](#named-definitions)
  - [Optional parameters](#optional-parameters)
  - [Prototypes](#prototypes)
  - [Cleanup](#cleanup)

## Installing

```shell
go get -u github.com/defval/inject/v2
```

## Documentation

Let's learn to use `inject` by example. We will code a simple
application that processes HTTP requests.

### Providing

To start, we will need to create two fundamental types: server and
router. We will create a simple constructors that initialize this.

Our constructors:

```go
// NewServer creates a http server with provided mux as handler.
func NewServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Handler: mux,
	}
}

// NewServeMux creates a new http serve mux.
func NewServeMux() *http.ServeMux {
	return &http.ServeMux{}
}
```

> Supported constructor signature:
>
> ```go
> func([dep1, dep2, depN]) (result, [cleanup, error])
> ```

Now let's teach a container to build these types.

```go
// Collect container parameters, build and compile container.
container := inject.New(
	inject.Provide(NewServer),  // provide http server
	inject.Provide(NewServeMux) // provide http serve mux
)
```

The function `New()` parse our constructors and compile dependency
graph.

> Container panics if it could not compile. I think that panic at the
> initialization of the application and not in runtime is usual.

> Result dependencies will be lazy-loaded. If no one requires a type
> from the container it will not be constructed.

### Extraction

We can extract the built server from the container. For this, define the
variable of extracted type and pass variable pointer to `Extract`
function.

> If extracted type not found or the process of building instance cause
> error, `Extract` return error.

If no error occurred, we can use the variable as if we had built it
yourself.

```go
// declare type variable
var server *http.Server
// extracting
err := container.Extract(&server)
if err != nil {
	// check extraction error
}

server.ListenAndServe()
```

### Interfaces and groups

Let's add some endpoints to our application.


```go
// NewAuthEndpoint creates a auth http endpoint.
func NewAuthEndpoint() *AuthEndpoint {
	return &AuthEndpoint{}
}

// AuthEndpoint is a http endpoint for auth.
type AuthEndpoint struct {}

// Login tries authenticate a user and write result using the writer.
func (a *AuthEndpoint) Login(writer http.ResponseWriter, request *http.Request) {
	// implementation
}
```

```go
// NewUserEndpoint creates a user http endpoint.
func NewUserEndpoint() *UserEndpoint {
	return &UserEndpoint{}
}

// UserEndpoint is a http endpoint for user.
type UserEndpoint struct {}

// Retrieve loads of user data and writes it using a writer.
func (e *UserEndpoint) Retrieve(writer http.ResponseWriter, request *http.Request) {
    // implementation
}
```

Change `*http.ServeMux` constructor for register endpoint routes.

```go
// NewServeMux creates a new http serve mux and register user endpoint.
func NewServeMux(auth *AuthEndpoint, users *UserEndpoint) *http.ServeMux {
	mux := &http.ServeMux{}
	mux.HandleFunc("/user", users.Retrieve)
	mux.HandleFunc("/auth", auth.Login)
	return mux
}
```

Updated container initialization code:

```go
container := inject.New(
	inject.Provide(NewServer),        // provide http server
	inject.Provide(NewServeMux)       // provide http serve mux
	// endpoints
	inject.Provide(NewUserEndpoint),  // provide user endpoint
	inject.Provide(NewAuthEndpoint),  // provide auth endpoint
)
```

Container knows that building mux requires `AuthEndpoint` and
`UserEndpoint`. And construct it for our `*http.ServeMux` on demand.

> Frequently, dependency injection is used to bind a concrete
> implementation for an interface.

Our endpoints have typical behavior. It is registering routes. Let's
create an interface for it:

```go
// Endpoint is an interface that can register its routes.
type Endpoint interface {
	RegisterRoutes(mux *http.ServeMux)
}
```

And implement it:

```go
// RegisterRoutes is a Endpoint interface implementation.
func (a *AuthEndpoint) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", a.Login)
}

// RegisterRoutes is a Endpoint interface implementation.
func (e *UserEndpoint) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user", e.Retrieve)
}
```

Now we can provide endpoint implementation as `Endpoint` interface. For
a container to know that as an implementation of `Endpoint` is necessary
to use, we use the option `inject.As()`. The argument of this option
must be a pointer to an interface like `new(Endpoint)`. This syntax may
seem strange, but I have not found a better way to specify the
interface.

```go
container := inject.New(
	inject.Provide(NewServer),        // provide http server
	inject.Provide(NewServeMux)       // provide http serve mux
	// endpoints
	inject.Provide(NewUserEndpoint, inject.As(new(Endpoint))),  // provide user endpoint
	inject.Provide(NewAuthEndpoint, inject.As(new(Endpoint))),  // provide auth endpoint
)
```

> Container groups all implementation of interface to `[]<interface>`
> group. For example, `inject.As(new(Endpoint)` automatically creates a
> group `[]Endpoint`.

We can use it in our mux. See updated code:

```go
// NewServeMux creates a new http serve mux.
func NewServeMux(endpoints []Endpoint) *http.ServeMux {
	mux := &http.ServeMux{}

	for _, endpoint := range endpoints {
		endpoint.RegisterRoutes(mux)
	}

	return mux
}
```

> If you have only one implementation of an interface, then you can use
> the interface instead of the implementation. It contributes to
> writing more testable code and not contrary to "return structs,
> accept interfaces" principle.

## Advanced features

### Named definitions

TBD

### Optional parameters

TBD

### Prototypes

TBD

### Cleanup

TBD
