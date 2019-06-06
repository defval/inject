package provider

import (
	"reflect"

	"github.com/defval/inject/internal/graph"
)

// NewConstructorProvider
func NewDirectProvider(rawProvider interface{}) *DirectProvider {
	return &DirectProvider{
		value: reflect.ValueOf(rawProvider),
	}
}

// ObjectProvider
type DirectProvider struct {
	value reflect.Value
}

func (p *DirectProvider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	return p.value, nil
}

func (p *DirectProvider) ResultType() reflect.Type {
	return p.value.Type()
}

func (p *DirectProvider) Arguments() (args []graph.Key) {
	return args
}
