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

func notEqPtr(t *testing.T, expected interface{}, actual interface{}) {
	require.NotEqual(t, fmt.Sprintf("%p", expected), fmt.Sprintf("%p", actual))
}

func TestContainer_ProvideConstructor(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		container, err := New(
			Provide("string"),
			Provide(func(s string) bool {
				require.Equal(t, "string", s)
				return true
			}),
		)

		require.NoError(t, err)

		var result string
		require.NoError(t, container.Populate(&result))
		require.Equal(t, "string", result)

		var b bool
		require.NoError(t, container.Populate(&b))
	})

	t.Run("struct", func(t *testing.T) {
		_, err := New(
			Provide(struct{}{}),
		)

		require.NoError(t, err)
	})

	t.Run("slice", func(t *testing.T) {
		container, err := New(
			Provide([]int64{32, 30, 31}),
		)

		require.NoError(t, err)

		var result []int64
		require.NoError(t, container.Populate(&result))
		require.Len(t, result, 3)
	})

	t.Run("chan", func(t *testing.T) {
		var ch = make(chan struct{})

		container, err := New(
			Provide(ch),
			Provide(func(ch chan struct{}) bool {
				close(ch)
				return true
			}),
		)

		require.NoError(t, err)

		var b bool
		require.NoError(t, container.Populate(&b))
		_, more := <-ch
		require.False(t, more)
	})

	t.Run("map", func(t *testing.T) {
		var m = map[string]string{"test": "test"}

		container, err := New(
			Provide(m),
			Provide(func(arg map[string]string) bool {
				require.Equal(t, m, arg)
				return true
			}),
		)

		require.NoError(t, err)

		var b bool
		require.NoError(t, container.Populate(&b))
		require.True(t, b)
	})

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

	t.Run("constructors with dependency and build error", func(t *testing.T) {
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

	t.Run("named interface", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "first"}
			}, WithName("first")),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}, WithName("second")),
		)

		require.NoError(t, err)
		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr, Name("second")))
		require.Equal(t, "second", addr.Zone)
		require.NoError(t, container.Populate(&addr, Name("first")))
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

	t.Run("constructor with two types", func(t *testing.T) {
		_, err := New(
			Provide(func() (*net.TCPAddr, *net.UDPAddr) {
				return &net.TCPAddr{}, &net.UDPAddr{}
			}),
		)

		require.EqualError(t, err, "could not compile container: provide failed: constructor must be a function with value and optional error as result")
	})

	t.Run("duplicate", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}),
		)

		require.EqualError(t, err, "could not compile container: *net.TCPAddr: use named definition if you have several instances of the same type")
	})

	t.Run("unknown argument", func(t *testing.T) {
		_, err := New(
			Provide(func(addr *net.TCPAddr) net.Addr {
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

		container, err := New(
			Provide(func() *http.ServeMux {
				return defaultMux
			}),
			Provide(func() *http.ServeMux {
				return anotherMux
			}, WithName("another")),
			Provide(&Muxes{}),
			Provide(func(muxes *Muxes) *http.Server {
				defaultServer.Handler = muxes.DefaultMux
				return defaultServer
			}),
			Provide(func(muxes *Muxes) *http.Server {
				anotherServer.Handler = muxes.AnotherMux
				return anotherServer
			}, WithName("another")),
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

		container, err := New(
			Provide(func() *http.ServeMux {
				return defaultMux
			}),
			Provide(func() *http.ServeMux {
				return anotherMux
			}, WithName("another")),
			Provide(&Muxes{}, Exported()),
			Provide(func(muxes *Muxes) *http.Server {
				defaultServer.Handler = muxes.DefaultMux
				return defaultServer
			}),
			Provide(func(muxes *Muxes) *http.Server {
				anotherServer.Handler = muxes.AnotherMux
				return anotherServer
			}, WithName("another")),
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

	t.Run("struct with unknown field", func(t *testing.T) {
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

func TestContainer_ProvideStructValue(t *testing.T) {
	t.Run("struct with exported option", func(t *testing.T) {
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

		var s = Server{}
		var m = Muxes{}

		container, err := New(
			Provide(defaultMux),
			Provide(anotherMux, WithName("another")),
			Provide(s, Exported()),
			Provide(func(muxes Muxes) *http.Server {
				defaultServer.Handler = muxes.DefaultMux
				return defaultServer
			}),
			Provide(func(muxes Muxes) *http.Server {
				anotherServer.Handler = muxes.AnotherMux
				return anotherServer
			}, WithName("another")),
			Provide(m, Exported()),
		)

		require.NoError(t, err)

		var server Server
		require.NoError(t, container.Populate(&server))

		notEqPtr(t, &defaultServer, &server.DefaultServer)
		notEqPtr(t, &defaultServer.Handler, &defaultMux)

		notEqPtr(t, &anotherServer, &server.AnotherServer)
		notEqPtr(t, &anotherServer.Handler, &anotherMux)
	})
}

func TestContainer_ProvideAs(t *testing.T) {
	t.Run("provide as interface", func(t *testing.T) {
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

	t.Run("provide as named interface", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := New(
			Provide(func() *net.TCPAddr {
				return defaultAddr
			}, As(new(net.Addr))),
			Provide(func() *net.TCPAddr {
				return anotherAddr
			}, As(new(net.Addr)), WithName("another")),
		)

		require.NoError(t, err)

		var addr net.Addr
		require.NoError(t, container.Populate(&addr))

		eqPtr(t, defaultAddr, addr)

		require.NoError(t, container.Populate(&addr, Name("another")))
		eqPtr(t, anotherAddr, addr)
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
			}, WithName("first")),
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "second"}
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.EqualError(t, container.Populate(&addr, Name("second")), "type *net.TCPAddr not provided")
	})
}

func TestContainer_Group(t *testing.T) {
	t.Run("group with different implementation types", func(t *testing.T) {
		var tcpAddr = &net.TCPAddr{}
		var udpAddr = &net.UDPAddr{}

		container, err := New(
			Provide(func() *net.TCPAddr {
				return tcpAddr
			}, As(new(net.Addr))),
			Provide(func() *net.UDPAddr {
				return udpAddr
			}, As(new(net.Addr))),
			Provide(func(addrs []net.Addr) bool {
				eqPtr(t, tcpAddr, addrs[0])
				eqPtr(t, udpAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.NoError(t, container.Populate(&addrs))
		require.Len(t, addrs, 2)

		eqPtr(t, tcpAddr, addrs[0])
		eqPtr(t, udpAddr, addrs[1])

		var result bool
		require.NoError(t, container.Populate(&result))
		require.True(t, result)

	})

	t.Run("group with one implementation type", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := New(
			Provide(func() *net.TCPAddr {
				return defaultAddr
			}, As(new(net.Addr))),
			Provide(func() *net.TCPAddr {
				return anotherAddr
			}, As(new(net.Addr)), WithName("another")),
			Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.NoError(t, container.Populate(&addrs))
		require.Len(t, addrs, 2)

		eqPtr(t, defaultAddr, addrs[0])
		eqPtr(t, anotherAddr, addrs[1])

		var result bool
		require.NoError(t, container.Populate(&result))
		require.True(t, result)

		var addr net.Addr

		require.NoError(t, container.Populate(&addr))
		eqPtr(t, defaultAddr, addr)

		require.NoError(t, container.Populate(&addr, Name("another")))
		eqPtr(t, anotherAddr, addr)
	})

	t.Run("complex group", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := New(
			Provide("127.0.0.1"),
			Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return defaultAddr
			}, As(new(net.Addr))),
			Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return anotherAddr
			}, As(new(net.Addr)), WithName("another")),
			Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.NoError(t, container.Populate(&addrs))
		require.Len(t, addrs, 2)

		eqPtr(t, defaultAddr, addrs[0])
		eqPtr(t, anotherAddr, addrs[1])

		var result bool
		require.NoError(t, container.Populate(&result))
		require.True(t, result)

		var addr net.Addr

		require.NoError(t, container.Populate(&addr))
		eqPtr(t, defaultAddr, addr)

		require.NoError(t, container.Populate(&addr, Name("another")))
		eqPtr(t, anotherAddr, addr)
	})

	t.Run("complex group with dependency error", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := New(
			Provide(func() (string, error) {
				return "", errors.Errorf("build error")
			}),
			Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return defaultAddr
			}, As(new(net.Addr))),
			Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return anotherAddr
			}, As(new(net.Addr)), WithName("another")),
			Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.EqualError(t, container.Populate(&addrs), "string: build error")
	})

	t.Run("complex group with dependency error", func(t *testing.T) {
		var defaultAddr = &net.TCPAddr{}
		var anotherAddr = &net.TCPAddr{}

		container, err := New(
			Provide(func() (string, error) {
				return "127.0.0.1", nil
			}),
			Provide(func(addr string) (*net.TCPAddr, error) {
				require.Equal(t, "127.0.0.1", addr)
				return defaultAddr, errors.New("build error")
			}, As(new(net.Addr))),
			Provide(func(addr string) *net.TCPAddr {
				require.Equal(t, "127.0.0.1", addr)
				return anotherAddr
			}, As(new(net.Addr)), WithName("another")),
			Provide(func(addrs []net.Addr) bool {
				eqPtr(t, defaultAddr, addrs[0])
				eqPtr(t, anotherAddr, addrs[1])
				return len(addrs) == 2
			}),
		)

		require.NoError(t, err)

		var addrs []net.Addr
		require.EqualError(t, container.Populate(&addrs), "*net.TCPAddr: build error")
	})
}

func TestContainer_Cycle(t *testing.T) {
	t.Run("simple cycle", func(t *testing.T) {
		_, err := New(
			Provide("string"),
			Provide(func(string, int32) bool {
				return true
			}),
			Provide(func(bool) int64 {
				return 64
			}),
			Provide(func(int64) int32 {
				return 32
			}),
		)

		require.EqualError(t, err, "could not compile container: detect cycle: string: bool: int64: int32: bool")
	})
}
