package injector

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlices(t *testing.T) {
	type iface interface{}

	var ifacetype = reflect.TypeOf(new(iface)).Elem()
	assert.Equal(t, "injector.iface", fmt.Sprintf("%s", ifacetype))
	// reflect.SliceOf(reflect.TypeOf(i)) // not work

	var fn = func(is []iface) {}

	var arg = reflect.TypeOf(fn).In(0).Elem()
	assert.Equal(t, "injector.iface", fmt.Sprintf("%s", arg))

	assert.Equal(t, ifacetype, arg)
}

type Writer struct {
}

func (w *Writer) Write(p []byte) (n int, err error) {
	panic("implement me")
}

func TestAppendPointerToInterfaceSlice(t *testing.T) {
	var s = make([]io.Writer, 0)

	var v = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(new(io.Writer)).Elem()), 0, 0)

	v = reflect.Append(v, reflect.ValueOf(new(Writer)))

	s = append(s, new(Writer))

	assert.Len(t, s, 1)
}
