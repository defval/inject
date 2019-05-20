package inject

import (
	"github.com/pkg/errors"
)

// definitions
type definitions struct {
	keys            []key
	definitions     map[key]*definition
	implementations map[key][]*definition
	groups          map[key][]*definition
}

// add
func (s *definitions) add(def *definition) (err error) {
	if _, ok := s.definitions[def.key]; ok {
		return errors.Errorf("%s: use named definition if you have several instances of the same type", def.key) // todo: value.String()
	}

	s.keys = append(s.keys, def.key)
	s.definitions[def.key] = def

	for _, k := range def.implements {
		s.implementations[k] = append(s.implementations[k], def)

		if _, ok := s.groups[k]; !ok {
			s.groups[k] = make([]*definition, 0, 8)
		}

		groupKey := createGroupKey(k)

		s.groups[groupKey] = append(s.groups[groupKey], def)
	}

	return nil
}

// get
func (s *definitions) get(k key) (_ *definition, err error) {
	if def, ok := s.definitions[k]; ok {
		return def, nil
	}

	if len(s.implementations[k]) > 0 {
		return s.implementations[k][0], nil // todo: return element
	}

	return nil, errors.Errorf("type %s not provided", k)
}

// all
func (s *definitions) all() (defs []*definition) {
	for _, k := range s.keys {
		defs = append(defs, s.definitions[k])
	}

	return defs
}

func (s *definitions) clearGroups() {
	for k, group := range s.groups {
		if len(group) == 1 {
			delete(s.groups, k)
		}
	}
}
