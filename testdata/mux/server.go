package mux

import (
	"net/http"
)

// NewServer ...
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
