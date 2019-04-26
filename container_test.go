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
		var container = &Container{}

		container.Provide(func() *http.Server {
			return &http.Server{}
		})

		container.Provide(func(server *http.Server) *net.TCPAddr {
			return &net.TCPAddr{
				Zone: "test",
			}
		})

		require.NoError(t, container.Compile())

		var addr *net.TCPAddr
		err := container.Populate(&addr)
		require.NoError(t, err)
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("function with error", func(t *testing.T) {
		var container = &Container{}

		container.Provide(func() (*net.TCPAddr, error) {
			return &net.TCPAddr{
				Zone: "test",
			}, errors.New("build error")
		})

		require.NoError(t, container.Compile())

		var addr *net.TCPAddr
		err := container.Populate(&addr)
		require.EqualError(t, err, "*net.TCPAddr: build error")
	})

	t.Run("function with nil error", func(t *testing.T) {
		var container = &Container{}

		container.Provide(func() (*net.TCPAddr, error) {
			return &net.TCPAddr{
				Zone: "test",
			}, nil
		})

		require.NoError(t, container.Compile())

		var addr *net.TCPAddr
		err := container.Populate(&addr)
		require.NoError(t, err)
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("function without arguments", func(t *testing.T) {
		var container = &Container{}

		// todo: improve error message
		container.Provide(func() {})
	})

	// todo: implement struct provide
	t.Run("struct", func(t *testing.T) {
		var container = &Container{}

		type StructProvider struct {
			TCPAddr *net.TCPAddr `inject:""`
			UDPAddr *net.UDPAddr `inject:""`
		}

		container.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{Zone: "tcp"}
		})

		container.Provide(func() *net.UDPAddr {
			return &net.UDPAddr{Zone: "udp"}
		})

		container.Provide(&StructProvider{})

		require.NoError(t, container.Compile())

		var sp *StructProvider
		require.NoError(t, container.Populate(&sp))
		require.Equal(t, "tcp", sp.TCPAddr.Zone)
		require.Equal(t, "udp", sp.UDPAddr.Zone)
	})
}

func TestBuilder_ProvideAs(t *testing.T) {
	t.Run("provide as", func(t *testing.T) {
		var container = &Container{}

		container.Provide(
			func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "test",
				}
			},
			As(new(net.Addr)),
		)

		require.NoError(t, container.Compile())

		var addr net.Addr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "test", addr.(*net.TCPAddr).Zone)
	})

	t.Run("provide as struct", func(t *testing.T) {
		var container = &Container{}

		container.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{}
		}, As(http.Server{}))
	})

	t.Run("provide as struct pointer", func(t *testing.T) {
		var container = &Container{}

		container.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{}
		}, As(new(http.Server)))
	})

	t.Run("provide as not implemented interface", func(t *testing.T) {
		var container = &Container{}

		container.Provide(func() *net.TCPAddr {
			return &net.TCPAddr{}
		}, As(new(http.Handler)))
	})
}

//
// func TestBuilder_ProvideName(t *testing.T) {
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
