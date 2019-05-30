package inject_test

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/defval/inject"
)

// Example
func ExampleUsage() {
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
