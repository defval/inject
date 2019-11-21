package inject_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defval/inject"
)

func TestContainer(t *testing.T) {
	var HTTPBundle = inject.Bundle(
		inject.Provide(ProvideAddr("0.0.0.0", "8080")),
		inject.Provide(NewMux, inject.As(new(http.Handler))),
		inject.Provide(NewHTTPServer, inject.Prototype(), inject.WithName("server")),
	)

	c := inject.New(HTTPBundle)

	var server1 *http.Server
	err := c.Extract(&server1, inject.Name("server"))
	require.NoError(t, err)

	var server2 *http.Server
	err = c.Extract(&server2, inject.Name("server"))
	require.NoError(t, err)
}

// Addr
type Addr string

// ProvideAddr
func ProvideAddr(host string, port string) func() Addr {
	return func() Addr {
		return Addr(net.JoinHostPort(host, port))
	}
}

// NewHTTPServer
func NewHTTPServer(addr Addr, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    string(addr),
		Handler: handler,
	}
}

// NewMux
func NewMux() *http.ServeMux {
	return &http.ServeMux{}
}
