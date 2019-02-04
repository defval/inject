package injector

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/defval/injector/testdata/controllers"
	"github.com/defval/injector/testdata/mux"
	"github.com/defval/injector/testdata/order"
	"github.com/defval/injector/testdata/product"
	"github.com/defval/injector/testdata/storage/memory"
	"github.com/stretchr/testify/assert"
)

// InjectionTestCase
type InjectionTestCase struct {
	Name       string
	Injections []interface{}
	Bindings   [][]interface{}
	Error      string
}

// injectionTable
var injectionTable = []InjectionTestCase{
	{
		Name: "DudeTest",
		Injections: []interface{}{
			func(s string) bool {
				return s == "dude"
			},
			func() string {
				return "dude"
			},
		},
	},
	{
		Name: "StringInt64",
		Injections: []interface{}{
			func(s string) bool {
				return s == "value:28071990"
			},
			func(value int64) string {
				return fmt.Sprintf("%s:%d", "value", value)
			},
			func() int64 {
				return 28071990
			},
		},
	},
	{
		Name: "DuplicateType",
		Injections: []interface{}{
			func() string {
				return "string"
			},
			func() string {
				return "string"
			},
		},
		Error: "string already injected",
	},
	{
		Name: "InjectStructPointer",
		Injections: []interface{}{
			func(server *http.Server, handler http.Handler) bool {
				return server.Handler == handler
			},
			func(handler http.Handler) *http.Server {
				return &http.Server{
					Handler: handler,
				}
			},
			func() http.Handler {
				return http.NewServeMux()
			},
		},
	},
	{
		Name: "PopulateNotExistingType",
		Injections: []interface{}{
			func(v bool) string {
				return ""
			},
		},
		Error: "bool not found",
	},
	// {
	// 	Name: "InjectError",
	// 	Injections: []interface{}{
	// 		func(s string) bool {
	// 			return true
	// 		},
	// 		func() (string, error) {
	// 			return "", errors.New("dude was gone")
	// 		},
	// 	},
	// 	Error: "string: dude was gone",
	// },
}

// TestNew
func TestNew(t *testing.T) {

	// Injection
	t.Run("Injection", func(t *testing.T) {
		for _, row := range injectionTable {
			t.Run(row.Name, func(t *testing.T) {
				var injector, err = New(
					Provide(row.Injections...),
				)

				var result bool
				if row.Error == "" {
					if err = injector.Populate(&result); err != nil {
						assert.FailNow(t, err.Error())
					}

					assert.EqualValues(t, true, result)
				} else {
					assert.EqualError(t, err, row.Error)
				}

			})
		}
	})

	t.Run("Injection", func(t *testing.T) {
		var container, err = New(
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

		if err != nil {
			assert.FailNow(t, err.Error())
		}

		var server *http.Server
		if err = container.Populate(&server); err != nil {
			assert.Fail(t, err.Error())
		}

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
