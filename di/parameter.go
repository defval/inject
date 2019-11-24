package di

import (
	"fmt"
	"reflect"
)

// parameter
type parameter interface {
	resultKey() key
	register(c *Container, dependant key)
	load(c *Container) (reflect.Value, error)
}

// parameterRequired
type parameterRequired struct{ key }

// parameterOptional
type parameterOptional struct{ key }

// parameterEmbed
type parameterEmbed struct{ key }

// parameterList
type parameterList []parameter

// load loads all parameters presented in parameter list.
func (pl parameterList) load(c *Container) ([]reflect.Value, error) {
	var values []reflect.Value
	for _, p := range pl {
		pvalue, err := p.load(c)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", p.resultKey(), err)
		}

		values = append(values, pvalue)
	}

	return values, nil
}

// providerParameterList
type providerParameterList struct {
	parameterList
	providerKey key
}

// add adds parameter into provider parameter list
func (l *providerParameterList) add(p parameter) {
	l.parameterList = append(l.parameterList, p)
}

// register registers parameters in container
func (l providerParameterList) register(c *Container) {
	for _, p := range l.parameterList {
		// register parameter in container
		p.register(c, l.providerKey)
	}
}
