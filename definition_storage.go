package inject

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// storage
type storage struct {
	keys        []key
	definitions map[key]*definition            // key - type with name
	ifaces      map[reflect.Type][]*definition // key - interface
}

// Add
func (s *storage) Add(def *definition) (err error) {
	if _, ok := s.definitions[def.key]; ok {
		return errors.Errorf("%s: use named definition if you have several instances of the same type", def.key) // todo: value.String()
	}

	s.keys = append(s.keys, def.key)
	s.definitions[def.key] = def

	for _, typ := range def.implements {
		s.ifaces[typ] = append(s.ifaces[typ], def)
	}

	return nil
}

// Definition
// first: find in definitions
// second: find in interfaces with name
// third: find in group with elem of slice
func (s *storage) Definition(k key) (_ []*definition, err error) {
	if def, ok := s.definitions[k]; ok {
		return []*definition{def}, nil
	}

	// definition iface
	if defs, ok := s.ifaces[k.typ]; ok {
		for _, def := range defs {
			if def.key.name == k.name {
				return []*definition{def}, nil
			}
		}
	}

	// definition group
	if defs, ok := s.group(k); ok {
		return defs, nil
	}

	return nil, errors.Errorf("type %s not provided", k)
}

// Value
func (s *storage) Value(k key) (v reflect.Value, err error) {
	defs, err := s.Definition(k)

	if err != nil {
		return v, errors.WithStack(err)
	}

	v = k.Value()

	if len(defs) == 1 {
		var args []reflect.Value

		def := defs[0]

		for _, argKey := range def.in {
			arg, err := s.Value(argKey)

			if err != nil {
				return v, errors.WithStack(err)
			}

			args = append(args, arg)
		}

		instance, err := def.Create(args)
		if err != nil {
			return v, errors.Wrapf(err, "%s", def)
		}

		v.Set(instance)

		return v, nil
	}

	if k.IsGroup() {
		for _, def := range defs {
			var args []reflect.Value

			for _, argKey := range def.in {
				arg, err := s.Value(argKey)

				if err != nil {
					return v, errors.WithStack(err)
				}

				args = append(args, arg)
			}

			instance, err := def.Create(args)
			if err != nil {
				return v, errors.Wrapf(err, "%s", def)
			}

			v.Set(reflect.Append(v, instance))
		}

		return v, nil
	}

	return v, errors.Errorf("type %s not provided", k)
}

// groupExists
func (s *storage) group(k key) (_ []*definition, exists bool) {
	if k.IsGroup() {
		_, ok := s.ifaces[k.typ.Elem()]
		return s.ifaces[k.typ.Elem()], ok
	}

	return nil, false
}

// All
func (s *storage) All() (defs []*definition) {
	for _, k := range s.keys {
		defs = append(defs, s.definitions[k])
	}

	return defs
}

// checkCycles
func (s *storage) checkCycles() (err error) {
	// verify cycles
	for _, n := range s.All() {
		if n.visited == visitMarkUnmarked {
			if err = s.Visit(n); err != nil {
				return errors.Wrap(err, "detect cycle")
			}
		}
	}

	return nil
}

// Visit.
func (s *storage) Visit(d *definition) (err error) {
	if d.visited == visitMarkPermanent {
		return
	}

	if d.visited == visitMarkTemporary {
		return fmt.Errorf("%s", d.key)
	}

	d.visited = visitMarkTemporary

	for _, out := range d.out {
		defs, err := s.Definition(out)
		if err != nil {
			return errors.WithStack(err)
		}

		// visit arguments
		for _, def := range defs {
			if err = s.Visit(def); err != nil {
				return errors.Wrapf(err, "%s", d.key)
			}
		}
	}

	d.visited = visitMarkPermanent

	return nil
}
