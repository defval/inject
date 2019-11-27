package main

import (
	"net/http"

	"github.com/defval/inject/v2"
)

func main() {
	container := inject.New(
		inject.Provide(NewServer),
		inject.Provide(NewServeMux),
		inject.Provide(NewAuthEndpoint, inject.As(new(Endpoint))),
		inject.Provide(NewUserEndpoint, inject.As(new(Endpoint))),
	)

	var server *http.Server
	err := container.Extract(&server)
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}

// NewServer creates a http server with provided mux as handler.
func NewServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Handler: mux,
	}
}

// NewServeMux creates a new http serve mux.
func NewServeMux(endpoints []Endpoint) *http.ServeMux {
	mux := &http.ServeMux{}

	for _, endpoint := range endpoints {
		endpoint.RegisterRoutes(mux)
	}

	return mux
}

// Endpoint
type Endpoint interface {
	RegisterRoutes(mux *http.ServeMux)
}

// AuthEndpoint
type AuthEndpoint struct {
}

// NewAuthEndpoint creates a auth http endpoint.
func NewAuthEndpoint() *AuthEndpoint {
	return &AuthEndpoint{}
}

func (a *AuthEndpoint) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", a.Login)
}

// Login control user authentication.
func (a *AuthEndpoint) Login(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("LoginEndpoint"))
}

// UserEndpoint is a http endpoint for user.
type UserEndpoint struct {
}

// NewUserEndpoint
func NewUserEndpoint() *UserEndpoint {
	return &UserEndpoint{}
}

// Register is a method for Endpoint implementation.
func (e *UserEndpoint) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user", e.Retrieve)
}

func (e *UserEndpoint) Retrieve(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("UserEndpoint"))
}
