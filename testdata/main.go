package main

import (
	"log"
	"net/http"

	"github.com/defval/injector"
	"github.com/defval/injector/testdata/controllers"
	"github.com/defval/injector/testdata/mux"
	"github.com/defval/injector/testdata/order"
	"github.com/defval/injector/testdata/product"
	"github.com/defval/injector/testdata/storage/memory"
)

func main() {
	var container, err = injector.New(
		// HTTP
		injector.Provide(
			mux.NewHandler,
			mux.NewServer,
		),
		// Product
		injector.Provide(
			controllers.NewProductController,
			memory.NewProductRepository,
		),
		// Order
		injector.Provide(
			memory.NewOrderRepository,
			order.NewInteractor,
			controllers.NewOrderController,
		),

		// Binds
		injector.Bind(new(order.Repository), &memory.OrderRepository{}),
		injector.Bind(new(product.Repository), &memory.ProductRepository{}),
		injector.Bind(new(http.Handler), &mux.Handler{}),

		// Controllers
		injector.Group(new(mux.Controller),
			&controllers.ProductController{},
			&controllers.OrderController{},
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	var server *http.Server
	if err = container.Populate(&server); err != nil {
		log.Fatal(err)
	}

	log.Println("Successful run")

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
