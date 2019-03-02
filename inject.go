package inject

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// New creates new container
func New(options ...Option) (_ *Injector, err error) {
	var injector = &Injector{
		nodes: make(map[reflect.Type]node),
	}

	for _, opt := range options {
		opt.apply(injector)
	}

	if err = injector.processProviders(); err != nil {
		return nil, errors.WithStack(err)
	}

	if err = injector.processBindings(); err != nil {
		return nil, errors.WithStack(err)
	}

	if err = injector.processGroups(); err != nil {
		return nil, errors.WithStack(err)
	}

	if err = injector.connectNodes(); err != nil {
		return nil, errors.WithStack(err)
	}

	log.Printf("BUILDED")
	log.Println()

	return injector, nil
}

// Injector ...
type Injector struct {
	lock sync.Mutex

	providers []interface{}
	bindings  [][]interface{}
	groups    []*group

	nodes map[reflect.Type]node
}

// Populate
func (i *Injector) Populate(targets ...interface{}) (err error) {
	for _, target := range targets {
		var targetValue = reflect.ValueOf(target).Elem()

		var node node
		if node, err = i.get(targetValue.Type()); err != nil {
			return errors.WithStack(err)
		}

		var instance reflect.Value
		if instance, err = node.Instance(0); err != nil {
			return errors.Wrapf(err, "%s", targetValue.Type())
		}

		targetValue.Set(instance)

		log.Println()
	}

	return nil
}

func (i *Injector) processProviders() (err error) {
	for _, provider := range i.providers {
		var provide *provideNode
		if provide, err = newProvide(provider); err != nil {
			return errors.WithStack(err)
		}

		if err = i.add(provide); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (i *Injector) processBindings() (err error) {
	for _, binding := range i.bindings {
		if len(binding) == 2 {
			var bind = newBind(binding[0], binding[1])

			if err = i.add(bind); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

func (i *Injector) processGroups() (err error) {
	for _, group := range i.groups {
		var node *groupNode
		if node, err = newGroup(group.of, group.members...); err != nil {
			return errors.WithStack(err)
		}

		if err = i.add(node); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (i *Injector) connectNodes() (err error) {
	for _, node := range i.nodes {
		for _, arg := range node.Args() {
			arg, err := i.get(arg)

			if err != nil {
				return errors.WithStack(err)
			}

			if err = i.connect(arg, node); err != nil {
				return errors.WithStack(err)
			}
		}
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

func (i *Injector) get(typ reflect.Type) (node node, _ error) {
	var found bool
	if node, found = i.nodes[typ]; !found {
		return nil, fmt.Errorf("%s not found", typ)
	}

	return node, nil
}

// func (i *Injector) out(n *provideNode) ([]node, error) {
// 	var successors []node
//
// 	_, found := i.get(n.resultType)
// 	if found != nil {
// 		return successors, fmt.Errorf("%s not found", n.resultType)
// 	}
//
// 	for _, v := range n.out {
// 		successors = append(successors, v)
// 	}
//
// 	return successors, nil
// }
//
// func (i *Injector) in(n *provideNode) ([]node, error) {
// 	var predecessors []node
//
// 	_, found := i.get(n.resultType)
// 	if found != nil {
// 		return predecessors, fmt.Errorf("%s not found", n.resultType)
// 	}
//
// 	for _, v := range n.in {
// 		predecessors = append(predecessors, v)
// 	}
//
// 	return predecessors, nil
// }

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
