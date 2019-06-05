package object_test

import (
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defval/inject/internal/provider/object"
)

// StructWithOnlyFields
type StructDependency struct {
	privateField  string
	Server        *http.Server `inject:"test"`
	privateField2 string
	Mux           *http.ServeMux
	privateField3 string
	TCPAddr       *net.TCPAddr `inject:""`
	privateField4 string
	privateField5 string
	UDPAddr       *net.UDPAddr
	privateField6 string
	AnotherServer *http.Server `anotherTag:"another"`
}

func TestStructPointerProvider_Arguments(t *testing.T) {
	t.Run("all fields with tag are arguments", func(t *testing.T) {
		p, _ := object.New(&StructDependency{})

		args := p.Arguments()
		require.Len(t, args, 2)
		require.Equal(t, "*http.Server", args[0].Type.String())
		require.Equal(t, "test", args[0].Name)
		require.Equal(t, "*net.TCPAddr", args[1].Type.String())
		require.Equal(t, "", args[1].Name)
	})

	t.Run("with provider.Exported() option all public fields are arguments", func(t *testing.T) {
		p, _ := object.New(&StructDependency{}, object.Exported())

		args := p.Arguments()
		require.Len(t, args, 5)
		require.Equal(t, "*http.Server", args[0].Type.String())
		require.Equal(t, "test", args[0].Name)
		require.Equal(t, "*http.ServeMux", args[1].Type.String())
		require.Equal(t, "", args[1].Name)
		require.Equal(t, "*net.TCPAddr", args[2].Type.String())
		require.Equal(t, "", args[2].Name)
		require.Equal(t, "*net.UDPAddr", args[3].Type.String())
		require.Equal(t, "", args[3].Name)
		require.Equal(t, "*http.Server", args[4].Type.String())
		require.Equal(t, "", args[4].Name)
	})

	t.Run("change tag works correctly", func(t *testing.T) {
		p, _ := object.New(&StructDependency{}, object.Tag("anotherTag"))
		args := p.Arguments()
		require.Len(t, args, 1)
		require.Equal(t, "*http.Server", args[0].Type.String())
	})
}

func TestStructPointerProvider_Provide(t *testing.T) {
	t.Run("provide", func(t *testing.T) {
		p, _ := object.New(&StructDependency{})

		server := &http.Server{}
		addr := &net.TCPAddr{}

		args := []reflect.Value{
			reflect.ValueOf(server),
			reflect.ValueOf(addr),
		}

		v, err := p.Provide(args)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%p", server), fmt.Sprintf("%p", v.Interface().(*StructDependency).Server))
	})
}
