package inject

import (
	"net"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Provide(t *testing.T) {
	t.Run("function", func(t *testing.T) {
		var builder = &Container{}

		require.NoError(t, builder.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{
				Zone: "test",
			}
		}))

		var addr *net.TCPAddr
		err := builder.Populate(&addr)
		require.NoError(t, err)
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("function with nil error", func(t *testing.T) {
		var builder = &Container{}

		require.NoError(t, builder.Provide(func() (*net.TCPAddr, error) {
			return &net.TCPAddr{
				Zone: "test",
			}, nil
		}))

		var addr *net.TCPAddr
		err := builder.Populate(&addr)
		require.NoError(t, err)
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("function with error", func(t *testing.T) {
		var builder = &Container{}

		require.NoError(t, builder.Provide(func() (*net.TCPAddr, error) {
			return &net.TCPAddr{
				Zone: "test",
			}, errors.New("build error")
		}))

		var addr *net.TCPAddr
		err := builder.Populate(&addr)
		require.EqualError(t, err, "*net.TCPAddr: build error")
	})

	t.Run("function without arguments", func(t *testing.T) {
		var builder = &Container{}

		// todo: improve error message
		require.EqualError(t, builder.Provide(func() {}), "provide failed: provider must be a function with returned value and optional error")
	})

	// todo: implement struct provide
	t.Run("struct", func(t *testing.T) {
		var builder = &Container{}

		type StructProvider struct {
		}

		require.EqualError(t, builder.Provide(&StructProvider{}), "provide failed: struct provider not implemented yet")
	})
}

func TestBuilder_ProvideAs(t *testing.T) {
	t.Run("provide as", func(t *testing.T) {
		var builder = &Container{}

		require.NoError(t,
			builder.Provide(
				func() *net.TCPAddr {
					return &net.TCPAddr{
						Zone: "test",
					}
				},
				As(new(net.Addr)),
			),
		)

		var addr net.Addr
		require.NoError(t, builder.Populate(&addr))
		require.Equal(t, "test", addr.(*net.TCPAddr).Zone)
	})

	t.Run("provide as struct", func(t *testing.T) {
		var builder = &Container{}

		require.EqualError(t, builder.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{}
		}, As(http.Server{})), "provide failed: argument for As() must be pointer to interface type, got http.Server")
	})

	t.Run("provide as struct pointer", func(t *testing.T) {
		var builder = &Container{}

		require.EqualError(t, builder.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{}
		}, As(new(http.Server))), "provide failed: argument for As() must be pointer to interface type, got *http.Server")
	})

	t.Run("provide as not implemented interface", func(t *testing.T) {
		var builder = &Container{}

		require.EqualError(t, builder.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{}
		}, As(new(http.Handler))), "provide failed: *net.TCPAddr not implement http.Handler interface")
	})
}

//
// func TestBuilder_ProvideName(t *testing.T) {
// 	t.Run("provide two named implementations as one interface", func(t *testing.T) {
// 		var builder = &Container{}
//
// 		require.NoError(t, builder.Provide(func() *net.TCPAddr {
// 			return &net.TCPAddr{}
// 		}, As(new(net.Addr)), Name("tcp")))
//
// 		require.NoError(t, builder.Provide(func() *net.UDPAddr {
// 			return &net.UDPAddr{}
// 		}, As(new(net.Addr)), Name("udp")))
// 	})
//
// 	t.Run("provide two implementations as one interface without name", func(t *testing.T) {
// 		var builder = &Container{}
//
// 		require.NoError(t, builder.Provide(func() *net.TCPAddr {
// 			return &net.TCPAddr{}
// 		}, As(new(net.Addr))))
//
// 		require.NoError(t, builder.Provide(func() *net.UDPAddr {
// 			return &net.UDPAddr{}
// 		}, As(new(net.Addr))))
// 	})
// }
