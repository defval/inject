package inject

import (
	"fmt"
	"log"
	"reflect"

	"github.com/pkg/errors"
)

// New creates new container
func New(options ...Option) (_ *Injector, err error) {
	var injector = &Injector{
		nodes: &nodeStorage{
			keys:  []reflect.Type{},
			nodes: map[reflect.Type]*oldNode{},
		},
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

	if err = injector.verifyCycles(); err != nil {
		return nil, errors.Wrap(err, "cycle not allowed")
	}

	log.Printf("BUILDED")

	return injector, nil
}

// Injector ...
type Injector struct {
	providers []interface{}
	bindings  []*bind
	groups    []*group

	nodes *nodeStorage
}

// Populate
func (i *Injector) Populate(targets ...interface{}) (err error) {
	for _, target := range targets {
		var targetValue = reflect.ValueOf(target).Elem()

		var node *oldNode
		if node, err = i.get(targetValue.Type()); err != nil {
			return errors.WithStack(err)
		}

		var instance reflect.Value
		if instance, err = node.get(0); err != nil {
			return errors.Wrapf(err, "%s", targetValue.Type())
		}

		targetValue.Set(instance)
	}

	return nil
}

func (i *Injector) processProviders() (err error) {
	for _, provider := range i.providers {
		var provide *oldNode
		if provide, err = newProvider(provider); err != nil {
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
		var bind *oldNode
		if bind, err = newBind(binding.iface, binding.implementation); err != nil {
			return errors.WithStack(err)
		}

		if err = i.add(bind); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (i *Injector) processGroups() (err error) {
	for _, group := range i.groups {
		var node *oldNode
		if node, err = newGroup(group.iface, group.implementations...); err != nil {
			return errors.WithStack(err)
		}

		if err = i.add(node); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (i *Injector) connectNodes() (err error) {
	for _, node := range i.nodes.all() {
		for _, arg := range node.args {
			arg, err := i.get(arg)

			if err != nil {
				return errors.WithStack(err)
			}

			for _, currentOut := range arg.out {
				if currentOut == node {
					return fmt.Errorf("%v already injected into %v", arg.resultType, node.resultType)
				}
			}

			arg.addOut(node)
			node.addIn(arg)
		}
	}

	return nil
}

func (i *Injector) add(node *oldNode) (err error) {
	log.Printf("INJECT: %s", node.resultType)

	if err = i.nodes.add(node); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (i *Injector) get(typ reflect.Type) (node *oldNode, _ error) {
	var found bool
	if node, found = i.nodes.get(typ); !found {
		return nil, fmt.Errorf("%s not found", typ)
	}

	return node, nil
}

func (i *Injector) verifyCycles() (err error) {
	for _, n := range i.nodes.all() {
		if n.visited == visitMarkUnmarked {
			if err = n.visit(); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

// group
type group struct {
	iface           interface{}
	implementations []interface{}
}

type bind struct {
	iface          interface{}
	implementation interface{}
}
