package injector

import (
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/defval/injector/testdata/controllers"
	"github.com/defval/injector/testdata/mux"
	"github.com/defval/injector/testdata/order"
	"github.com/defval/injector/testdata/product"
	"github.com/defval/injector/testdata/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Injection", func(t *testing.T) {
		var err error

		var container = New(
			// HTTP
			Provide(
				mux.NewHandler,
				mux.NewServer,
			),
			// Product
			Provide(
				controllers.NewProductController,
				memory.NewProductRepository,
			),
			// Order
			Provide(
				memory.NewOrderRepository,
				order.NewInteractor,
				controllers.NewOrderController,
			),
			// Controllers
			Bind(new(order.Repository), new(memory.OrderRepository)),
			Bind(new(product.Repository), new(memory.ProductRepository)),

			Bind(new(mux.Controller),
				new(controllers.ProductController),
				new(controllers.OrderController),
			),
			Bind(new(http.Handler), new(mux.Handler)),
		)

		if err := container.Error(); err != nil {
			assert.FailNow(t, "container build error")
		}

		var server *http.Server
		if err = container.Populate(&server); err != nil {
			assert.Fail(t, err.Error())
		}

		spew.Dump(server)

		assert.NotNil(t, server)
	})
	//
	// t.Run("PopulateMultipleInterfaceImplementation", func(t *testing.T) {
	// 	type stringer interface {
	// 		s() string
	// 	}
	//
	// 	var stringers []stringer
	//
	// 	var container = New(
	// 		Provide(
	// 			func() FirstStringer {
	// 				return FirstStringer("first")
	// 			},
	// 			func() SecondStringer {
	// 				return SecondStringer("first")
	// 			},
	// 		),
	// 		Bind(new(stringer), new(FirstStringer), new(SecondStringer)),
	// 		Populate(&stringers),
	// 	)
	//
	// 	// Err
	// 	if err := container.Error(); err != nil {
	// 		assert.FailNow(t, "%s", err)
	// 	}
	//
	// 	assert.Len(t, stringers, 2)
	//
	// 	for _, s := range stringers {
	// 		assert.Implements(t, new(stringer), s)
	// 	}
	// })
}
