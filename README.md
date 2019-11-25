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

--------

This container implementation inspired by
[google/wire](https://github.com/google/wire),
[uber-go/fx](https://github.com/uber-go/fx) and
[uber-go/dig](https://github.com/uber-go/dig).


## Installing

```shell
go get -u github.com/defval/inject/v2
```

## Documentation

### Providing

Let's code a simple application that processes HTTP requests. For this,
we need a server and a router. We take the server and mux from the
standard library.

```go
// NewServer creates a new http server with provided handler. 
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

Now let's teach a container to build these types.

```go
// Collect container parameters, build and compile container.
container := inject.New(
	inject.Provide(NewServer),  // provide http server
	inject.Provide(NewServeMux) // provide http serve mux
)
```

### Extraction

Now, we can extract the built server from the container. For this,
define the variable of extracted type and pass variable pointer to
`Extract` function.

```
var server *http.Server
err := container.Extract(&server)
```

If extracted type not found or the process of building instance cause
error, `Extract` return error.

If no error occurred, we can use the variable as if we had built it
yourself.

### Inversion of control

Let's look on the code without dependency injection and with it.

In common case programmer control how constructor dependencies are
injected in (error checking omitted):

```go
// data layer
db := NewDatabaseConnection()
profiles := NewProfileRepository(db)
accounts := NewAccountRepository(db)
// api layer
mux := NewServeMux()
server := NewServer(mux)

profileEndpoint := NewProfileEndpoint(mux, profiles)

// func NewProfileEndpoint(mux, profiles) *ProfileEndpoint {
//     endpoint := &ProfileEndpoint{profiles}
//     mux.Get("/profile", endpoint.Retrieve)
//     mux.Post("/profile", endpoint.Create)
// }

accountEndpoint := NewAccountEndpoint(mux, accountRepository)

// func NewAccountEndpoint(mux, accounts) *AccountEndpoint {
//     endpoint := &AccountEndpoint{accounts}
//     mux.Get("/profile", profileEndpoint.Retrieve)
//     mux.Post("/profile", profileEndpoint.Create)
// }

// messaging layer
userLoggedInSubscription := NewUserLoggedInSubscription(profileRepository)
// service layer
userInteractor := NewUserInteractor(profileRepository, accountRepository)
// start
server.ListenAndServe()
```

With using of dependency injection, the container decides it.

```go
container := inject.New(
	// data layer
	inject.Provide(NewDatabaseConnection),
	inject.Provide(NewProfileRepository),
	inject.Provide(NewAccountRepository),
	// api layer
	inject.Provide(NewServeMux),
	
	// func NewServeMux(endpoints []Endpoint) *ServeMux {
	//     mux := &http.ServeMux{}
	//     for _, endpoint := range endpoints {
	//         endpoint.RegisterRoutes(mux)
	//     }
	// }
	
	inject.Provide(NewServer),
	inject.Provide(NewProfileEndpoint, inject.As(new(Endpoint))),
	inject.Provide(NewAccountEndpoint, inject.As(new(Endpoint))),
	
	// type Endpoint interface {
	//     RegisterRoutes(mux *http.ServeMux)
	// }
	
	// messaging layer
	inject.Provide(NewUserLoggedInSubscription),
	// service layer
	inject.Provide(NewUserInteractor),
)

var server *http.Server
container.Extract(&server)
server.ListenAndServe()
```


Some call it magic, but it's an inversion of control.

### Implementation

For a container to know that as an implementation of `http.Handler` it
is necessary to use `*http.ServeMux`, we use the option `inject.As()`.
The argument of this option must be a pointer to an interface like
`new(http.Handler)`. This syntax may seem strange, but I have not found
a better way to specify the interface.

```go
inject.Provide(NewServeMux, inject.As(new(http.Handler)))
```


