package inject

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func eqPtr(t *testing.T, expected interface{}, actual interface{}) {
	require.Equal(t, fmt.Sprintf("%p", expected), fmt.Sprintf("%p", actual))
}

func TestContainer_ProvideConstructor(t *testing.T) {
	t.Run("constructors with dependency without errors", func(t *testing.T) {
		var server = &http.Server{}
		var mux = &http.ServeMux{}

		container, err := New(
			Provide(func() *http.ServeMux {
				return mux
			}),
			Provide(func(mux *http.ServeMux) *http.Server {
				server.Handler = mux
				return server
			}),
			Provide(func(s *http.Server) bool {
				eqPtr(t, server, s)
				return true
			}),
		)

		require.NoError(t, err)

		// check populate
		var populatedServer *http.Server
		require.NoError(t, container.Populate(&populatedServer))
		require.NotNil(t, server)

		eqPtr(t, populatedServer, server)
		eqPtr(t, mux, server.Handler)

		var r bool
		require.NoError(t, container.Populate(&r))
	})

	t.Run("constructors with dependency with nil errors", func(t *testing.T) {
		var server = &http.Server{}
		var mux = &http.ServeMux{}

		container, err := New(
			Provide(func() (*http.ServeMux, error) {
				return mux, nil
			}),
			Provide(func(mux *http.ServeMux) (*http.Server, error) {
				server.Handler = mux
				return server, nil
			}),
			Provide(func(s *http.Server) bool {
				eqPtr(t, server, s)
				return true
			}),
		)

		require.NoError(t, err)

		// check populate
		var populatedServer *http.Server
		require.NoError(t, container.Populate(&populatedServer))
		require.NotNil(t, server)

		eqPtr(t, populatedServer, server)
		eqPtr(t, mux, server.Handler)

		var r bool
		require.NoError(t, container.Populate(&r))
	})

	t.Run("constructors with dependency with build error", func(t *testing.T) {
		var server = &http.Server{}
		var mux = &http.ServeMux{}

		container, err := New(
			Provide(func() (*http.ServeMux, error) {
				return mux, errors.New("build error")
			}),
			Provide(func(mux *http.ServeMux) (*http.Server, error) {
				server.Handler = mux
				return server, nil
			}),
			Provide(func(s *http.Server) bool {
				require.Equal(t, server, s)
				return true
			}),
		)

		require.NoError(t, err)

		// check populate
		var populatedServer *http.Server
		require.EqualError(t, container.Populate(&populatedServer), "*http.ServeMux: build error")
	})

	t.Run("two instance of one type with names", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, Name("first")),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}, Name("second")),
		)

		require.NoError(t, err)
		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr, PopulateName("second")))
		require.Equal(t, "second", addr.Zone)
		require.NoError(t, container.Populate(&addr, PopulateName("first")))
		require.Equal(t, "first", addr.Zone)
	})

	t.Run("provide nil", func(t *testing.T) {
		_, err := New(
			Provide(nil),
		)

		require.EqualError(t, err, "could not compile container: could not provide nil")
	})

	t.Run("constructor provide nil", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return nil
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.EqualError(t, container.Populate(&addr), "*net.TCPAddr: nil provided")
	})

	t.Run("constructor without result", func(t *testing.T) {
		_, err := New(
			Provide(func() {}),
		)

		// todo: improve error message
		require.EqualError(t, err, "could not compile container: provide failed: constructor must be a function with value and optional error as result")
	})

	t.Run("constructor with incorrect signature", func(t *testing.T) {
		_, err := New(
			Provide(func() (*net.TCPAddr, *net.UDPAddr) {
				return &net.TCPAddr{}, &net.UDPAddr{}
			}),
		)

		require.EqualError(t, err, "could not compile container: provide failed: constructor must be a function with value and optional error as result")
	})

	t.Run("provide duplicate type", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
		)

		require.EqualError(t, err, "could not compile container: could not add definition: *net.TCPAddr: use named definition if you have several instances of the same type")
	})

	t.Run("unknown injection type", func(t *testing.T) {
		_, err := New(
			Provide(func(addr *net.TCPAddr) net.Addr {
				return net.Addr(addr)
			}),
		)

		require.EqualError(t, err, "could not compile container: type *net.TCPAddr not provided")
	})

	t.Run("incorrect value type", func(t *testing.T) {
		_, err := New(
			Provide("string"),
		)

		require.EqualError(t, err, "could not compile container: provide failed: constructor must be a function with value and optional error as result")
	})
}

func TestContainer_ProvideStructPointer(t *testing.T) {
	t.Run("with tags", func(t *testing.T) {
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

		container, err := New(
			Provide(func() *http.ServeMux {
				return defaultMux
			}),
			Provide(func() *http.ServeMux {
				return anotherMux
			}, Name("another")),
			Provide(&Muxes{}),
			Provide(func(muxes *Muxes) *http.Server {
				defaultServer.Handler = muxes.DefaultMux
				return defaultServer
			}),
			Provide(func(muxes *Muxes) *http.Server {
				anotherServer.Handler = muxes.AnotherMux
				return anotherServer
			}, Name("another")),
			Provide(&Server{}),
		)

		require.NoError(t, err)

		var server *Server
		require.NoError(t, container.Populate(&server))

		eqPtr(t, defaultServer, server.DefaultServer)
		eqPtr(t, defaultServer.Handler, defaultMux)

		eqPtr(t, anotherServer, server.AnotherServer)
		eqPtr(t, anotherServer.Handler, anotherMux)
	})

	t.Run("with exported option", func(t *testing.T) {
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

		container, err := New(
			Provide(func() *http.ServeMux {
				return defaultMux
			}),
			Provide(func() *http.ServeMux {
				return anotherMux
			}, Name("another")),
			Provide(&Muxes{}, Exported()),
			Provide(func(muxes *Muxes) *http.Server {
				defaultServer.Handler = muxes.DefaultMux
				return defaultServer
			}),
			Provide(func(muxes *Muxes) *http.Server {
				anotherServer.Handler = muxes.AnotherMux
				return anotherServer
			}, Name("another")),
			Provide(&Server{}, Exported()),
		)

		require.NoError(t, err)

		var server *Server
		require.NoError(t, container.Populate(&server))

		eqPtr(t, defaultServer, server.DefaultServer)
		eqPtr(t, defaultServer.Handler, defaultMux)

		eqPtr(t, anotherServer, server.AnotherServer)
		eqPtr(t, anotherServer.Handler, anotherMux)
	})
}

func TestContainer_ProvideStructValue(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		var addr = net.TCPAddr{}

		container, err := New(
			Provide(addr),
		)

		require.NoError(t, err)

		var populatedAddr net.TCPAddr
		require.NoError(t, container.Populate(&populatedAddr))
		require.Equal(t, addr, populatedAddr)
		require.NotEqual(t, fmt.Sprintf("%p", &addr), fmt.Sprintf("%p", &populatedAddr))
	})

	t.Run("struct with not provided field", func(t *testing.T) {
		type StructProvider struct {
			TCPAddr *net.TCPAddr
			UDPAddr *net.UDPAddr
			String  string
		}

		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "tcp"}
			}),
			Provide(func() *net.UDPAddr {
				return &net.UDPAddr{Zone: "udp"}
			}),
			Provide(&StructProvider{}, Exported()),
		)

		require.Nil(t, container)
		require.EqualError(t, err, "could not compile container: type string not provided") // todo: improve message
	})
}

func TestContainer_ProvideAs(t *testing.T) {
	t.Run("provide as", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "test",
				}
			}, As(new(net.Addr))),
		)

		require.NoError(t, err)

		var addr net.Addr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "test", addr.(*net.TCPAddr).Zone)
	})

	t.Run("provide as struct", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, As(http.Server{})),
		)

		require.EqualError(t, err, "could not compile container: provide failed: argument for As() must be pointer to interface type, got http.Server")
	})

	t.Run("provide as struct pointer", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, As(new(http.Server))),
		)

		require.EqualError(t, err, "could not compile container: provide failed: argument for As() must be pointer to interface type, got *http.Server")
	})

	t.Run("provide as not implemented interface", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, As(new(http.Handler))),
		)

		require.EqualError(t, err, "could not compile container: provide failed: *net.TCPAddr not implement http.Handler interface")
	})

	t.Run("provide as interface with struct injection", func(t *testing.T) {
		type TestStruct struct {
			Addr net.Addr `inject:""`
		}

		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "zone",
				}
			}, As(new(net.Addr))),
			Provide(&TestStruct{}),
		)

		require.NoError(t, err)

		var s *TestStruct
		require.NoError(t, container.Populate(&s))
		require.NotNil(t, s.Addr)
		require.Equal(t, "zone", s.Addr.(*net.TCPAddr).Zone)
	})
}

func TestContainer_Package(t *testing.T) {
	t.Run("package", func(t *testing.T) {
		container, err := New(
			Package(
				Provide(func() *net.TCPAddr {
					return &net.TCPAddr{
						Zone: "zone",
						Port: 5432,
					}
				}),
			),
			Package(
				Provide(func(addr *net.TCPAddr) string {
					return addr.String()
				}),
			),
		)

		require.NoError(t, err)
		var s string
		require.NoError(t, container.Populate(&s))
		require.Equal(t, s, "%zone:5432")
	})
}

func TestContainer_Populate(t *testing.T) {
	t.Run("not pointer", func(t *testing.T) {
		container, err := New(
			Provide(func() string {
				return "string"
			}),
			Provide(func() int32 {
				return 32
			}),
		)

		require.NoError(t, err)

		var s string
		require.NoError(t, container.Populate(&s))
		require.Equal(t, s, "string")

		var i32 int32
		require.NoError(t, container.Populate(&i32))
		require.Equal(t, i32, int32(32))
	})

	t.Run("not existing type", func(t *testing.T) {
		container, err := New()

		require.NoError(t, err)

		var s string
		require.EqualError(t, container.Populate(&s), "type string not provided")
	})

	t.Run("nil", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
		)

		require.NoError(t, err)

		require.EqualError(t, container.Populate(nil), "populate target must be a not nil pointer")
	})

	t.Run("not provided named type", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, Name("first")),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.EqualError(t, container.Populate(&addr, PopulateName("second")), "type *net.TCPAddr not provided")
	})
}

func TestContainer_Group(t *testing.T) {
	t.Run("inject group", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "tcp",
				}
			}, As(new(net.Addr))),
			Provide(func() *net.UDPAddr {
				return &net.UDPAddr{
					Zone: "udp",
				}
			}, As(new(net.Addr))),
			Provide(func(addrs []net.Addr) bool {
				require.Equal(t, "tcp", addrs[0].(*net.TCPAddr).Zone)
				require.Equal(t, "udp", addrs[1].(*net.UDPAddr).Zone)
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)
		var result bool
		require.NoError(t, container.Populate(&result))
		require.True(t, result)

	})

	t.Run("different types", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "tcp",
				}
			}, As(new(net.Addr))),
			Provide(func() *net.UDPAddr {
				return &net.UDPAddr{
					Zone: "udp",
				}
			}, As(new(net.Addr))),
		)

		require.NoError(t, err)
		var addrs []net.Addr
		require.NoError(t, container.Populate(&addrs))
		require.Len(t, addrs, 2)
		require.Equal(t, "tcp", addrs[0].(*net.TCPAddr).Zone)
		require.Equal(t, "udp", addrs[1].(*net.UDPAddr).Zone)
	})

	t.Run("one type without name", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, As(new(net.Addr))),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}, As(new(net.Addr))),
		)

		require.Nil(t, container)
		require.EqualError(t, err, "could not compile container: could not add definition: *net.TCPAddr: use named definition if you have several instances of the same type")
	})

	t.Run("default value of group", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, Name("first"), As(new(net.Addr))),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}, As(new(net.Addr))),
		)

		require.NoError(t, err)
		var addr net.Addr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "second", addr.(*net.TCPAddr).Zone)
	})

	t.Run("named value of group", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, Name("first"), As(new(net.Addr))),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}, As(new(net.Addr))),
		)

		require.NoError(t, err)
		var addr net.Addr
		require.NoError(t, container.Populate(&addr, PopulateName("first")))
		require.Equal(t, "first", addr.(*net.TCPAddr).Zone)
	})
}

func TestContainer_Cycle(t *testing.T) {
	t.Run("simple cycle", func(t *testing.T) {
		_, err := New(
			Provide(func(string) bool {
				return true
			}),
			Provide(func(bool) int64 {
				return 64
			}),
			Provide(func(int64) int32 {
				return 32
			}),
			Provide(func(int32) string {
				return "string"
			}),
		)

		require.EqualError(t, err, "could not compile container: detect cycle: bool: int64: int32: string: bool")
	})
}
