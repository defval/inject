package controllers

import (
	"net/http"

	"github.com/defval/injector/testdata/mux"
	"github.com/defval/injector/testdata/order"
)

// NewOrderController ...
func NewOrderController(interactor *order.Interactor) *OrderController {
	return &OrderController{
		interactor: interactor,
	}
}

// OrderController ...
type OrderController struct {
	interactor *order.Interactor
}

// Routes ...
func (c *OrderController) Routes(router mux.Router) {
	router.Add("/orders", c.List)
}

// List ...
func (c *OrderController) List(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("order list"))
}
