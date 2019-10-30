package dot_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/defval/inject/internal/dag"
	"github.com/defval/inject/internal/dag/dot"
)

func TestFromGraph(t *testing.T) {
	g := dag.New()

	foo := dag.NewNode("foo", "foo")
	bar := dag.NewNode("bar", "bar")

	_ = g.Add(foo)
	_ = g.Add(bar)

	_ = foo.ConnectWith(bar)

	{
		foo := dag.NewNode("first_foo", "foo")
		bar := dag.NewNode("first_bar", "bar")

		sub := g.Subgraph("first")
		_ = sub.Add(foo)
		_ = sub.Add(bar)

		{
			foo2 := dag.NewNode("second_foo", "foo")
			bar2 := dag.NewNode("second_bar", "bar")

			_ = foo2.ConnectWith(bar2)
		}
	}

	dg := dot.FromGraph(g)

	var buffer bytes.Buffer
	_, _ = dg.WriteTo(&buffer)

	fmt.Print(buffer.String())
}
