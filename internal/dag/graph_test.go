package dag_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defval/inject/internal/dag"
)

func TestNew(t *testing.T) {
	g := dag.New()

	require.NoError(t, g.Add("foo"), "must add foo node")
	require.NoError(t, g.Add("bar"), "must add bar node")

	foo, err := g.Node("foo")
	require.NoError(t, err, "must find node")
	require.Equal(t, "foo:value", foo.Value().(string), "loaded value equals")

	bar, err := g.Node("bar")
	require.NoError(t, err, "must find node")
	require.Equal(t, "bar:value", bar.Value().(string), "loaded value equals")

	require.NoError(t, foo.ConnectWith(bar), "must add ancestor")

	require.Len(t, bar.Ancestors(), 1)
	require.Len(t, foo.Descendants(), 1)
}
