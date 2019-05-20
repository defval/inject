# Inject
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)

Dependency injection container allows you to inject dependencies into constructors or
structures without the need to having specify each constructor argument manually.

## Injection features

- inject result of constructor
- inject tagged struct fields
- inject public struct fields
- inject as interface
- inject interface groups
- inject default value of interface group
- inject named definition

## WIP

- documentation
- inject named definition into constructor

## Full example

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
		inject.Apply(RegisterRoutes),
	)

	if err != nil {
		log.Fatalln(err)
	}

	var server *http.Server
	if err = container.Populate(&server); err != nil {
		log.Fatalln(err)
	}

	var stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	<-stop

	if err = server.Close(); err != nil {
		log.Fatalln(err)
	}
}

// NewHTTPServer
func NewHTTPServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
	}
}

// NewServeMux
func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}

// RegisterRoutes
func RegisterRoutes(mux *http.ServeMux, controllers []Controller) {
	for _, ctrl := range controllers {
		ctrl.RegisterRoutes(mux)
	}
}

// Controller
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

// UserController
type UserController struct {
}

func (c *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("user"))
	})
}

// UserController
type AccountController struct {
}

func (c *AccountController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/account", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("account"))
	})
}


```