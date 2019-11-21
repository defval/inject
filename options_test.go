package inject

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestProvideOptions
func TestProvideOptions(t *testing.T) {
	opts := &providerOptions{}

	for _, opt := range []ProvideOption{
		WithName("test"),
		As(new(http.Handler)),
		Prototype(),
	} {
		opt.apply(opts)
	}

	require.Equal(t, &providerOptions{
		name:       "test",
		provider:   nil,
		interfaces: []interface{}{new(http.Handler)},
		prototype:  true,
	}, opts)
}

func TestExtractOptions(t *testing.T) {
	opts := &extractOptions{}

	for _, opt := range []ExtractOption{
		Name("test"),
	} {
		opt.apply(opts)
	}

	require.Equal(t, &extractOptions{
		name:   "test",
		target: nil,
	}, opts)
}
