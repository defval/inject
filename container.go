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
	// errIncorrectProviderType.
	errIncorrectProviderType = errors.New("value must be a function with value and optional error as result")

	// errIncorrectModifierSignature.
	errIncorrectModifierSignature = errors.New("modifier must be a function with optional error as result")
)

// errorInterface type for error interface implementation checking
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// New creates new container with provided options.
// Fore more information about container options see `Option` type.
func New(options ...Option) (_ *Container, err error) {
	var c = &Container{
		storage: &definitions{
			keys:            make([]key, 0, 8),
			definitions:     make(map[key]*definition, 8),
			implementations: make(map[key][]*definition, 8),
			groups:          make(map[key][]*definition, 8),
		},
	}

	// apply options.
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

// Container.
type Container struct {
	logger Logger

	providers []*providerOptions
	modifiers []*modifierOptions

	storage *definitions
}

// Populate populates given target pointer with type instance provided in container.
func (c *Container) Populate(target interface{}, options ...PopulateOption) (err error) {
	rv := reflect.ValueOf(target)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("populate target must be a not nil pointer")
	}

	if !rv.IsValid() {
		return errors.New("could not populate nil")
	}

	rv = rv.Elem()

	var po = &populateOptions{
		target: rv,
	}

	for _, opt := range options {
		opt.apply(po)
	}

	if rv.Kind() == reflect.Slice {
		return c.populateSlice(rv)
	}

	if err := c.populate(rv, po.name); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// populate
func (c *Container) populate(value reflect.Value, name string) (err error) {
	k := key{
		typ:  value.Type(),
		name: name,
	}

	var def *definition
	if def, err = c.storage.get(k); err != nil {
		return errors.WithStack(err)
	}

	instance, err := def.load()

	if err != nil {
		return errors.Wrapf(err, "%s", k)
	}

	value.Set(instance)

	return nil
}

func (c *Container) populateSlice(targetValue reflect.Value) (err error) {
	k := key{typ: targetValue.Type()}

	if _, ok := c.storage.groups[k]; !ok {
		return errors.Errorf("%s group not exists", targetValue.Type())
	}

	for _, def := range c.storage.groups[k] {
		instance, err := def.load()

		if err != nil {
			return errors.WithStack(err)
		}

		targetValue.Set(reflect.Append(targetValue, instance))
	}

	return nil
}

// compile.
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
		for _, k := range def.provider.args() {
			in, err := c.storage.get(k)

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

	c.storage.clearGroups()

	// apply modifiers
	for _, mo := range c.modifiers {
		if err = c.apply(mo); err != nil {
			return err
		}
	}

	return nil
}

// apply.
func (c *Container) apply(mo *modifierOptions) (err error) {
	if mo.modifier == nil {
		return errors.New("nil modifier")
	}

	// todo: validation
	mv := reflect.ValueOf(mo.modifier)
	mt := mv.Type()

	if err = checkModifier(mv); err != nil {
		return errors.WithStack(err)
	}

	var args []reflect.Value
	for i := 0; i < mt.NumIn(); i++ {
		rv := reflect.New(mt.In(i)).Elem()

		if rv.Kind() == reflect.Slice {
			if err = c.populateSlice(rv); err != nil {
				return errors.WithStack(err)
			}
		} else {
			if err = c.populate(rv, ""); err != nil {
				return errors.WithStack(err)
			}
		}

		args = append(args, rv)
	}

	var result = mv.Call(args)

	if len(result) == 1 {
		return errors.Wrap(result[0].Interface().(error), "apply error")
	}

	return nil
}

// checkModifier.
func checkModifier(mv reflect.Value) (err error) {
	if mv.Kind() != reflect.Func {
		return errors.WithStack(errIncorrectModifierSignature)
	}

	var modifierType = mv.Type()

	if modifierType.NumOut() > 1 {
		return errors.WithStack(errIncorrectModifierSignature)
	}

	if modifierType.NumOut() == 1 && !modifierType.Out(0).Implements(errorInterface) {
		return errors.WithStack(errIncorrectModifierSignature)
	}

	return nil
}

// providerOptions.
type providerOptions struct {
	provider   interface{}
	name       string
	implements []interface{}
}

// modifierOptions.
type modifierOptions struct {
	modifier interface{}
}

// populateOptions
type populateOptions struct {
	target reflect.Value
	name   string
}
