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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// InjectionTestCase
type InjectionTestCase struct {
	Name    string
	Options []Option
	Error   string
}

// testCases
var testCases = []InjectionTestCase{
	{
		Name:    "PopulateNotExist",
		Options: []Option{},
		Error:   "bool not found",
	},
	{
		Name: "InjectErrorInGroup",
		Options: []Option{
			Provide(
				func(addrs []net.Addr) bool {
					return true
				},
				func() (*net.UDPAddr, error) {
					return &net.UDPAddr{}, nil
				},
				func() (*net.TCPAddr, error) {
					return nil, errors.New("dude was gone")
				},
			),
			Group(new(net.Addr), &net.UDPAddr{}, &net.TCPAddr{}),
		},
		Error: "bool: []net.Addr: *net.TCPAddr: dude was gone",
	},
	{
		Name: "GroupNilOf",
		Options: []Option{
			Provide(
				func(addrs []net.Addr) bool {
					return true
				},
				func() (*net.UDPAddr, error) {
					return &net.UDPAddr{}, nil
				},
				func() (*net.TCPAddr, error) {
					return nil, errors.New("dude was gone")
				},
			),
			Group(net.Addr(nil), &net.UDPAddr{}, &net.TCPAddr{}),
		},
		Error: "group of must be a interface pointer like new(http.Handler)",
	},
	{
		Name: "InjectError",
		Options: []Option{
			Provide(
				func(s string) bool {
					return true
				},
				func() (string, error) {
					return "", errors.New("dude was gone")
				},
			),
		},
		Error: "bool: string: dude was gone",
	},
	{
		Name: "InjectIncorrectErrorArgument",
		Options: []Option{
			Provide(
				func() (*http.Server, *http.ServeMux) {
					return nil, nil
				},
			),
		},
		Error: "injection argument must be a function with returned value and optional error",
	},
	{
		Name: "EmptyFunction",
		Options: []Option{
			Provide(
				func() {},
			),
		},
		Error: "injection argument must be a function with returned value and optional error",
	},
	{
		Name: "InjectNil",
		Options: []Option{
			Provide(
				nil,
				nil,
				nil,
			),
		},

		Error: "nil could not be injected",
	},
	{
		Name: "InjectNilInterface",
		Options: []Option{
			Provide(
				new(http.Handler),
			),
		},
		Error: "inject argument must be a function, got *http.Handler",
	},
	{
		Name: "DudeTest",
		Options: []Option{
			Provide(
				func(s string) bool {
					return s == "dude"
				},
				func() string {
					return "dude"
				},
			),
		},
	},
	{
		Name: "StringInt64",
		Options: []Option{
			Provide(
				func(s string) bool {
					return s == "value:28071990"
				},
				func(value int64) string {
					return fmt.Sprintf("%s:%d", "value", value)
				},
				func() int64 {
					return 28071990
				},
			),
		},
	},
	{
		Name: "DuplicateType",
		Options: []Option{
			Provide(
				func() string {
					return "string"
				},
				func() string {
					return "string"
				},
			),
		},
		Error: "string already injected",
	},
	{
		Name: "InjectPointer",
		Options: []Option{
			Provide(
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
			),
		},
	},
	{
		Name: "PopulateNotExistingType",
		Options: []Option{
			Provide(
				func(v bool) string {
					return ""
				},
			),
		},
		Error: "bool not found",
	},
	{
		Name: "InjectNotExistingType",
		Options: []Option{
			Provide(
				func(handler http.Handler) bool {
					_, ok := handler.(*http.ServeMux)
					return ok
				},
				func() *http.ServeMux {
					return http.NewServeMux()
				},
			),
		},
		Error: "http.Handler not found",
	},
	{
		Name: "Bind",
		Options: []Option{
			Provide(
				func(handler http.Handler) bool {
					_, ok := handler.(*http.ServeMux)
					return ok
				},
				func() *http.ServeMux {
					return http.NewServeMux()
				},
			),
			Bind(new(http.Handler), &http.ServeMux{}),
		},
	},
	{
		Name: "BindDuplicate",
		Options: []Option{

			Provide(
				func() bool {
					return true
				},
				func() *http.ServeMux {
					return http.NewServeMux()
				},
				func() *rpc.Server {
					return rpc.NewServer()
				},
			),
			Bind(new(http.Handler), &http.ServeMux{}),
			Bind(new(http.Handler), &rpc.Server{}),
		},
		Error: "http.Handler already injected",
	},
	{
		Name: "Group",
		Options: []Option{
			Provide(
				func(addrs []net.Addr) bool {
					return len(addrs) == 2
				},
				func() *net.TCPAddr {
					return &net.TCPAddr{}
				},
				func() *net.UDPAddr {
					return &net.UDPAddr{}
				},
			),
			Group(new(net.Addr), &net.TCPAddr{}, &net.UDPAddr{}),
		},
	},
	{
		Name: "BuildError",
		Options: []Option{
			Provide(
				func(server *http.Server) bool {
					return true
				},
				func() (*http.Server, error) {
					return nil, fmt.Errorf("build error")
				},
			),
		},
		Error: "bool: *http.Server: build error",
	},
}

// TestNew
func TestInjector(t *testing.T) {

	// Injection
	for _, row := range testCases {
		t.Run(row.Name, func(t *testing.T) {
			var injector, err = New(
				row.Options...,
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
				Group(new(mux.Controller),
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
