package inject_test

import (
	"log"
	"net/http"
	"os"

	"github.com/defval/inject/v2"
)

func Example() {
	// build container
	container := inject.New(
		// inject constructor
		inject.Provide(NewLogger),
		inject.Provide(NewServer),

		// inject as interface
		inject.Provide(NewRouter,
			inject.As(new(http.Handler)), // *http.Server mux implements http.Handler interface
		),

		// controller interface group
		inject.Provide(NewAccountController,
			inject.As(new(Controller)), // add AccountController to controller group
			inject.WithName("account"),
		),
		inject.Provide(NewAuthController,
			inject.As(new(Controller)), // add AuthController to controller group
			inject.WithName("auth"),
		),
	)

	// extract server from container
	var server *http.Server
	if err := container.Extract(&server); err != nil {
		panic(err)
	}

	// Output:
	// Logger loaded!
	// Create router!
	// AccountController registered!
	// AuthController registered!
	// Router created!
	// Server created!
}

// NewLogger
func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)
	defer logger.Println("Logger loaded!")

	return logger
}

// NewServer
func NewServer(logger *log.Logger, handler http.Handler) *http.Server {
	defer logger.Println("Server created!")
	return &http.Server{
		Handler: handler,
	}
}

// NewRouter
func NewRouter(logger *log.Logger, controllers []Controller) *http.ServeMux {
	logger.Println("Create router!")
	defer logger.Println("Router created!")

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

// NewAccountController
func NewAccountController(logger *log.Logger) *AccountController {
	return &AccountController{Logger: logger}
}

// RegisterRoutes
func (c *AccountController) RegisterRoutes(mux *http.ServeMux) {
	c.Logger.Println("AccountController registered!")

	// register your routes
}

// AuthController
type AuthController struct {
	Logger *log.Logger
}

// NewAuthController
func NewAuthController(logger *log.Logger) *AuthController {
	return &AuthController{Logger: logger}
}

// RegisterRoutes
func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	c.Logger.Println("AuthController registered!")

	// register your routes
}
