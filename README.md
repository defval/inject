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

## Installing

```shell
go get -u github.com/defval/inject
```

## Features

- Arguments injection
- Tagged and public struct fields injection
- Inject struct as interfaces
- Named definitions
- Replacing

## Example

```go
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/defval/inject"
)

func main() {
	// build container
	container, err := inject.New(
		// inject constructor
		inject.Provide(NewLogger),
		inject.Provide(NewServer),

		// inject as interface
		inject.Provide(NewRouter,
			inject.As(new(http.Handler)), // *http.Server mux implements http.Handler interface
		),

		// controller interface group
		inject.Provide(&AccountController{},
			inject.As(new(Controller)), // add AccountController to controller group
			inject.Exported(),          // inject all exported fields
		),
		inject.Provide(&AuthController{},
			inject.As(new(Controller)), // add AuthController to controller group
			inject.Exported(),          // inject all exported fields
		),
	)

	// build error
	if err != nil {
		panic(err)
	}

	// extract server from container
	var server *http.Server
	if err = container.Extract(&server); err != nil {
		panic(err)
	}

	// start server
	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// NewLogger
func NewLogger() *log.Logger {
	return log.New(os.Stderr, "", 0)
}

// NewServer
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
	}
}

// NewRouter
func NewRouter(controllers []Controller) *http.ServeMux {
	mux := &http.ServeMux{}

	for _, ctrl := range controllers {
		ctrl.RegisterRoutes(mux)
	}

	return mux
}

// Controller
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

// AccountController
type AccountController struct {
	Logger *log.Logger
}

// RegisterRoutes
func (c *AccountController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/account", func(writer http.ResponseWriter, request *http.Request) {
		c.Logger.Println("Got account request!")

		_, _ = io.WriteString(writer, "account")
	})
}

// AuthController
type AuthController struct {
	Logger *log.Logger
}

// RegisterRoutes
func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth", func(writer http.ResponseWriter, request *http.Request) {
		c.Logger.Println("Got auth request!")

		_, _ = io.WriteString(writer, "auth")
	})
}

```