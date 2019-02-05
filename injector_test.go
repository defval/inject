package injector

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
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
		Name: "InjectPointer",
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
	{
		Name: "InjectNotExistingType",
		Injections: []interface{}{
			func(handler http.Handler) bool {
				_, ok := handler.(*http.ServeMux)
				return ok
			},
			func() *http.ServeMux {
				return http.NewServeMux()
			},
		},
		Error: "http.Handler not found",
	},
	{
		Name: "Bind",
		Injections: []interface{}{
			func(handler http.Handler) bool {
				_, ok := handler.(*http.ServeMux)
				return ok
			},
			func() *http.ServeMux {
				return http.NewServeMux()
			},
		},
		Bindings: [][]interface{}{
			{new(http.Handler), &http.ServeMux{}},
		},
	},
	{
		Name: "BindDuplicate",
		Injections: []interface{}{
			func() bool {
				return true
			},
			func() *http.ServeMux {
				return http.NewServeMux()
			},
			func() *rpc.Server {
				return rpc.NewServer()
			},
		},
		Bindings: [][]interface{}{
			{new(http.Handler), &http.ServeMux{}},
			{new(http.Handler), &rpc.Server{}},
		},
		Error: "http.Handler already injected",
	},
	{
		Name: "BindGroup",
		Injections: []interface{}{
			func(addrs []net.Addr) bool {
				return len(addrs) == 2
			},
			func() *net.TCPAddr {
				return &net.TCPAddr{}
			},
			func() *net.UDPAddr {
				return &net.UDPAddr{}
			},
		},
		Bindings: [][]interface{}{
			{new(net.Addr), &net.TCPAddr{}, &net.UDPAddr{}},
		},
	},
	{
		Name: "BuildError",
		Injections: []interface{}{
			func(server *http.Server) bool {
				return true
			},
			func() (*http.Server, error) {
				return nil, fmt.Errorf("build error")
			},
		},
		Error: "*http.Server build error: build error",
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
func TestInjector(t *testing.T) {

	// Injection
	for _, row := range injectionTable {
		t.Run(row.Name, func(t *testing.T) {
			var options []Option

			options = append(options, Provide(row.Injections...))

			for _, bindingSet := range row.Bindings {
				options = append(options, Bind(bindingSet...))
			}

			var injector, err = New(
				options...,
			)

			if err != nil && row.Error == "" {
				assert.FailNow(t, err.Error())
			}

			var result bool
			if row.Error == "" {
				if err = injector.Populate(&result); err != nil {
					assert.FailNow(t, err.Error())
				}

				assert.EqualValues(t, true, result)
			} else {
				if err != nil {
					assert.EqualError(t, err, row.Error)
				} else {
					assert.EqualError(t, injector.Populate(&result), row.Error)
				}
			}

		})
	}
}

func TestApp(t *testing.T) {
	t.Run("DummyApplication", func(t *testing.T) {
		var container, err = New(
			// HTTP
			Bundle(
				Provide(
					mux.NewHandler,
					mux.NewServer,
				),
				Bind(new(http.Handler), new(mux.Handler)),

				// Controllers
				Bind(new(mux.Controller),
					new(controllers.ProductController),
					new(controllers.OrderController),
				),
			),

			// Product
			Bundle(
				Provide(
					controllers.NewProductController,
					memory.NewProductRepository,
				),
				Bind(new(product.Repository), new(memory.ProductRepository)),
			),

			// Order
			Bundle(
				Provide(
					memory.NewOrderRepository,
					order.NewInteractor,
					controllers.NewOrderController,
				),
				Bind(new(order.Repository), new(memory.OrderRepository)),
			),
		)

		if err != nil {
			assert.FailNow(t, err.Error())
		}

		var server *http.Server
		if err = container.Populate(&server); err != nil {
			assert.FailNow(t, err.Error())
		}

		assert.NotNil(t, server)

		var cs []mux.Controller
		if err = container.Populate(&cs); err != nil {
			assert.FailNow(t, err.Error())
		}

		assert.Len(t, cs, 2)

		var products order.Repository
		if err = container.Populate(&products); err != nil {
			assert.FailNow(t, err.Error())
		}

		assert.NotNil(t, products)
	})
}
