package inject

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

const (
	nodeTypeProvider = iota + 1
	nodeTypeGroup
	nodeTypeBind
)

const (
	visitMarkUnmarked = iota
	visitMarkTemporary
	visitMarkPermanent
)

// newProvider
func newProvider(ctor interface{}) (_ *node, err error) {
	if ctor == nil {
		return nil, errors.New("nil could not be injected")
	}

	var ctype = reflect.TypeOf(ctor)

	if ctype.Kind() != reflect.Func {
		return nil, errors.Errorf("inject argument must be a function, got %s", ctype.String())
	}

	if ctype.NumOut() <= 0 || ctype.NumOut() > 2 {
		return nil, errors.Errorf("injection argument must be a function with returned value and optional error")
	}

	if ctype.NumOut() == 2 && ctype.Out(1).String() != "error" {
		return nil, errors.Errorf("injection argument must be a function with returned value and optional error")
	}

	var arguments = make([]reflect.Type, 0)

	for i := 0; i < ctype.NumIn(); i++ {
		arguments = append(arguments, ctype.In(i))
	}

	return &node{
		nodeType:   nodeTypeProvider,
		provider:   reflect.ValueOf(ctor),
		resultType: ctype.Out(0),
		args:       arguments,
	}, nil
}

func newGroup(iface interface{}, implementations ...interface{}) (_ *node, err error) {
	if iface == nil {
		return nil, errors.Errorf("group iface must be a interface pointer like new(http.Handler), got nil")
	}

	var ifaceType = reflect.TypeOf(iface)

	if ifaceType.Kind() != reflect.Ptr || ifaceType.Elem().Kind() != reflect.Interface {
		return nil, errors.Errorf("group iface must be a interface pointer like new(http.Handler), got %s", ifaceType)
	}

	var args []reflect.Type
	for _, member := range implementations {
		args = append(args, reflect.TypeOf(member))
	}

	return &node{
		nodeType:   nodeTypeGroup,
		resultType: reflect.SliceOf(ifaceType.Elem()),
		args:       args,
	}, nil
}

func newBind(iface interface{}, implementation interface{}) (_ *node, err error) {
	if iface == nil {
		return nil, errors.Errorf("bind iface must be a interface pointer like new(http.Handler), got nil")
	}

	var ifaceType = reflect.TypeOf(iface)

	if ifaceType.Kind() != reflect.Ptr || ifaceType.Elem().Kind() != reflect.Interface {
		return nil, errors.Errorf("bind iface must be a interface pointer like new(http.Handler), got %s", ifaceType)
	}

	var args []reflect.Type
	var implType = reflect.TypeOf(implementation)
	args = append(args, implType)

	return &node{
		nodeType:   nodeTypeBind,
		resultType: ifaceType.Elem(),
		args:       args,
	}, nil
}

// node
type node struct {
	visited  int
	nodeType int
	provider reflect.Value // only for nodes with provider type

	resultType reflect.Type
	args       []reflect.Type // arguments types

	in  []*node
	out []*node

	instance *reflect.Value
}

// addIn
func (n *node) addIn(node *node) {
	n.in = append(n.in, node)
}

// addOut
func (n *node) addOut(node *node) {
	n.out = append(n.out, node)
}

func (n *node) get(depth int) (value reflect.Value, err error) {
	if n.instance != nil {
		return *n.instance, nil
	}

	switch n.nodeType {
	case nodeTypeProvider:
		log.Print(Pad, strings.Repeat(LevelSymbol, depth), n.resultType.String())

		var values []reflect.Value
		for _, in := range n.in {
			var value reflect.Value
			if value, err = in.get(depth + 1); err != nil {
				return value, errors.Wrapf(err, "%s", in.resultType)
			}
			values = append(values, value)
		}

		var result = n.provider.Call(values)
		n.instance = &result[0]

		if len(result) == 2 {
			if result[1].IsNil() {
				return *n.instance, nil
			}

			return *n.instance, errors.WithStack(result[1].Interface().(error))
		}

		return *n.instance, nil
	case nodeTypeGroup:
		log.Print(Pad, strings.Repeat(LevelSymbol, depth), n.resultType.String())

		var values []reflect.Value
		for _, in := range n.in {
			var value reflect.Value
			if value, err = in.get(depth + 1); err != nil {
				return value, errors.Wrapf(err, "%s", in.resultType)
			}
			values = append(values, value)
		}

		elemSlice := reflect.MakeSlice(n.resultType, 0, 10)
		elemSlice = reflect.Append(elemSlice, values...)

		n.instance = &elemSlice

		return *n.instance, err
	case nodeTypeBind:
		log.Print(Pad, strings.Repeat(LevelSymbol, depth), n.resultType.String())

		if value, err = n.in[0].get(depth + 1); err != nil {
			return value, errors.Wrapf(err, "%s", n.in[0].resultType)
		}

		n.instance = &value

		return *n.instance, err
	}

	panic("unknown node type")
}

func (n *node) visit() (err error) {
	if n.visited == visitMarkPermanent {
		return
	}

	if n.visited == visitMarkTemporary {
		return fmt.Errorf("%s", n.resultType)
	}

	n.visited = visitMarkTemporary

	for _, outNode := range n.out {
		if err = outNode.visit(); err != nil {
			return errors.Wrapf(err, "%s", n.resultType)
		}
	}

	n.visited = visitMarkPermanent

	return nil
}

const Pad = " "
const LevelSymbol = "|  "
