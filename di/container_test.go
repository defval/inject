package di_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defval/inject/v2/di"
	"github.com/defval/inject/v2/di/internal/ditest"
)

func TestContainerCompileErrors(t *testing.T) {
	t.Run("dependency cycle cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewCycleFooBar)
		c.MustProvide(ditest.NewBar)
		c.MustCompileError("Cycle detected")
	})

	t.Run("not existing dependency cause compile error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewBar)
		c.MustCompileError("*ditest.Bar: dependency *ditest.Foo not exists in container")
	})
}

func TestContainerProvideErrors(t *testing.T) {
	t.Run("provide string cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError("string", "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `string`")
	})

	t.Run("provide nil cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(nil, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `nil`")
	})

	t.Run("provide struct pointer cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(&ditest.Foo{}, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `*ditest.Foo`")
	})

	t.Run("provide constructor without result cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(ditest.ConstructorWithoutResult, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `github.com/defval/inject/v2/di/internal/ditest.ConstructorWithoutResult`")
	})

	t.Run("provide constructor with many results cause panic", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(ditest.ConstructorWithManyResults, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `github.com/defval/inject/v2/di/internal/ditest.ConstructorWithManyResults`")
	})

	t.Run("provide constructor with incorrect result error argument", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvideError(ditest.ConstructorWithIncorrectResultError, "The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `github.com/defval/inject/v2/di/internal/ditest.ConstructorWithIncorrectResultError`")
	})

	t.Run("provide duplicate", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvideError(ditest.NewFoo, "The `*ditest.Foo` type already exists in container")
	})

	t.Run("provide as not implemented interface cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvideError(ditest.NewBar, "*ditest.Bar not implement ditest.Barer", new(ditest.Barer))
	})

	t.Run("provide as not interface cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvideError(ditest.NewBar, "*ditest.Foo: not a pointer to interface", new(ditest.Foo))
	})
}

func TestContainerExtractErrors(t *testing.T) {
	t.Run("container panic on trying extract before compilation", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		var extracted *ditest.Foo
		c.MustExtractError(&extracted, "container not compiled")
	})

	t.Run("extract into string cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustCompile()
		c.MustExtractError("string", "extract target must be a pointer, got `string`")
	})

	t.Run("extract into struct cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustCompile()
		c.MustExtractError(struct{}{}, "extract target must be a pointer, got `struct {}`")
	})

	t.Run("extract into nil cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustCompile()
		c.MustExtractError(nil, "extract target must be a pointer, got `nil`")
	})

	t.Run("container does not find type because its named", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvideWithName("foo", ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtractError(&extracted, "type `*ditest.Foo` not exists in container")
	})

	t.Run("extract returns error because dependency constructing failed", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.CreateFooConstructorWithError(errors.New("internal error")))
		c.MustProvide(ditest.NewBar)
		c.MustCompile()
		var bar *ditest.Bar
		c.MustExtractError(&bar, "*ditest.Foo: internal error")
	})

	t.Run("extract interface with multiple implementations cause error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.NewBar, new(ditest.Fooer))
		c.MustProvide(ditest.NewBaz, new(ditest.Fooer))
		c.MustCompile()

		var extracted ditest.Fooer
		c.MustExtractError(&extracted, "ditest.Fooer have sereral implementations")
	})
}

func TestContainerProvide(t *testing.T) {
	t.Run("container successfully accept simple constructor", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
	})

	t.Run("container successfully accept constructor with error", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.CreateFooConstructorWithError(nil))
	})

	t.Run("container successfully accept constructor with cleanup function", func(t *testing.T) {
		c := NewTestContainer(t)

		cleanup := func() {}
		c.MustProvide(ditest.CreateFooConstructorWithCleanup(cleanup))
	})

}

func TestContainerExtract(t *testing.T) {
	t.Run("container extract correct pointer", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtractPtr(foo, &extracted)
	})

	t.Run("container extract same pointer on each extraction", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted1 *ditest.Foo
		c.MustExtractPtr(foo, &extracted1)

		var extracted2 *ditest.Foo
		c.MustExtractPtr(foo, &extracted2)
	})

	t.Run("container extract instance if error is nil", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.CreateFooConstructorWithError(nil))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtract(&extracted)
	})

	t.Run("container extract instance if cleanup and error is nil", func(t *testing.T) {
		c := NewTestContainer(t)

		c.MustProvide(ditest.CreateFooConstructorWithCleanupAndError(nil, nil))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtract(&extracted)
	})

	t.Run("container extract correct named pointer", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvideWithName("foo", ditest.CreateFooConstructor(foo))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtractWithName("foo", &extracted)
	})

	t.Run("container extract correct interface implementation", func(t *testing.T) {
		c := NewTestContainer(t)
		bar := &ditest.Bar{}
		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.CreateBarConstructor(bar), new(ditest.Fooer))
		c.MustCompile()

		var extracted ditest.Fooer
		c.MustExtractPtr(bar, &extracted)
	})

	t.Run("container creates group from interface and extract it", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.NewBar, new(ditest.Fooer))
		c.MustProvide(ditest.NewBaz, new(ditest.Fooer))
		c.MustCompile()

		var group []ditest.Fooer
		c.MustExtract(&group)
		require.Len(t, group, 2)
	})

	t.Run("container extract new instance of prototype by each extraction", func(t *testing.T) {
		c := NewTestContainer(t)
		c.MustProvide(ditest.NewFoo)
		c.MustProvidePrototype(ditest.NewBar)
		c.MustCompile()

		var extracted1 *ditest.Bar
		c.MustExtract(&extracted1)
		var extracted2 *ditest.Bar
		c.MustExtract(&extracted2)

		c.MustNotEqualPointer(extracted1, extracted2)
	})

	t.Run("container extract correct parameter bag for type", func(t *testing.T) {
		c := NewTestContainer(t)

		c.Provide(di.ProvideParams{
			Provider: ditest.NewFooWithName,
			Parameters: di.ParameterBag{
				"name": "test",
			},
		})

		c.MustCompile()

		var foo *ditest.Foo
		err := c.Extract(di.ExtractParams{
			Target: &foo,
		})

		require.NoError(t, err)
		require.Equal(t, "test", foo.Name)
	})

	t.Run("container extract correct parameter bag for type", func(t *testing.T) {
		c := NewTestContainer(t)

		c.Provide(di.ProvideParams{
			Name:     "named",
			Provider: ditest.NewFooWithName,
			Parameters: di.ParameterBag{
				"name": "test",
			},
		})

		c.MustCompile()

		var foo *ditest.Foo
		err := c.Extract(di.ExtractParams{
			Name:   "named",
			Target: &foo,
		})

		require.NoError(t, err)
		require.Equal(t, "test", foo.Name)
	})
}

func TestContainerResolve(t *testing.T) {
	t.Run("container resolve correct argument", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := &ditest.Foo{}
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.NewBar)
		c.MustCompile()

		var bar *ditest.Bar
		c.MustExtract(&bar)
		c.MustEqualPointer(foo, bar.Foo())
	})

	t.Run("container resolve correct interface implementation", func(t *testing.T) {
		c := NewTestContainer(t)

		foo := ditest.NewFoo()
		bar := ditest.NewBar(foo)

		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.CreateBarConstructor(bar), new(ditest.Fooer))
		c.MustProvide(ditest.NewQux)
		c.MustCompile()

		var qux *ditest.Qux
		c.MustExtract(&qux)
		c.MustEqualPointer(bar, qux.Fooer())
	})

	t.Run("container resolve correct group", func(t *testing.T) {
		c := NewTestContainer(t)

		c.MustProvide(ditest.NewFoo)
		c.MustProvide(ditest.NewBar, new(ditest.Fooer))
		c.MustProvide(ditest.NewBaz, new(ditest.Fooer))
		c.MustProvide(ditest.NewFooerGroup)
		c.MustCompile()

		var bar *ditest.Bar
		c.MustExtract(&bar)

		var baz *ditest.Baz
		c.MustExtract(&baz)

		var group *ditest.FooerGroup
		c.MustExtract(&group)
		require.Len(t, group.Fooers(), 2)
		c.MustEqualPointer(bar, group.Fooers()[0])
		c.MustEqualPointer(baz, group.Fooers()[1])
	})
}

func TestContainerResolveEmbedParameters(t *testing.T) {
	t.Run("container resolve embed parameters", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := ditest.NewFoo()
		bar := ditest.NewBar(foo)
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.CreateBarConstructor(bar))
		c.MustProvide(ditest.NewBazFromParameters)
		c.MustCompile()

		var extracted *ditest.Baz
		c.MustExtract(&extracted)
		c.MustEqualPointer(foo, extracted.Foo())
		c.MustEqualPointer(bar, extracted.Bar())
	})

	t.Run("container skip optional parameter", func(t *testing.T) {
		c := NewTestContainer(t)
		foo := ditest.NewFoo()
		c.MustProvide(ditest.CreateFooConstructor(foo))
		c.MustProvide(ditest.NewBazFromParameters)
		c.MustCompile()

		var extracted *ditest.Baz
		c.MustExtract(&extracted)
		c.MustEqualPointer(foo, extracted.Foo())
		require.Nil(t, extracted.Bar())
	})

	t.Run("container resolve optional not existing group as nil", func(t *testing.T) {
		c := NewTestContainer(t)
		type Params struct {
			di.Parameter
			Handlers []http.Handler `di:"optional"`
		}
		c.MustProvide(func(params Params) bool {
			return params.Handlers == nil
		})
		c.MustCompile()
		var extracted bool
		c.MustExtract(&extracted)
		require.True(t, extracted)
	})

	t.Run("container skip private fields in parameter", func(t *testing.T) {
		c := NewTestContainer(t)
		type Param struct {
			di.Parameter
			private    []http.Handler `di:"optional"`
			Addrs      []net.Addr     `di:"optional"`
			HaveNotTag string
		}
		c.MustProvide(func(param Param) bool {
			return param.Addrs == nil
		})
		c.MustCompile()
		var extracted bool
		c.MustExtract(&extracted)
		require.True(t, extracted)
	})
}

func TestContainerCleanup(t *testing.T) {
	t.Run("container run cleanup function after container close", func(t *testing.T) {
		c := NewTestContainer(t)

		var cleanupCalled bool
		cleanup := func() {
			cleanupCalled = true
		}

		c.MustProvide(ditest.CreateFooConstructorWithCleanup(cleanup))
		c.MustCompile()

		var extracted *ditest.Foo
		c.MustExtract(&extracted)
		c.Cleanup()

		require.True(t, cleanupCalled)
	})
}

// NewTestContainer
func NewTestContainer(t *testing.T) *TestContainer {
	return &TestContainer{t, di.New()}
}

// TestContainer
type TestContainer struct {
	t *testing.T
	*di.Container
}

func (c *TestContainer) MustProvide(provider interface{}, as ...interface{}) {
	require.NotPanics(c.t, func() {
		c.Provide(di.ProvideParams{
			Provider:   provider,
			Interfaces: as,
		})
	}, "provide should not panic")
}

func (c *TestContainer) MustProvidePrototype(provider interface{}, as ...interface{}) {
	require.NotPanics(c.t, func() {
		c.Provide(di.ProvideParams{
			Provider:    provider,
			Interfaces:  as,
			IsPrototype: true,
		})
	})
}

func (c *TestContainer) MustProvideWithName(name string, provider interface{}, as ...interface{}) {
	require.NotPanics(c.t, func() {
		c.Provide(di.ProvideParams{
			Name:       name,
			Provider:   provider,
			Interfaces: as,
		})
	})
}

func (c *TestContainer) MustProvideError(provider interface{}, msg string, as ...interface{}) {
	require.PanicsWithValue(c.t, msg, func() {
		c.Provide(di.ProvideParams{
			Provider:   provider,
			Interfaces: as,
		})
	})
}

func (c *TestContainer) MustCompile() {
	require.NotPanics(c.t, func() {
		c.Compile()
	})
}

func (c *TestContainer) MustCompileError(msg string) {
	require.PanicsWithValue(c.t, msg, func() {
		c.Compile()
	})
}

func (c *TestContainer) MustExtract(target interface{}) {
	require.NoError(c.t, c.Extract(di.ExtractParams{
		Target: target,
	}))
}

func (c *TestContainer) MustExtractWithName(name string, target interface{}) {
	require.NoError(c.t, c.Extract(di.ExtractParams{
		Name:   name,
		Target: target,
	}))
}

func (c *TestContainer) MustExtractError(target interface{}, msg string) {
	require.EqualError(c.t, c.Extract(di.ExtractParams{
		Target: target,
	}), msg)
}

func (c *TestContainer) MustExtractWithNameError(name string, target interface{}, msg string) {
	require.EqualError(c.t, c.Extract(di.ExtractParams{
		Name:   name,
		Target: target,
	}), msg)
}

// MustExtractPtr extract value from container into target and check that target and expected pointers are equal.
func (c *TestContainer) MustExtractPtr(expected, target interface{}) {
	c.MustExtract(target)

	// indirect
	actual := reflect.ValueOf(target).Elem().Interface()
	c.MustEqualPointer(expected, actual)
}

func (c *TestContainer) MustExtractPtrWithName(expected interface{}, name string, target interface{}) {
	c.MustExtractWithName(name, target)

	actual := reflect.ValueOf(target).Elem().Interface()
	c.MustEqualPointer(expected, actual)
}

func (c *TestContainer) MustEqualPointer(expected interface{}, actual interface{}) {
	require.Equal(c.t,
		fmt.Sprintf("%p", actual),
		fmt.Sprintf("%p", expected),
		"actual and expected pointers should be equal",
	)
}

func (c *TestContainer) MustNotEqualPointer(expected interface{}, actual interface{}) {
	require.NotEqual(c.t,
		fmt.Sprintf("%p", actual),
		fmt.Sprintf("%p", expected),
		"actual and expected pointers should not be equal",
	)
}
