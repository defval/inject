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
	// errIncorrectFunctionProviderSignature.
	errIncorrectFunctionProviderSignature = errors.New("constructor must be a function with value and optional error as result")

	// errIncorrectModifierSignature.
	errIncorrectModifierSignature = errors.New("modifier must be a function with optional error as result")
)

// errorInterface type for error interface implementation checking
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// New creates new container with provided options.
// Fore more information about container options see `Option` type.
func New(options ...Option) (_ *Container, err error) {
	var c = &Container{
		storage: &storage{
			keys:        make([]key, 0, 8),
			definitions: make(map[key]*definition, 8),
			ifaces:      make(map[reflect.Type][]*definition, 8),
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
	// modifiers []*modifierOptions

	storage *storage
}

// Populate populates given target pointer with type instance provided in container.
func (c *Container) Populate(target interface{}, options ...PopulateOption) (err error) {
	rv := reflect.ValueOf(target)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("populate target must be a not nil pointer")
	}

	rv = rv.Elem()

	var po = &populateOptions{
		target: rv,
	}

	for _, opt := range options {
		opt.apply(po)
	}

	k := key{
		typ:  po.target.Type(),
		name: po.name,
	}

	newValue, err := c.storage.Value(k)
	if err != nil {
		return errors.WithStack(err)
	}

	rv.Set(newValue)

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

		if err = c.storage.Add(def); err != nil {
			return errors.Wrap(err, "could not add definition")
		}
	}

	// connect storage
	for _, def := range c.storage.All() {
		// value arguments
		for _, argKey := range def.provider.args() {
			def.in = append(def.in, argKey)

			args, err := c.storage.Definition(argKey)

			if err != nil {
				return errors.WithStack(err)
			}

			for _, argDef := range args {
				argDef.out = append(argDef.out, def.key)
			}
		}
	}

	if err = c.storage.checkCycles(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//
// // apply.
// func (c *Container) apply(mo *modifierOptions) (err error) {
// 	if mo.modifier == nil {
// 		return errors.New("nil modifier")
// 	}
//
// 	// todo: validation
// 	mv := reflect.ValueOf(mo.modifier)
// 	mt := mv.Type()
//
// 	if err = checkModifier(mv); err != nil {
// 		return errors.WithStack(err)
// 	}
//
// 	var args []reflect.Value
// 	for i := 0; i < mt.NumIn(); i++ {
// 		rv := reflect.New(mt.In(i)).Elem()
//
// 		if rv.Kind() == reflect.Slice {
// 			if err = c.populateSlice(rv); err != nil {
// 				return errors.WithStack(err)
// 			}
// 		} else {
// 			if err = c.populate(rv, ""); err != nil {
// 				return errors.WithStack(err)
// 			}
// 		}
//
// 		args = append(args, rv)
// 	}
//
// 	var result = mv.Call(args)
//
// 	if len(result) == 1 {
// 		return errors.Wrap(result[0].Interface().(error), "apply error")
// 	}
//
// 	return nil
// }
//
// // checkModifier.
// func checkModifier(mv reflect.Value) (err error) {
// 	if mv.Kind() != reflect.Func {
// 		return errors.WithStack(errIncorrectModifierSignature)
// 	}
//
// 	var modifierType = mv.Type()
//
// 	if modifierType.NumOut() > 1 {
// 		return errors.WithStack(errIncorrectModifierSignature)
// 	}
//
// 	if modifierType.NumOut() == 1 && !modifierType.Out(0).Implements(errorInterface) {
// 		return errors.WithStack(errIncorrectModifierSignature)
// 	}
//
// 	return nil
// }

// providerOptions.
type providerOptions struct {
	provider           interface{}
	name               string
	implements         []interface{}
	injectPublicFields bool
}

//
// // modifierOptions.
// type modifierOptions struct {
// 	modifier interface{}
// }

// populateOptions
type populateOptions struct {
	target reflect.Value
	name   string
}
