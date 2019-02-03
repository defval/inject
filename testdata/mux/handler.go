package mux

import (
	"net/http"
)

// Router
type Router interface {
	Add(pattern string, handler func(w http.ResponseWriter, r *http.Request))
}

// NewHandler
func NewHandler(controllers []Controller) *Handler {
	var h = &Handler{
		mux: http.NewServeMux(),
	}

	for _, controller := range controllers {
		controller.Routes(h)
	}

	return h
}

// Handler ...
type Handler struct {
	mux *http.ServeMux
}

// Add
func (h *Handler) Add(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	h.mux.HandleFunc(pattern, handler)
}

// ServeHTTP
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
