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
	logger := log.New(os.Stderr, "providing", log.Ldate|log.Ltime)

	container, err := inject.New(
		inject.Provide(logger),
		inject.Provide(NewHTTPServer),
		inject.Provide(NewRouter),
	)

	if err != nil {
		logger.Fatalln(err)
	}

	var server *http.Server
	if err = container.Extract(&server); err != nil {
		logger.Fatalln(err)
	}

	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err = server.ListenAndServe(); err != nil {
			logger.Fatalln(err)
		}

		stop <- syscall.SIGTERM
	}()

}

// NewRouter
func NewRouter() *http.ServeMux {
	return &http.ServeMux{}
}

// NewHTTPServer is http server constructor
func NewHTTPServer() *http.Server {
	return &http.Server{}
}
