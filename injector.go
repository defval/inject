package injector

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

// New creates new container
func New(options ...Option) (_ *Injector, err error) {
	var c = &Injector{
		nodes: make(map[reflect.Type]node),
	}

	for _, opt := range options {
		opt.apply(c)
	}

	for _, provider := range c.providers {
		if err = c.add(newProvide(provider)); err != nil {
			return nil, err
		}
	}

	for _, binding := range c.binders {
		if len(binding) == 2 {
			if err = c.add(newBind(binding[0], binding[1])); err != nil {
				return nil, err
			}
		} else {
			if err = c.add(newGroup(binding...)); err != nil {
				return nil, err
			}
		}
	}

	for _, n := range c.nodes {
		for _, arg := range n.Args() {
			tail, err := c.getNode(arg)

			if err != nil {
				return nil, err
			}

			if err = c.connect(tail, n); err != nil {
				return nil, err
			}
		}
	}

	log.Printf("BUILDED")
	log.Println()

	return c, nil
}

// Injector ...
type Injector struct {
	lock      sync.Mutex
	providers []interface{}
	binders   [][]interface{}

	nodes map[reflect.Type]node
}

// Populate
func (i *Injector) Populate(targets ...interface{}) (err error) {
	for _, target := range targets {
		var v = reflect.ValueOf(target).Elem()

		var node node
		if node, err = i.getNode(v.Type()); err != nil {
			return fmt.Errorf("could not populate `%s`", v.Type())
		}

		var instance = node.Instance(0)

		v.Set(instance)
	}

	return nil
}

func (i *Injector) add(node node) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	log.Printf("INJECT: %s", node.Type())

	if _, ok := i.nodes[node.Type()]; ok {
		return fmt.Errorf("%s already injected", node.Type())
	}

	i.nodes[node.Type()] = node

	return nil
}

func (i *Injector) connect(n1, n2 node) error {
	dependencyExist := false
	nodeExists := false

	i.lock.Lock()
	defer i.lock.Unlock()

	for _, cur := range i.nodes {
		if cur == n1 {
			dependencyExist = true
		}
		if cur == n2 {
			nodeExists = true
		}
	}

	if !dependencyExist {
		return fmt.Errorf("%s not found", n1.Type())
	}
	if !nodeExists {
		return fmt.Errorf("%s not found", n2.Type())
	}

	for _, n := range n1.Out() {
		if n == n2 {
			return fmt.Errorf("%v already injected in to %v", n1.Type(), n2.Type())
		}
	}

	n1.AddOut(n2)
	n2.AddIn(n1)

	return nil
}

func (i *Injector) getNode(typ reflect.Type) (node node, _ error) {
	var found bool
	if node, found = i.nodes[typ]; !found {
		return nil, fmt.Errorf("%s not found", typ)
	}

	return node, nil
}

func (i *Injector) out(n *provideNode) ([]node, error) {
	var successors []node

	_, found := i.getNode(n.resultType)
	if found != nil {
		return successors, fmt.Errorf("%s not found", n.resultType)
	}

	for _, v := range n.out {
		successors = append(successors, v)
	}

	return successors, nil
}

// predecessors return vertices that are in of a given vertex.
func (i *Injector) in(n *provideNode) ([]node, error) {
	var predecessors []node

	_, found := i.getNode(n.resultType)
	if found != nil {
		return predecessors, fmt.Errorf("%s not found", n.resultType)
	}

	for _, v := range n.in {
		predecessors = append(predecessors, v)
	}

	return predecessors, nil
}

//
// // String implements stringer interface.
// //
// // Prints an string representation of this Instance.
// func (c *Injector) String() string {
// 	// resultType := fmt.Sprintf("DAG Vertices: %c - Edges: %c\n", c.Order(), c.Size())
//
// 	var s []node
// 	for _, node := range c.nodes {
// 		s = append(s, node)
// 	}
//
// 	sort.Slice(s, func(i, j int) bool {
// 		return s[i].instability() < s[j].instability()
// 	})
//
// 	var result string
// 	result += fmt.Sprintf("Nodes:\n")
// 	for _, node := range s {
// 		result += fmt.Sprintf("%s", node)
// 	}
//
// 	return result
// }
