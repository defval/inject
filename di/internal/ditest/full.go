package ditest

import (
	"log"
	"net/http"
	"os"
)

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
