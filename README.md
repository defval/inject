# Inject
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)

## Features

- constructor injection
- tagged struct fields injection
- inject as interface

## WIP

- Named definitions

## Usage

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

//
func main() {
	container, err := inject.New(
		inject.Provide(NewHTTPServer),
		inject.Provide(NewServeMux, inject.As(new(http.Handler))),
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
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/echo", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("echo"))
	})
}

```