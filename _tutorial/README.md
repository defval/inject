# Tutorial

Let's code a simple application that processes HTTP requests.

## Providing

To start, we will need to create two types. We will create a simple constructors
that initialize this.

Supported constructor signature:

```go
// depN - our dependencies
// result - initialized result
// cleanup - result cleanup function
// error - initialize error
func([dep1, dep2, depN]) (result, [cleanup, error])
```

Our types:

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

Now let's teach a container to build these types.

```go
// Collect container parameters, build and compile container.
container := inject.New(
	inject.Provide(NewServer),  // provide http server
	inject.Provide(NewServeMux) // provide http serve mux
)
```

### Extraction

We can extract the built server from the container. For this, define the
variable of extracted type and pass variable pointer to `Extract`
function.

```
var server *http.Server
err := container.Extract(&server)
```

If extracted type not found or the process of building instance cause
error, `Extract` return error.

If no error occurred, we can use the variable as if we had built it
yourself.

## Interfaces and groups

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

Our endpoints have typical behavior. It is registering routes. Let's
create an interface for it.

```go
// Endpoint is an interface that can register its routes.
type Endpoint interface {
	RegisterRoutes(mux *http.ServeMux)
}
```

Now we can provide endpoint implementation as `Endpoint` interface.
Interface implementation:

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

For a container to know that as an implementation of `Endpoint` is
necessary to use, we use the option `inject.As()`. The argument of this
option must be a pointer to an interface like `new(Endpoint)`. This
syntax may seem strange, but I have not found a better way to specify
the interface.

```go
container := inject.New(
	inject.Provide(NewServer),        // provide http server
	inject.Provide(NewServeMux)       // provide http serve mux
	// endpoints
	inject.Provide(NewUserEndpoint, inject.As(new(Endpoint))),  // provide user endpoint
	inject.Provide(NewAuthEndpoint, inject.As(new(Endpoint))),  // provide auth endpoint
)
```

Container groups all implementations of `Endpoint` interface into
`[]Endpoint` group. We may use it in our mux. See updated code:

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

> If you have only one implementation of an interface then you can use
> the interface instead of the implementation. This contributes to
> writing more testable code.

