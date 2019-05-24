package inject

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// storage.
type storage struct {
	keys        []key
	definitions map[key]*definition
	ifaces      map[reflect.Type][]*definition
}

// Add.
func (s *storage) Add(def *definition) (err error) {
	if _, ok := s.definitions[def.Key]; ok {
		return errors.Errorf("%s: use named definition if you have several instances of the same type", def.Key)
	}

	s.keys = append(s.keys, def.Key)
	s.definitions[def.Key] = def

	for _, typ := range def.Implements {
		s.ifaces[typ] = append(s.ifaces[typ], def)
	}

	return nil
}

// Replace
func (s *storage) Replace(def *definition) (err error) {
	if len(def.Implements) == 0 {
		return errors.Errorf("%s: no one interface has been replaced, use `inject.As()` for specify it", def.Key)
	}

	if _, ok := s.definitions[def.Key]; ok {
		*s.definitions[def.Key] = *def
		return nil
	}

	for _, typ := range def.Implements {
		k := key{
			typ:  typ,
			name: def.Key.name,
		}

		defs, err := s.Get(k)
		if err != nil {
			return errors.WithStack(err)
		}

		for i := range defs {
			*defs[i] = *def
		}
	}

	return nil
}

// Get.
func (s *storage) Get(k key) (_ []*definition, err error) {
	if def, ok := s.definitions[k]; ok {
		return []*definition{def}, nil
	}

	// definition iface
	if defs, ok := s.ifaces[k.typ]; ok {
		for _, def := range defs {
			if def.Key.name == k.name {
				return []*definition{def}, nil
			}
		}
	}

	// definition group
	if defs, ok := s.Group(k); ok {
		return defs, nil
	}

	return nil, errors.Errorf("type %s not provided", k)
}

// Group.
func (s *storage) Group(k key) (_ []*definition, exists bool) {
	if k.IsGroup() {
		_, ok := s.ifaces[k.typ.Elem()]
		return s.ifaces[k.typ.Elem()], ok
	}

	return nil, false
}

// Value.
func (s *storage) Value(k key) (v reflect.Value, err error) {
	defs, err := s.Get(k)

	if err != nil {
		return v, errors.WithStack(err)
	}

	v = k.Value()

	if len(defs) == 1 {
		var args []reflect.Value

		def := defs[0]

		for _, argKey := range def.In {
			arg, err := s.Value(argKey)

			if err != nil {
				return v, errors.WithStack(err)
			}

			args = append(args, arg)
		}

		instance, err := def.Create(args)
		if err != nil {
			return v, errors.Wrapf(err, "%s", def.Key)
		}

		v.Set(instance)

		return v, nil
	}

	// if k.IsGroup() {
	for _, def := range defs {
		var args []reflect.Value

		for _, argKey := range def.In {
			arg, err := s.Value(argKey)

			if err != nil {
				return v, errors.WithStack(err)
			}

			args = append(args, arg)
		}

		instance, err := def.Create(args)
		if err != nil {
			return v, errors.Wrapf(err, "%s", def.Key)
		}

		v.Set(reflect.Append(v, instance))
	}

	return v, nil
	// }

	// return v, errors.Errorf("type %s not provided", k)
}

// All.
func (s *storage) All() (defs []*definition) {
	for _, k := range s.keys {
		defs = append(defs, s.definitions[k])
	}

	return defs
}

// CheckCycles.
func (s *storage) CheckCycles() (err error) {
	// verify cycles
	for _, n := range s.All() {
		if n.visited == visitMarkUnmarked {
			if err = s.visit(n); err != nil {
				return errors.Wrap(err, "detect cycle")
			}
		}
	}

	return nil
}

// visit.
func (s *storage) visit(d *definition) (err error) {
	if d.visited == visitMarkPermanent {
		return
	}

	if d.visited == visitMarkTemporary {
		return fmt.Errorf("%s", d.Key)
	}

	d.visited = visitMarkTemporary

	for _, out := range d.Out {
		defs, _ := s.Get(out)
		// if err != nil {
		// 	return errors.WithStack(err)
		// }

		// visit arguments
		for _, def := range defs {
			if err = s.visit(def); err != nil {
				return errors.Wrapf(err, "%s", d.Key)
			}
		}
	}

	d.visited = visitMarkPermanent

	return nil
}
