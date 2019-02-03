package controllers

import (
	"net/http"

	"github.com/defval/injector/testdata/mux"
)

// NewProductController ...
func NewProductController() *ProductController {
	return &ProductController{}
}

// ProductController ...
type ProductController struct {
}

// RegisterRoutes ...
func (c *ProductController) Routes(router mux.Router) {
	router.Add("/products", c.List)
}

func (c *ProductController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("product list"))
}
