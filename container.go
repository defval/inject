package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	// ErrIncorrectProviderType
	ErrIncorrectProviderType = errors.New("provider must be a function with value and optional error as result")

	// ErrIncorrectModifierSignature
	ErrIncorrectModifierSignature = errors.New("modifier must be a function with optional error as result")
)

// errorInterface type for error interface implementation checking
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// New creates new container with provided options.
func New(options ...Option) (_ *Container, err error) {
	var container = &Container{
		nodes:           make(map[key]*definition),
		implementations: make(map[key][]*definition),
	}

	for _, opt := range options {
		opt.apply(container)
	}

	if container.logger == nil {
		container.logger = &defaultLogger{}
	}

	if err = container.compile(); err != nil {
		return nil, errors.Wrapf(err, "could not compile container")
	}

	return container, nil
}

// Container
type Container struct {
	logger          Logger
	providers       []*providerOptions
	modifiers       []*modifierOptions
	nodes           map[key]*definition
	implementations map[key][]*definition
}

// Populate
func (c *Container) Populate(target interface{}, options ...ProvideOption) (err error) {
	var targetValue = reflect.ValueOf(target).Elem()

	var def *definition
	if def, err = c.get(key{typ: targetValue.Type()}); err != nil {
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
func (c *Container) add(def *definition) (err error) {
	if _, ok := c.nodes[def.key]; ok {
		return errors.Wrapf(err, "%s already provided", def.provider.resultType)
	}

	c.nodes[def.key] = def

	for _, key := range def.implements {
		c.implementations[key] = append(c.implementations[key], def)
	}

	c.logger.Printf("Provide: %s", def)

	return nil
}

// get
func (c *Container) get(key key) (_ *definition, err error) {
	if def, ok := c.nodes[key]; ok {
		return def, nil
	}

	if len(c.implementations[key]) > 0 {
		return c.implementations[key][0], nil // todo: return element
	}

	return nil, errors.Errorf("%s not provided yet", key)
}

// compile
func (c *Container) compile() (err error) {
	// register providers
	for _, po := range c.providers {
		if po.provider == nil {
			return errors.New("could not provide nil")
		}

		var def *definition
		if def, err = po.definition(); err != nil {
			return errors.Wrapf(err, "provide failed")
		}

		if err = c.add(def); err != nil {
			return errors.Wrap(err, "could not add node")
		}
	}

	// connect nodes
	for _, def := range c.nodes {
		// load arguments
		for _, key := range def.provider.args {
			in, err := c.get(key)

			if err != nil {
				return errors.WithStack(err)
			}

			def.in = append(def.in, in)
		}
	}

	// apply modifiers
	for _, mo := range c.modifiers {
		if err = mo.apply(c); err != nil {
			return err
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

// modifierOptions
type modifierOptions struct {
	modifier interface{}
}

// apply
func (o *modifierOptions) apply(container *Container) (err error) {
	if o.modifier == nil {
		return errors.New("nil modifier")
	}

	// todo: validation
	var modifierValue = reflect.ValueOf(o.modifier)

	if modifierValue.Kind() != reflect.Func {
		return errors.WithStack(ErrIncorrectModifierSignature)
	}

	var modifierType = modifierValue.Type()

	if modifierType.NumOut() > 1 {
		return errors.WithStack(ErrIncorrectModifierSignature)
	}

	if modifierType.NumOut() == 1 && !modifierType.Out(0).Implements(errorInterface) {
		return errors.WithStack(ErrIncorrectModifierSignature)
	}

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

	var result = modifierValue.Call(args)

	if len(result) == 1 {
		return errors.Wrap(result[0].Interface().(error), "apply error")
	}

	return nil
}
