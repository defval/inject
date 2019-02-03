package injector

import (
	"log"
	"reflect"
	"runtime"
)

// node
type node interface {
	Instance() reflect.Value
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
func newProvide(ctor interface{}) *provideNode {
	var ctype = reflect.TypeOf(ctor)
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
	}
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
func (n *provideNode) Instance() reflect.Value {
	if n.instance != nil {
		return *n.instance
	}

	log.Printf("Build: %s", n.resultType)

	var values []reflect.Value
	for _, in := range n.in {
		log.Printf("-- inject: %s", in.Type())
		values = append(values, in.Instance())
	}

	var result = n.cvalue.Call(values)
	n.instance = &result[0]

	return *n.instance
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

func newGroup(bindings ...interface{}) *groupNode {
	if len(bindings) < 2 {
		panic("need two types to bind")
	}

	var args []reflect.Type
	for _, impl := range bindings[1:] {
		var implType = reflect.TypeOf(impl)
		args = append(args, implType)
	}

	return &groupNode{
		baseNode: &baseNode{
			resultType: reflect.SliceOf(reflect.TypeOf(bindings[0]).Elem()),
			args:       args,
		},
	}
}

type groupNode struct {
	*baseNode
}

// Instance
func (n *groupNode) Instance() reflect.Value {
	if n.instance != nil {
		return *n.instance
	}

	log.Printf("Group: %s", n.resultType)

	var values []reflect.Value
	for _, in := range n.in {
		log.Printf("-- inject: %s", in.Type())
		values = append(values, in.Instance())
	}

	elemSlice := reflect.MakeSlice(n.resultType, 0, 10)
	elemSlice = reflect.Append(elemSlice, values...)

	n.instance = &elemSlice

	return *n.instance
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

func (n *bindNode) Instance() reflect.Value {
	if n.instance != nil {
		return *n.instance
	}

	log.Printf("Bind: %s", n.resultType)
	log.Printf("-- inject: %s", n.in[0].Type())

	var instance = n.in[0].Instance()

	n.instance = &instance

	return *n.instance
}
