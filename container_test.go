package inject_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/defval/inject"
)

func eqPtr(t *testing.T, expected interface{}, actual interface{}) {
	require.Equal(t, fmt.Sprintf("%p", expected), fmt.Sprintf("%p", actual))
}

func TestContainer_ProvideConstructor(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide("string"),
			inject.Provide(func(s string) bool {
				require.Equal(t, "string", s)
				return true
			}),
		)

		require.NoError(t, err)

		var result string
		require.NoError(t, container.Extract(&result))
		require.Equal(t, "string", result)

		var b bool
		require.NoError(t, container.Extract(&b))
	})

	t.Run("struct", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(struct{}{}),
		)

		require.NoError(t, err)
	})

	t.Run("slice", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide([]int64{32, 30, 31}),
		)

		require.NoError(t, err)

		var result []int64
		require.NoError(t, container.Extract(&result))
		require.Len(t, result, 3)
	})

	t.Run("chan", func(t *testing.T) {
		var ch = make(chan struct{})

		container, err := inject.New(
			inject.Provide(ch),
			inject.Provide(func(ch chan struct{}) bool {
				close(ch)
				return true
			}),
		)

		require.NoError(t, err)

		var b bool
		require.NoError(t, container.Extract(&b))
		_, more := <-ch
		require.False(t, more)
	})

	t.Run("map", func(t *testing.T) {
		var m = map[string]string{"test": "test"}

		container, err := inject.New(
			inject.Provide(m),
			inject.Provide(func(arg map[string]string) bool {
				require.Equal(t, m, arg)
				return true
			}),
		)

		require.NoError(t, err)

		var b bool
		require.NoError(t, container.Extract(&b))
		require.True(t, b)
	})

	t.Run("constructors with dependency without errors", func(t *testing.T) {
		var server = &http.Server{}
		var mux = &http.ServeMux{}

		container, err := inject.New(
			inject.Provide(func() *http.ServeMux {
				return mux
			}),
			inject.Provide(func(mux *http.ServeMux) *http.Server {
				server.Handler = mux
				return server
			}),
			inject.Provide(func(s *http.Server) bool {
				eqPtr(t, server, s)
				return true
			}),
		)

		require.NoError(t, err)

		var extractedServer *http.Server
		require.NoError(t, container.Extract(&extractedServer))
		require.NotNil(t, server)

		eqPtr(t, extractedServer, server)
		eqPtr(t, mux, server.Handler)

		var r bool
		require.NoError(t, container.Extract(&r))
	})

	t.Run("constructors with dependency with nil errors", func(t *testing.T) {
		var server = &http.Server{}
		var mux = &http.ServeMux{}

		container, err := inject.New(
			inject.Provide(func() (*http.ServeMux, error) {
				return mux, nil
			}),
			inject.Provide(func(mux *http.ServeMux) (*http.Server, error) {
				server.Handler = mux
				return server, nil
			}),
			inject.Provide(func(s *http.Server) bool {
				eqPtr(t, server, s)
				return true
			}),
		)

		require.NoError(t, err)

		var extractedServer *http.Server
		require.NoError(t, container.Extract(&extractedServer))
		require.NotNil(t, server)

		eqPtr(t, extractedServer, server)
		eqPtr(t, mux, server.Handler)

		var r bool
		require.NoError(t, container.Extract(&r))
	})

	t.Run("constructors with dependency and build error", func(t *testing.T) {
		var server = &http.Server{}
		var mux = &http.ServeMux{}

		container, err := inject.New(
			inject.Provide(func() (*http.ServeMux, error) {
				return mux, errors.New("build error")
			}),
			inject.Provide(func(mux *http.ServeMux) (*http.Server, error) {
				server.Handler = mux
				return server, nil
			}),
			inject.Provide(func(s *http.Server) bool {
				require.Equal(t, server, s)
				return true
			}),
		)

		require.NoError(t, err)

		var extractedServer *http.Server
		require.EqualError(t, container.Extract(&extractedServer), "*http.ServeMux: build error")
	})

	t.Run("named interface", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, inject.WithName("first")),
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}, inject.WithName("second")),
		)

		require.NoError(t, err)
		var addr *net.TCPAddr
		require.NoError(t, container.Extract(&addr, inject.Name("second")))
		require.Equal(t, "second", addr.Zone)
		require.NoError(t, container.Extract(&addr, inject.Name("first")))
		require.Equal(t, "first", addr.Zone)
	})

	t.Run("provide nil", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(nil),
		)

		require.EqualError(t, err, "could not compile container: could not provide nil")
	})

	t.Run("constructor provide nil", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return nil
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.EqualError(t, container.Extract(&addr), "*net.TCPAddr: nil provided")
	})

	t.Run("constructor without result", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() {}),
		)

		require.EqualError(t, err, "could not compile container: provide failed: github.com/defval/inject_test.TestContainer_ProvideConstructor.func12.1 must have at least one return value")
	})

	t.Run("constructor more than two return values", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() (*net.TCPAddr, *net.UDPAddr, *http.Server) {
				return nil, nil, nil
			}),
		)

		require.EqualError(t, err, "could not compile container: provide failed: github.com/defval/inject_test.TestContainer_ProvideConstructor.func13.1: constructor may have maximum two return values")
	})

	t.Run("constructor with two types", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() (*net.TCPAddr, *net.UDPAddr) {
				return &net.TCPAddr{}, &net.UDPAddr{}
			}),
		)

		require.EqualError(t, err, "could not compile container: provide failed: github.com/defval/inject_test.TestContainer_ProvideConstructor.func14.1: second argument of constructor must be error, got *net.UDPAddr")
	})

	t.Run("duplicate", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
		)

		require.EqualError(t, err, "could not compile container: *net.TCPAddr: use named definition if you have several instances of the same type")
	})

	t.Run("unknown argument", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func(addr *net.TCPAddr) net.Addr {
				return net.Addr(addr)
			}),
		)

		require.EqualError(t, err, "could not compile container: type *net.TCPAddr not provided")
	})
}

func TestContainer_ProvideStructPointer(t *testing.T) {
	t.Run("struct pointer with tags", func(t *testing.T) {
		var defaultMux = &http.ServeMux{}
		var anotherMux = &http.ServeMux{}

		var defaultServer = &http.Server{}
		var anotherServer = &http.Server{}

		type Server struct {
			private  string
			private2 string

			DefaultServer *http.Server `inject:""`        // default server
			AnotherServer *http.Server `inject:"another"` // another server

			Public  string
			Public2 string
		}

		type Muxes struct {
			DefaultMux *http.ServeMux `inject:""`
			private    string
			private2   string
			Public     string
			Public2    string
			AnotherMux *http.ServeMux `inject:"another"`
		}

		container, err := inject.New(
			inject.Provide(func() *http.ServeMux {
				return defaultMux
			}),
			inject.Provide(func() *http.ServeMux {
				return anotherMux
			}, inject.WithName("another")),
			inject.Provide(&Muxes{}),
			inject.Provide(func(muxes *Muxes) *http.Server {
				defaultServer.Handler = muxes.DefaultMux
				return defaultServer
			}),
			inject.Provide(func(muxes *Muxes) *http.Server {
				anotherServer.Handler = muxes.AnotherMux
				return anotherServer
			}, inject.WithName("another")),
			inject.Provide(&Server{}),
		)

		require.NoError(t, err)

		var server *Server
		require.NoError(t, container.Extract(&server))

		eqPtr(t, defaultServer, server.DefaultServer)
		eqPtr(t, defaultServer.Handler, defaultMux)

		eqPtr(t, anotherServer, server.AnotherServer)
		eqPtr(t, anotherServer.Handler, anotherMux)
	})

	t.Run("struct pointer with exported option", func(t *testing.T) {
		var defaultMux = &http.ServeMux{}
		var anotherMux = &http.ServeMux{}

		var defaultServer = &http.Server{}
		var anotherServer = &http.Server{}

		type Server struct {
			private  string
			private2 string

			DefaultServer *http.Server
			AnotherServer *http.Server `inject:"another"` // another server
		}

		type Muxes struct {
			DefaultMux *http.ServeMux
			private    string
			private2   string
			AnotherMux *http.ServeMux `inject:"another"`
		}

		container, err := inject.New(
			inject.Provide(func() *http.ServeMux {
				return defaultMux
			}),
			inject.Provide(func() *http.ServeMux {
				return anotherMux
			}, inject.WithName("another")),
			inject.Provide(&Muxes{}, inject.Exported()),
			inject.Provide(func(muxes *Muxes) *http.Server {
				defaultServer.Handler = muxes.DefaultMux
				return defaultServer
			}),
			inject.Provide(func(muxes *Muxes) *http.Server {
				anotherServer.Handler = muxes.AnotherMux
				return anotherServer
			}, inject.WithName("another")),
			inject.Provide(&Server{}, inject.Exported()),
		)

		require.NoError(t, err)

		var server *Server
		require.NoError(t, container.Extract(&server))

		eqPtr(t, defaultServer, server.DefaultServer)
		eqPtr(t, defaultServer.Handler, defaultMux)

		eqPtr(t, anotherServer, server.AnotherServer)
		eqPtr(t, anotherServer.Handler, anotherMux)
	})

	t.Run("struct with unknown field", func(t *testing.T) {
		type StructProvider struct {
			TCPAddr *net.TCPAddr
			UDPAddr *net.UDPAddr
			String  string
		}

		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "tcp"}
			}),
			inject.Provide(func() *net.UDPAddr {
				return &net.UDPAddr{Zone: "udp"}
			}),
			inject.Provide(&StructProvider{}, inject.Exported()),
		)

		require.Nil(t, container)
		require.EqualError(t, err, "could not compile container: type string not provided") // todo: improve message
	})
}

func TestContainer_ProvideStructValue(t *testing.T) {
	t.Run("struct with exported option", func(t *testing.T) {
		var defaultServer = &http.Server{}
		var anotherServer = &http.Server{}

		type Server struct {
			DefaultServer *http.Server
			AnotherServer *http.Server `inject:"another"` // another server
		}

		var servers = Server{}

		container, err := inject.New(
			inject.Provide(defaultServer),
			inject.Provide(anotherServer, inject.WithName("another")),
			inject.Provide(servers, inject.Exported()),
		)

		require.NoError(t, err)

		var server Server
		require.NoError(t, container.Extract(&server))

		eqPtr(t, defaultServer, server.DefaultServer)
		eqPtr(t, anotherServer, server.AnotherServer)
	})
}

func TestContainer_ProvideAs(t *testing.T) {
	t.Run("provide as interface", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "test",
				}
			}, inject.As(new(net.Addr))),
		)

		require.NoError(t, err)

		var addr net.Addr
		require.NoError(t, container.Extract(&addr))
		require.Equal(t, "test", addr.(*net.TCPAddr).Zone)
	})

	t.Run("provide as named interface", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return defaultAddr
			}, inject.As(new(net.Addr))),
			inject.Provide(func() *net.TCPAddr {
				return anotherAddr
			}, inject.As(new(net.Addr)), inject.WithName("another")),
		)

		require.NoError(t, err)

		var addr net.Addr
		require.NoError(t, container.Extract(&addr))

		eqPtr(t, defaultAddr, addr)

		require.NoError(t, container.Extract(&addr, inject.Name("another")))
		eqPtr(t, anotherAddr, addr)
	})

	t.Run("provide as struct", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, inject.As(http.Server{})),
		)

		require.EqualError(t, err, "could not compile container: provide failed: argument for As() must be pointer to interface type, got http.Server")
	})

	t.Run("provide as struct pointer", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, inject.As(new(http.Server))),
		)

		require.EqualError(t, err, "could not compile container: provide failed: argument for As() must be pointer to interface type, got *http.Server")
	})

	t.Run("provide as not implemented interface", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, inject.As(new(http.Handler))),
		)

		require.EqualError(t, err, "could not compile container: provide failed: *net.TCPAddr not implement http.Handler interface")
	})

	t.Run("provide as interface with struct injection", func(t *testing.T) {
		type TestStruct struct {
			Addr net.Addr `inject:""`
		}

		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "zone",
				}
			}, inject.As(new(net.Addr))),
			inject.Provide(&TestStruct{}),
		)

		require.NoError(t, err)

		var s *TestStruct
		require.NoError(t, container.Extract(&s))
		require.NotNil(t, s.Addr)
		require.Equal(t, "zone", s.Addr.(*net.TCPAddr).Zone)
	})
}

func TestContainer_Bundle(t *testing.T) {
	t.Run("bundle", func(t *testing.T) {
		container, err := inject.New(
			inject.Bundle(
				inject.Provide(func() *net.TCPAddr {
					return &net.TCPAddr{
						Zone: "zone",
						Port: 5432,
					}
				}),
			),
			inject.Bundle(
				inject.Provide(func(addr *net.TCPAddr) string {
					return addr.String()
				}),
			),
		)

		require.NoError(t, err)
		var s string
		require.NoError(t, container.Extract(&s))
		require.Equal(t, s, "%zone:5432")
	})
}

func TestContainer_Extract(t *testing.T) {
	t.Run("not pointer", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide(func() string {
				return "string"
			}),
			inject.Provide(func() int32 {
				return 32
			}),
		)

		require.NoError(t, err)

		var s string
		require.NoError(t, container.Extract(&s))
		require.Equal(t, s, "string")

		var i32 int32
		require.NoError(t, container.Extract(&i32))
		require.Equal(t, i32, int32(32))
	})

	t.Run("not existing type", func(t *testing.T) {
		container, err := inject.New()

		require.NoError(t, err)

		var s string
		require.EqualError(t, container.Extract(&s), "type string not provided")
	})

	t.Run("nil", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
		)

		require.NoError(t, err)

		require.EqualError(t, container.Extract(nil), "extract target must be a pointer")
	})

	t.Run("not provided named type", func(t *testing.T) {
		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, inject.WithName("first")),
			inject.Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.EqualError(t, container.Extract(&addr, inject.Name("second")), "type *net.TCPAddr not provided")
	})
}

func TestContainer_Group(t *testing.T) {
	t.Run("group with one implementation", func(t *testing.T) {
		handler := &http.ServeMux{}

		container, err := inject.New(
			inject.Provide(func() *http.ServeMux {
				return handler
			}, inject.As(new(http.Handler))),
		)

		require.NoError(t, err)

		var handlers []http.Handler
		require.NoError(t, container.Extract(&handlers))
		eqPtr(t, handler, handlers[0])
	})

	t.Run("group with different implementation types", func(t *testing.T) {
		var tcpAddr = &net.TCPAddr{}
		var udpAddr = &net.UDPAddr{}

		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return tcpAddr
			}, inject.As(new(net.Addr))),
			inject.Provide(func() *net.UDPAddr {
				return udpAddr
			}, inject.As(new(net.Addr))),
			inject.Provide(func(addrs []net.Addr) bool {
				eqPtr(t, tcpAddr, addrs[0])
				eqPtr(t, udpAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.NoError(t, container.Extract(&addrs))
		require.Len(t, addrs, 2)

		eqPtr(t, tcpAddr, addrs[0])
		eqPtr(t, udpAddr, addrs[1])

		var result bool
		require.NoError(t, container.Extract(&result))
		require.True(t, result)

	})

	t.Run("group with one implementation type", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := inject.New(
			inject.Provide(func() *net.TCPAddr {
				return defaultAddr
			}, inject.As(new(net.Addr))),
			inject.Provide(func() *net.TCPAddr {
				return anotherAddr
			}, inject.As(new(net.Addr)), inject.WithName("another")),
			inject.Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.NoError(t, container.Extract(&addrs))
		require.Len(t, addrs, 2)

		eqPtr(t, defaultAddr, addrs[0])
		eqPtr(t, anotherAddr, addrs[1])

		var result bool
		require.NoError(t, container.Extract(&result))
		require.True(t, result)

		var addr net.Addr

		require.NoError(t, container.Extract(&addr))
		eqPtr(t, defaultAddr, addr)

		require.NoError(t, container.Extract(&addr, inject.Name("another")))
		eqPtr(t, anotherAddr, addr)
	})

	t.Run("complex group", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := inject.New(
			inject.Provide("127.0.0.1"),
			inject.Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return defaultAddr
			}, inject.As(new(net.Addr))),
			inject.Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return anotherAddr
			}, inject.As(new(net.Addr)), inject.WithName("another")),
			inject.Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.NoError(t, container.Extract(&addrs))
		require.Len(t, addrs, 2)

		eqPtr(t, defaultAddr, addrs[0])
		eqPtr(t, anotherAddr, addrs[1])

		var result bool
		require.NoError(t, container.Extract(&result))
		require.True(t, result)

		var addr net.Addr

		require.NoError(t, container.Extract(&addr))
		eqPtr(t, defaultAddr, addr)

		require.NoError(t, container.Extract(&addr, inject.Name("another")))
		eqPtr(t, anotherAddr, addr)
	})

	t.Run("complex group with dependency error", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := inject.New(
			inject.Provide(func() (string, error) {
				return "", errors.Errorf("build error")
			}),
			inject.Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return defaultAddr
			}, inject.As(new(net.Addr))),
			inject.Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return anotherAddr
			}, inject.As(new(net.Addr)), inject.WithName("another")),
			inject.Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.EqualError(t, container.Extract(&addrs), "string: build error")
	})

	t.Run("complex group with dependency error", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := inject.New(
			inject.Provide(func() (string, error) {
				return "127.0.0.1", nil
			}),
			inject.Provide(func(addr string) (*net.TCPAddr, error) {
				require.Equal(t, "127.0.0.1", addr)
				return defaultAddr, errors.New("build error")
			}, inject.As(new(net.Addr))),
			inject.Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return anotherAddr
			}, inject.As(new(net.Addr)), inject.WithName("another")),
			inject.Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.EqualError(t, container.Extract(&addrs), "*net.TCPAddr: build error")
	})
}

// Stringer
type Stringer struct {
	s string
}

func (s *Stringer) String() string {
	return s.s
}

// MockStringer
type MockStringer struct {
	s string
}

func (s *MockStringer) String() string {
	return s.s
}

func TestContainer_Replace(t *testing.T) {
	t.Run("replace by mock", func(t *testing.T) {
		var stringer = &Stringer{
			s: "default",
		}
		var mockStringer = &MockStringer{
			s: "mock",
		}

		container, err := inject.New(
			inject.Provide(func() *Stringer {
				return stringer
			}, inject.As(new(fmt.Stringer))),
			inject.Replace(func() *MockStringer {
				return mockStringer
			}, inject.As(new(fmt.Stringer))),
			inject.Provide(func(s fmt.Stringer) bool {
				eqPtr(t, mockStringer, s)
				return true
			}),
		)

		require.NoError(t, err)

		var s fmt.Stringer
		require.NoError(t, container.Extract(&s))
		eqPtr(t, s, mockStringer)

		var b bool
		require.NoError(t, container.Extract(&b))
	})

	t.Run("replace named interface by mock", func(t *testing.T) {
		var stringer = &Stringer{s: "default"}
		var anotherStringer = &Stringer{s: "another"}
		var mockStringer = &MockStringer{s: "mock"}

		container, err := inject.New(
			inject.Provide(func() *Stringer {
				return stringer
			}, inject.As(new(fmt.Stringer))),
			inject.Provide(func() *Stringer {
				return anotherStringer
			}, inject.As(new(fmt.Stringer)), inject.WithName("another")),
			inject.Replace(func() *MockStringer {
				return mockStringer
			}, inject.As(new(fmt.Stringer)), inject.WithName("another")),
		)

		require.NoError(t, err)

		var s fmt.Stringer
		require.NoError(t, container.Extract(&s))
		eqPtr(t, s, stringer)

		require.NoError(t, container.Extract(&s, inject.Name("another")))
		eqPtr(t, s, mockStringer)
	})

	t.Run("replace nil provider", func(t *testing.T) {
		_, err := inject.New(
			inject.Replace(nil),
		)

		require.EqualError(t, err, "could not compile container: replace provider could not be nil")
	})

	t.Run("replace without interfaces", func(t *testing.T) {
		_, err := inject.New(
			inject.Replace(func() fmt.Stringer {
				return &Stringer{}
			}),
		)

		require.EqualError(t, err, "could not compile container: fmt.Stringer: no one interface has been replaced, use `inject.As()` for specify it")
	})

	t.Run("replace incorrect constructor signature", func(t *testing.T) {
		_, err := inject.New(
			inject.Replace(func() {}),
		)

		require.EqualError(t, err, "could not compile container: provide failed: github.com/defval/inject_test.TestContainer_Replace.func5.1 must have at least one return value")
	})

	t.Run("replace already provided type", func(t *testing.T) {
		var stringer = &Stringer{s: "default"}
		var anotherStringer = &Stringer{s: "another"}

		container, err := inject.New(
			inject.Provide(func() *Stringer {
				return stringer
			}, inject.As(new(fmt.Stringer))),
			inject.Replace(func() *Stringer {
				return anotherStringer
			}, inject.As(new(fmt.Stringer))),
		)

		require.NoError(t, err)
		var s *Stringer
		require.NoError(t, container.Extract(&s))
		eqPtr(t, anotherStringer, s)

		var si fmt.Stringer
		require.NoError(t, container.Extract(&si))
		eqPtr(t, anotherStringer, si)
	})

	t.Run("replace unknown type", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide(func() *Stringer {
				return &Stringer{}
			}),
			inject.Replace(func() *MockStringer {
				return &MockStringer{}
			}, inject.As(new(fmt.Stringer))),
		)

		require.EqualError(t, err, "could not compile container: type fmt.Stringer not provided")
	})
}

func TestContainer_Cycle(t *testing.T) {
	t.Run("simple cycle", func(t *testing.T) {
		_, err := inject.New(
			inject.Provide("string"),
			inject.Provide(func(string, int32) bool {
				return true
			}),
			inject.Provide(func(bool) int64 {
				return 64
			}),
			inject.Provide(func(int64) int32 {
				return 32
			}),
		)

		require.EqualError(t, err, "could not compile container: detect cycle: string: bool: int64: int32: bool")
	})
}
