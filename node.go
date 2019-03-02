package inject

import (
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// node
type node interface {
	Instance(level int) (value reflect.Value, err error)
	Type() reflect.Type
	Args() []reflect.Type
	AddIn(node node)
	AddOut(node node)
	In() []node
	Out() []node
}

type baseNode struct {
	node
	resultType reflect.Type
	args       []reflect.Type // arguments types

	in  []node
	out []node

	instance *reflect.Value
}

// Type
func (n *baseNode) Type() reflect.Type {
	return n.resultType
}

// Args
func (n *baseNode) Args() []reflect.Type {
	return n.args
}

// AddIn
func (n *baseNode) AddIn(node node) {
	n.in = append(n.in, node)
}

// AddOut
func (n *baseNode) AddOut(node node) {
	n.out = append(n.out, node)
}

// In
func (n *baseNode) In() []node {
	return n.in
}

// Out
func (n *baseNode) Out() []node {
	return n.out
}

// newProvide
func newProvide(ctor interface{}) (_ *provideNode, err error) {
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

	var cvalue = reflect.ValueOf(ctor)
	var cptr = cvalue.Pointer()
	var cfunc = runtime.FuncForPC(cvalue.Pointer())
	var cname = cfunc.Name()
	var arguments = make([]reflect.Type, 0)

	for i := 0; i < ctype.NumIn(); i++ {
		arguments = append(arguments, ctype.In(i))
	}

	return &provideNode{
		ctype:  ctype,
		cvalue: reflect.ValueOf(ctor),
		cptr:   cptr,
		cfunc:  cfunc,
		cname:  cname,
		baseNode: &baseNode{
			resultType: ctype.Out(0),
			args:       arguments,
		},
	}, nil
}

// provideNode
type provideNode struct {
	ctype  reflect.Type
	cvalue reflect.Value
	cptr   uintptr
	cfunc  *runtime.Func
	cname  string

	*baseNode
}

// Instance
func (n *provideNode) Instance(level int) (_ reflect.Value, err error) {
	if n.instance != nil {
		return *n.instance, nil
	}

	log.Print(Pad, strings.Repeat(LevelSymbol, level), n.resultType.String())

	var values []reflect.Value
	for _, in := range n.in {
		var value reflect.Value
		if value, err = in.Instance(level + 1); err != nil {
			return value, errors.Wrapf(err, "%s", in.Type())
		}
		values = append(values, value)
	}

	var result = n.cvalue.Call(values)
	n.instance = &result[0]

	if len(result) == 2 {
		if result[1].IsNil() {
			return *n.instance, nil
		}

		return *n.instance, errors.WithStack(result[1].Interface().(error))
	}

	return *n.instance, nil
}

//
// func (n *provideNode) instability() float64 {
// 	if len(n.in) == 0 && len(n.out) == 0 {
// 		return -1
// 	}
//
// 	return float64(len(n.in) / (len(n.in) + len(n.out)))
// }
//
// func (n *provideNode) String() string {
// 	result := fmt.Sprintf("%s in: %d out: %d instability: %.2f\n", n.resultType.String(), len(n.in), len(n.out), n.instability())
//
// 	return result
// }

func newGroup(of interface{}, members ...interface{}) (_ *groupNode, err error) {
	if of == nil {
		return nil, errors.Errorf("group of must be a interface pointer like new(http.Handler)")
	}

	var args []reflect.Type
	for _, member := range members {
		args = append(args, reflect.TypeOf(member))
	}

	return &groupNode{
		baseNode: &baseNode{
			resultType: reflect.SliceOf(reflect.TypeOf(of).Elem()),
			args:       args,
		},
	}, nil
}

type groupNode struct {
	*baseNode
}

// Instance
func (n *groupNode) Instance(level int) (_ reflect.Value, err error) {
	if n.instance != nil {
		return *n.instance, err
	}

	log.Print(Pad, strings.Repeat(LevelSymbol, level), n.resultType.String())

	var values []reflect.Value
	for _, in := range n.in {
		var value reflect.Value
		if value, err = in.Instance(level + 1); err != nil {
			return value, errors.Wrapf(err, "%s", in.Type())
		}
		values = append(values, value)
	}

	elemSlice := reflect.MakeSlice(n.resultType, 0, 10)
	elemSlice = reflect.Append(elemSlice, values...)

	n.instance = &elemSlice

	return *n.instance, err
}

func newBind(target interface{}, source interface{}) *bindNode {
	var args []reflect.Type
	var implType = reflect.TypeOf(source)
	args = append(args, implType)

	return &bindNode{
		baseNode: &baseNode{
			resultType: reflect.TypeOf(target).Elem(),
			args:       args,
		},
	}
}

type bindNode struct {
	*baseNode
}

func (n *bindNode) Instance(level int) (value reflect.Value, err error) {
	if n.instance != nil {
		return *n.instance, err
	}

	log.Print(Pad, strings.Repeat(LevelSymbol, level), n.resultType.String())

	if value, err = n.in[0].Instance(level + 1); err != nil {
		return value, errors.Wrapf(err, "%s", n.in[0].Type())
	}

	n.instance = &value

	return *n.instance, err
}

const Pad = " "
const LevelSymbol = "|  "
