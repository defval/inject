# Injector
[![Build Status](https://img.shields.io/travis/defval/injector.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/injector)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/injector.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/injector)

```go
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
        // Controllers
        injector.Bind(new(order.Repository), new(memory.OrderRepository)),
        injector.Bind(new(product.Repository), new(memory.ProductRepository)),
    
        injector.Bind(new(mux.Controller),
            new(controllers.ProductController),
            new(controllers.OrderController),
        ),
        injector.Bind(new(http.Handler), new(mux.Handler)),
    )
    
    if err != nil {
    	log.Fatal(err.Error())
    }

    var server *http.Server
    if err = container.Populate(&server); err != nil {
        log.Fatal(err.Error())
    }

    if err = server.ListenAndServe(); err != nil {
        log.Fatal(err.Error())
    }
}

```

TODO:
- Test coverage
- Verify cycles
- Bundles
- Bind type to interfaces
- Replace dependency
