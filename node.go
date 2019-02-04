package injector

import (
	"log"
	"reflect"
	"runtime"
	"strings"
)

// node
type node interface {
	Instance(level int) reflect.Value
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
func (n *provideNode) Instance(level int) reflect.Value {
	if n.instance != nil {
		return *n.instance
	}

	// var builder = strings.Builder{}
	//
	// builder.WriteString(n.resultType.String())
	//
	// for i, in := range n.in {
	// 	builder.WriteString("(")
	// 	builder.WriteString(in.Type().String())
	// 	if i != len(n.in)-1 {
	// 		builder.WriteString(", ")
	// 	}
	// 	builder.WriteString(")")
	// }
	//
	// log.Println(builder.String())

	log.Print(Pad, strings.Repeat(LevelSymbol, level), n.resultType.String())

	var values []reflect.Value
	for _, in := range n.in {
		values = append(values, in.Instance(level+1))
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
func (n *groupNode) Instance(level int) reflect.Value {
	if n.instance != nil {
		return *n.instance
	}

	log.Print(Pad, strings.Repeat(LevelSymbol, level), n.resultType.String())

	var values []reflect.Value
	for _, in := range n.in {
		values = append(values, in.Instance(level+1))
	}

	// var builder = strings.Builder{}
	// builder.WriteString(n.resultType.String())
	//
	// for i, in := range n.in {
	// 	builder.WriteString("(")
	// 	builder.WriteString(in.Type().String())
	// 	if i != len(n.in)-1 {
	// 		builder.WriteString(", ")
	// 	}
	// 	builder.WriteString(")")
	// }
	//
	// log.Println(builder.String())

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

func (n *bindNode) Instance(level int) reflect.Value {
	if n.instance != nil {
		return *n.instance
	}

	log.Print(Pad, strings.Repeat(LevelSymbol, level), n.resultType.String())

	var instance = n.in[0].Instance(level + 1)

	n.instance = &instance

	return *n.instance
}

const Pad = " "
const LevelSymbol = "|  "
