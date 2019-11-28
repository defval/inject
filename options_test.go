package inject

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defval/inject/v2/di"
)

// TestProvideOptions
func TestProvideOptions(t *testing.T) {
	opts := &di.ProvideParams{
		Parameters: map[string]interface{}{},
	}

	for _, opt := range []ProvideOption{
		WithName("test"),
		As(new(http.Handler)),
		Prototype(),
		ParameterBag{
			"test": "test",
		},
	} {
		opt.apply(opts)
	}

	require.Equal(t, &di.ProvideParams{
		Name:        "test",
		Provider:    nil,
		Interfaces:  []interface{}{new(http.Handler)},
		IsPrototype: true,
		Parameters: map[string]interface{}{
			"test": "test",
		},
	}, opts)
}

func TestExtractOptions(t *testing.T) {
	opts := &di.ExtractParams{}

	for _, opt := range []ExtractOption{
		Name("test"),
	} {
		opt.apply(opts)
	}

	require.Equal(t, &di.ExtractParams{
		Name:   "test",
		Target: nil,
	}, opts)
}
