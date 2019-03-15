# Injector
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)

```go
package main

import (
	"log"
	"net/http"

	"github.com/defval/inject"
	"github.com/defval/inject/testdata/controllers"
	"github.com/defval/inject/testdata/mux"
	"github.com/defval/inject/testdata/order"
	"github.com/defval/inject/testdata/product"
	"github.com/defval/inject/testdata/storage/memory"
)

func main() {
	var container, err = inject.New(
		// HTTP
		inject.Provide(
			mux.NewHandler,
			mux.NewServer,
		),
		// Product
		inject.Provide(
			controllers.NewProductController,
			memory.NewProductRepository,
		),
		// Order
		inject.Provide(
			memory.NewOrderRepository,
			order.NewInteractor,
			controllers.NewOrderController,
		),

		// Binds
		inject.Bind(new(order.Repository), &memory.OrderRepository{}),
		inject.Bind(new(product.Repository), &memory.ProductRepository{}),
		inject.Bind(new(http.Handler), &mux.Handler{}),

		// Controllers
		inject.Group(new(mux.Controller),
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

```

TODO:
- ~~Test coverage~~
- ~~Verify cycles~~
- ~~Bundles~~
- ~~Bind type to interfaces~~
- Replace dependency
