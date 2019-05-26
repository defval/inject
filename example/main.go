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
		inject.Bundle(
			inject.Provide(NewHTTPServer),
			inject.Provide(NewServeMux, inject.As(new(http.Handler))),
			inject.Provide(&UserController{}, inject.As(new(Controller))),
			inject.Provide(&AccountController{}, inject.As(new(Controller))),
		).Namespace("test"),
	)

	if err != nil {
		log.Fatalln(err)
	}

	var server *http.Server
	if err = container.Populate(&server, inject.Namespace("test")); err != nil {
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
