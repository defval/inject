package inject

import (
	"log"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrIncorrectProviderType = errors.New("provider must be a function with returned value and optional error")
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// New
func New(options ...Option) (_ *Container, err error) {
	var container = &Container{}

	for _, opt := range options {
		opt.apply(container)
	}

	if err = container.compile(); err != nil {
		return nil, errors.Wrapf(err, "could not compile container")
	}

	return container, nil
}

// Container
type Container struct {
	init sync.Once

	providers       []*providerOptions
	modifiers       []*modifierOptions
	nodes           map[key]*definition
	implementations map[key][]*definition
}

// Populate
func (b *Container) Populate(target interface{}, options ...ProvideOption) (err error) {
	var targetValue = reflect.ValueOf(target).Elem()

	var def *definition
	if def, err = b.get(key{typ: targetValue.Type()}); err != nil {
		return errors.WithStack(err)
	}

	var instance reflect.Value
	if instance, err = def.instance(); err != nil {
		return errors.Wrapf(err, "%s", targetValue.Type())
	}

	targetValue.Set(instance)

	return nil
}

// add
func (b *Container) add(def *definition) (err error) {
	b.init.Do(func() {
		b.nodes = make(map[key]*definition, 8)
		b.implementations = make(map[key][]*definition, 8)
	})

	if _, ok := b.nodes[def.key]; ok {
		return errors.Wrapf(err, "%s already provided", def.provider.resultType)
	}

	b.nodes[def.key] = def

	for _, key := range def.implements {
		b.implementations[key] = append(b.implementations[key], def)
	}

	log.Printf("Provide: %s", def)

	return nil
}

func (b *Container) get(key key) (_ *definition, err error) {
	if def, ok := b.nodes[key]; ok {
		return def, nil
	}

	if len(b.implementations[key]) > 0 {
		return b.implementations[key][0], nil // todo: return element
	}

	return nil, errors.Errorf("%s not provided yet", key)
}

// Build
func (b *Container) compile() (err error) {
	// register providers
	for _, po := range b.providers {
		if po.provider == nil {
			return errors.New("could not provide nil")
		}

		var def *definition
		if def, err = po.definition(); err != nil {
			return errors.Wrapf(err, "provide failed")
		}

		if err = b.add(def); err != nil {
			return errors.Wrap(err, "could not add node")
		}
	}

	// connect nodes
	for _, def := range b.nodes {
		// load arguments
		for _, key := range def.provider.args {
			in, err := b.get(key)

			if err != nil {
				return errors.WithStack(err)
			}

			def.in = append(def.in, in)
		}
	}

	// apply modifiers
	for _, mo := range b.modifiers {
		if err = mo.apply(b); err != nil {
			return errors.Wrap(err, "could not apply modifier")
		}
	}

	return nil
}

// providerOptions
type providerOptions struct {
	provider   interface{}
	name       string
	implements []interface{}
}

// definition
func (o *providerOptions) definition() (_ *definition, err error) {
	ptype := reflect.TypeOf(o.provider)

	var wrapper *providerWrapper
	switch true {
	case ptype.Kind() == reflect.Func:
		wrapper, err = newFuncProvider(o.provider)
	case ptype.Kind() == reflect.Ptr && ptype.Elem().Kind() == reflect.Struct:
		wrapper, err = newStructProvider(o.provider)
	default:
		return nil, errors.WithStack(ErrIncorrectProviderType)
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var implements []key
	for _, iface := range o.implements {
		ifaceType := reflect.TypeOf(iface)

		if ifaceType.Kind() != reflect.Ptr || ifaceType.Elem().Kind() != reflect.Interface {
			return nil, errors.Errorf("argument for As() must be pointer to interface type, got %s", ifaceType)
		}

		ifaceTypeElem := ifaceType.Elem()

		if !wrapper.resultType.Implements(ifaceTypeElem) {
			return nil, errors.Errorf("%s not implement %s interface", wrapper.resultType, ifaceTypeElem)
		}

		implements = append(implements, key{typ: ifaceTypeElem})
	}

	return &definition{
		key: key{
			typ:  wrapper.resultType,
			name: o.name,
		},
		provider:   wrapper,
		implements: implements,
	}, nil
}

type modifierOptions struct {
	modifier interface{}
}

func (o *modifierOptions) apply(container *Container) (err error) {
	if o.modifier == nil {
		return errors.New("nil modifier")
	}

	// todo: validation
	var modifierValue = reflect.ValueOf(o.modifier)
	var modifierType = modifierValue.Type()

	var args []reflect.Value
	for i := 0; i < modifierType.NumIn(); i++ {
		// todo: add name
		var def *definition
		if def, err = container.get(key{typ: modifierType.In(i)}); err != nil {
			return errors.WithStack(err)
		}

		var arg reflect.Value
		if arg, err = def.instance(); err != nil {
			return errors.WithStack(err)
		}

		args = append(args, arg)
	}

	modifierValue.Call(args)

	return nil
}
