package inject

import (
	"net"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestContainer_Provide(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		container, err := New(
			Provide(func() *http.Server {
				return &http.Server{}
			}),
			Provide(func(server *http.Server) *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "test",
				}
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("constructor with nil error", func(t *testing.T) {
		container, err := New(
			Provide(func() (*net.TCPAddr, error) {
				return &net.TCPAddr{
					Zone: "test",
				}, nil
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("constructor with error", func(t *testing.T) {
		container, err := New(
			Provide(func() (*net.TCPAddr, error) {
				return &net.TCPAddr{
					Zone: "test",
				}, errors.New("build error")
			}),
			Provide(func(addr *net.TCPAddr) net.Addr {
				return addr
			}),
		)

		require.NoError(t, err)

		var addr net.Addr
		require.EqualError(t, container.Populate(&addr), "net.Addr: *net.TCPAddr: build error")
	})

	t.Run("struct", func(t *testing.T) {
		type StructProvider struct {
			TCPAddr *net.TCPAddr `inject:""`
			Public  string
			private string
			UDPAddr *net.UDPAddr `inject:""`
		}

		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "tcp"}
			}),
			Provide(func() *net.UDPAddr {
				return &net.UDPAddr{Zone: "udp"}
			}),
			Provide(&StructProvider{}),
		)

		require.NoError(t, err)

		var sp *StructProvider
		require.NoError(t, container.Populate(&sp))
		require.Equal(t, "tcp", sp.TCPAddr.Zone)
		require.Equal(t, "udp", sp.UDPAddr.Zone)
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
		require.EqualError(t, err, "could not compile container: provide failed: provider must be a function with value and optional error as result")
	})

	t.Run("constructor with incorrect signature", func(t *testing.T) {
		_, err := New(
			Provide(func() (*net.TCPAddr, *net.UDPAddr) {
				return &net.TCPAddr{}, &net.UDPAddr{}
			}),
		)

		require.EqualError(t, err, "could not compile container: provide failed: provider must be a function with value and optional error as result")
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

		require.EqualError(t, err, "could not compile container: could not add definition: *net.TCPAddr already provided")
	})

	t.Run("unknown injection type", func(t *testing.T) {
		_, err := New(
			Provide(func(addr *net.TCPAddr) net.Addr {
				return net.Addr(addr)
			}),
		)

		require.EqualError(t, err, "could not compile container: type *net.TCPAddr not provided")
	})

	t.Run("incorrect provider type", func(t *testing.T) {
		_, err := New(
			Provide("string"),
		)

		require.EqualError(t, err, "could not compile container: provide failed: provider must be a function with value and optional error as result")
	})

	t.Run("cycle", func(t *testing.T) {
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

func TestApply(t *testing.T) {
	t.Run("apply function", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) {
				addr.Zone = "two"
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "two", addr.Zone)
	})

	t.Run("apply without result", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) {
				addr.Zone = "two"
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "two", addr.Zone)
	})

	t.Run("apply error", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) (err error) {
				return errors.New("boom")
			}),
		)

		require.EqualError(t, err, "could not compile container: apply error: boom")
	})

	t.Run("apply incorrect function", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) (s string) {
				return "string"
			}),
		)

		require.EqualError(t, err, "could not compile container: modifier must be a function with optional error as result")
	})

	t.Run("nil", func(t *testing.T) {
		_, err := New(
			Apply(nil),
		)

		require.EqualError(t, err, "could not compile container: nil modifier")
	})

	t.Run("apply ptr", func(t *testing.T) {
		_, err := New(
			Apply(&net.TCPAddr{}),
		)

		require.EqualError(t, err, "could not compile container: modifier must be a function with optional error as result")
	})

	t.Run("more than one result", func(t *testing.T) {
		_, err := New(
			Apply(func() (string, error) {
				return "string", nil
			}),
		)

		require.EqualError(t, err, "could not compile container: modifier must be a function with optional error as result")
	})

	t.Run("use unknown type", func(t *testing.T) {
		_, err := New(
			Apply(func(*net.TCPAddr) {}),
		)

		require.EqualError(t, err, "could not compile container: type *net.TCPAddr not provided")
	})

	t.Run("apply argument instantiate error", func(t *testing.T) {
		_, err := New(
			Provide(func() (*net.TCPAddr, error) {
				return nil, errors.New("wow")
			}),
			Apply(func(*net.TCPAddr) {}),
		)

		require.EqualError(t, err, "could not compile container: *net.TCPAddr: wow")
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
	t.Run("not existing type", func(t *testing.T) {
		container, err := New()

		require.NoError(t, err)

		var s string
		require.EqualError(t, container.Populate(&s), "type string not provided")
	})
}

//
// func TestContainer_ProvideName(t *testing.T) {
// 	t.Run("provide two named implementations as one interface", func(t *testing.T) {
// 		var container = &Container{}
//
// 		require.NoError(t, container.Provide(func() *net.TCPAddr {
// 			return &net.TCPAddr{}
// 		}, As(new(net.Addr)), Name("tcp")))
//
// 		require.NoError(t, container.Provide(func() *net.UDPAddr {
// 			return &net.UDPAddr{}
// 		}, As(new(net.Addr)), Name("udp")))
// 	})
//
// 	t.Run("provide two implementations as one interface without name", func(t *testing.T) {
// 		var container = &Container{}
//
// 		require.NoError(t, container.Provide(func() *net.TCPAddr {
// 			return &net.TCPAddr{}
// 		}, As(new(net.Addr))))
//
// 		require.NoError(t, container.Provide(func() *net.UDPAddr {
// 			return &net.UDPAddr{}
// 		}, As(new(net.Addr))))
// 	})
// }
