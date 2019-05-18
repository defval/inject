package inject

import (
	"reflect"

	"github.com/pkg/errors"
)

const (
	visitMarkUnmarked = iota
	visitMarkTemporary
	visitMarkPermanent
)

var (
	// ErrIncorrectProviderType
	ErrIncorrectProviderType = errors.New("value must be a function with value and optional error as result")

	// ErrIncorrectModifierSignature
	ErrIncorrectModifierSignature = errors.New("modifier must be a function with optional error as result")
)

// errorInterface type for error interface implementation checking
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// New creates new container with provided options.
func New(options ...Option) (_ *Container, err error) {
	var c = &Container{
		storage: &definitions{
			keys:            make([]key, 0, 8),
			definitions:     make(map[key]*definition, 8),
			implementations: make(map[key][]*definition, 8),
		},
	}

	for _, opt := range options {
		opt.apply(c)
	}

	if c.logger == nil {
		c.logger = &defaultLogger{}
	}

	if err = c.compile(); err != nil {
		return nil, errors.Wrapf(err, "could not compile container")
	}

	return c, nil
}

// Container
type Container struct {
	logger Logger

	providers []*providerOptions
	modifiers []*modifierOptions

	storage *definitions
}

// Populate
func (c *Container) Populate(target interface{}, options ...ProvideOption) (err error) {
	rvalue := reflect.ValueOf(target)

	if !rvalue.IsValid() || rvalue.IsNil() {
		return errors.New("could not populate nil")
	}

	rvalue = rvalue.Elem()

	var def *definition
	if def, err = c.storage.get(key{typ: rvalue.Type()}); err != nil {
		return errors.WithStack(err)
	}

	var instance reflect.Value
	if instance, err = def.instance(); err != nil {
		return errors.Wrapf(err, "%s", rvalue.Type())
	}

	rvalue.Set(instance)

	return nil
}

// compile
func (c *Container) compile() (err error) {
	// register providers
	for _, po := range c.providers {
		if po.provider == nil {
			return errors.New("could not provide nil")
		}

		var def *definition
		if def, err = createDefinition(po); err != nil {
			return errors.Wrapf(err, "provide failed")
		}

		if err = c.storage.add(def); err != nil {
			return errors.Wrap(err, "could not add definition")
		}
	}

	// connect definitions
	for _, def := range c.storage.all() {
		// load arguments
		for _, key := range def.provider.arguments {
			in, err := c.storage.get(key)

			if err != nil {
				return errors.WithStack(err)
			}

			def.in = append(def.in, in)
			in.out = append(in.out, def)
		}
	}

	// verify cycles
	for _, n := range c.storage.all() {
		if n.visited == visitMarkUnmarked {
			if err = n.visit(); err != nil {
				return errors.Wrap(err, "detect cycle")
			}
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

// modifierOptions
type modifierOptions struct {
	modifier interface{}
}

// apply
func (o *modifierOptions) apply(c *Container) (err error) {
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
		if def, err = c.storage.get(key{typ: modifierType.In(i)}); err != nil {
			return errors.WithStack(err)
		}

		var arg reflect.Value
		if arg, err = def.instance(); err != nil {
			return errors.Wrapf(err, "%s", def)
		}

		args = append(args, arg)
	}

	var result = modifierValue.Call(args)

	if len(result) == 1 {
		return errors.Wrap(result[0].Interface().(error), "apply error")
	}

	return nil
}
