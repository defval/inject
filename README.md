# Inject
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)

Dependency injection container allows you to inject dependencies
into constructors or structures without the need to have specified
each argument manually.

This container implementation inspired by [google/wire](https://github.com/google/wire),
[uber-go/fx](https://github.com/uber-go/fx) and [uber-go/dig](https://github.com/uber-go/dig).

## Installing

```shell
go get -u github.com/defval/inject
```

## Features

- inject constructor arguments
- inject tagged struct fields
- inject public struct fields
- inject as interface
- inject interface groups
- inject default value of interface group
- inject named definition into structures
- replace interface implementation
- replace provided type

## WIP

- documentation
- inject named definition into constructor

## Usage

### Provide dependency

First of all, when creating a new container, you need to describe
how to create each instance of a dependency. To do this, use the container
option `inject.Provide()`. The first argument in this function is a `provider`.
It determines how to create dependency.

Provider can be a constructor function with optional error:

```go
// dependency constructor function
func NewDependency(dependency *pkg.AnotherDependency) *pkg.Dependency {
	return &pkg.Dependency{
		dependency: dependency,
	}
}

// and with possible initialization error
func NewAnotherDependency() (*pkg.AnotherDependency, error) {
	if dependency, err = initAnotherDependency(); err != nil {
		return nil, err
	}
	
	return dependency, nil
}

// container initialization code
container, err := New(
	Provide(NewDependency),
	Provide(NewAnotherDependency)
)
```

In this case, the container knows how to create `*pkg.AnotherDependency`
and can handle an instance creation error.

Also, a provider can be a structure pointer with public fields:

```go
// package pkg
type Dependency struct {
	AnotherDependency *pkg.AnotherDependency `inject:""`
}

// container initialization code
container, err := New(
	// another providing code..
	
	// pointer to structure
    Provide(&pkg.Dependency{}),
    // or structure value
    Provide(pkg.Dependency{})
)
```

In this case, the necessity of implementing specific fields are defined
with the tag `inject`.

#### Provide hints
- [Inject named definition]()
- [Inject all exported struct fields]()

## Example

```go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/defval/inject"
)

func main() {
	container, err := inject.New(
		inject.Provide(NewHTTPServer),
		inject.Provide(NewServeMux, inject.As(new(http.Handler))),
		inject.Provide(&UserController{}, inject.As(new(Controller))),
		inject.Provide(&AccountController{}, inject.As(new(Controller))),
	)

	if err != nil {
		log.Fatalln(err)
	}

	var server *http.Server
	if err = container.Populate(&server); err != nil {
		log.Fatalln(err)
	}

	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var done = make(chan struct{})

	go func() {
		defer func() {
			close(done)
		}()

		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	<-stop

	if err = server.Close(); err != nil {
		log.Fatalln(err)
	}

	<-done
}

// NewHTTPServer
func NewHTTPServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
	}
}

// NewServeMux
func NewServeMux(controllers []Controller) *http.ServeMux {
	mux := http.NewServeMux()

	for _, ctrl := range controllers {
		ctrl.RegisterRoutes(mux)
	}

	return mux
}

// Controller
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

// UserController
type UserController struct{}

func (c *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("user"))
	})
}

// UserController
type AccountController struct{}

func (c *AccountController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/account", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("account"))
	})
}

```