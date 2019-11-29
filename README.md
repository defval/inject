<img width="312"
src="https://github.com/defval/inject/raw/master/logo.png">[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Dependency%20injection%20container%20for%20Golang&url=https://github.com/defval/inject&hashtags=golang,go,di,dependency-injection)

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://godoc.org/github.com/defval/inject)
![Release](https://img.shields.io/github/tag/defval/inject.svg?label=release&color=24B898&logo=github&style=for-the-badge)
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)

## How will dependency injection help me?

Dependency injection is one form of the broader technique of inversion
of control. It is used to increase modularity of the program and make it
extensible.

## Disclaimer

I use `v2` version in production, but it in a pre-release state. I need
time to finish documentation and fix possible bugs.

You can see latest `v1`
[here](https://github.com/defval/inject/tree/v1.5.2).

## Contributing

> I will be glad if you contribute to this library. I don't know much
> English, so contributing to the documentation is very meaningful to me.

[![](https://sourcerer.io/fame/defval/defval/inject/images/0)](https://sourcerer.io/fame/defval/defval/inject/links/0)[![](https://sourcerer.io/fame/defval/defval/inject/images/1)](https://sourcerer.io/fame/defval/defval/inject/links/1)[![](https://sourcerer.io/fame/defval/defval/inject/images/2)](https://sourcerer.io/fame/defval/defval/inject/links/2)[![](https://sourcerer.io/fame/defval/defval/inject/images/3)](https://sourcerer.io/fame/defval/defval/inject/links/3)[![](https://sourcerer.io/fame/defval/defval/inject/images/4)](https://sourcerer.io/fame/defval/defval/inject/links/4)[![](https://sourcerer.io/fame/defval/defval/inject/images/5)](https://sourcerer.io/fame/defval/defval/inject/links/5)[![](https://sourcerer.io/fame/defval/defval/inject/images/6)](https://sourcerer.io/fame/defval/defval/inject/links/6)[![](https://sourcerer.io/fame/defval/defval/inject/images/7)](https://sourcerer.io/fame/defval/defval/inject/links/7)

## Contents

- [Installing](#installing)
- [Getting Started](#getting-started)
- [Tutorial](#tutorial)
  - [Providing](#providing)
  - [Extraction](#extraction)
  - [Interfaces and groups](#interfaces-and-groups)
- [Inversion of control](#inversion-of-control)
- [Advanced features](#advanced-features)
  - [Named definitions](#named-definitions)
  - [Optional parameters](#optional-parameters)
  - [Parameter Bag](#parameter-bag)
  - [Prototypes](#prototypes)
  - [Cleanup](#cleanup)
- [Contributing](#contributing)

## Installing

```shell
go get -u github.com/defval/inject/v2
```

## Getting Started


## Tutorial

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

> Note that by default, the container creates instances as a singleton.
> But you can change this behaviour. See [Prototypes](#prototypes).

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

Container knows that building `*http.ServeMux` requires `*AuthEndpoint`
and `*UserEndpoint` and construct it on demand.

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
> the interface instead of the implementation. It contributes to writing
> more testable code and not contrary to "return structs, accept
> interfaces" principle.

## Inversion of control

TBD

## Advanced features

### Named definitions

In some cases you have more than one instance of one type. For example
two instances of database: master - for writing, slave - for reading.

First way is a wrapping types:

```go
// MasterDatabase provide write database access.
type MasterDatabase struct {
	*Database
}

// SlaveDatabase provide read database access.
type SlaveDatabase struct {
	*Database
}
```

Second way is a using named definitions with `inject.WithName()` provide
option:

```go
// provide master database
inject.Provide(NewMasterDatabase, inject.WithName("master"))
// provide slave database
inject.Provide(NewSlaveDatabase, inject.WithName("slave"))
```

If you need to extract it from container use `inject.Name()` extract
option.

```go
var db *Database
container.Extract(&db, inject.Name("master"))
```

If you need to provide named definition in other constructor use
`di.Parameter` with embedding.

```go
// ServiceParameters
type ServiceParameters struct {
	di.Parameter
	
	// use `di` tag for the container to know that field need to be injected.
	MasterDatabase *Database `di:"master"`
	SlaveDatabase *Database  `di:"slave"`
}

// NewService creates new service with provided parameters.
func NewService(parameters ServiceParameters) *Service {
	return &Service{
		MasterDatabase:  parameters.MasterDatabase,
		SlaveDatabase: parameters.SlaveDatabase,
	}
}
```

### Optional parameters

Also `di.Parameter` provide ability to skip dependency if it not exists
in container.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Parameter
	
	Logger *Logger `di:"optional"`
}
```

> Constructors that declare dependencies as optional must handle the
> case of those dependencies being absent.

You can use naming and optional together.

```go
// ServiceParameter
type ServiceParameter struct {
	di.Parameter
	
	StdOutLogger *Logger `di:"stdout"`
	FileLogger   *Logger `di:"file,optional"`
}
```

### Parameter Bag

If you need to specify some parameters on definition level you can use
`inject.ParameterBag` provide option. This is a `map[string]interface{}`
that transforms to `di.ParameterBag` type.

```go
// Provide server with parameter bag
inject.Provide(NewServer, inject.ParameterBag{
	"addr": ":8080",
})

// NewServer create a server with provided parameter bag. Note: use di.ParameterBag type.
// Not inject.ParameterBag.
func NewServer(pb di.ParameterBag) *http.Server {
	return &http.Server{
		Addr: pb.RequireString("addr"),
	}
}
```

### Prototypes

If you want to create a new instance on each extraction use
`inject.Prototype()` provide option.

```go
inject.Provide(NewRequestContext, inject.Prototype())
```

> todo: real use case

### Cleanup

If a provider creates a value that needs to be cleaned up, then it can
return a closure to clean up the resource.

```go
func NewFile(log Logger, path Path) (*os.File, func(), error) {
    f, err := os.Open(string(path))
    if err != nil {
        return nil, nil, err
    }
    cleanup := func() {
        if err := f.Close(); err != nil {
            log.Log(err)
        }
    }
    return f, cleanup, nil
}
```

After `container.Cleanup()` call, it iterate over instances and call
cleanup function if it exists.

```go
container := inject.New(
	// ...
    inject.Provide(NewFile),
)

// do something
container.Cleanup() // file was closed
```

> Cleanup now work incorrectly with prototype providers.
