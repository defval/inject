package provider

import (
	"reflect"

	"github.com/defval/inject/internal/graph"
)

// NewDirectProvider todo
func NewDirectProvider(rawProvider interface{}) *DirectProvider {
	return &DirectProvider{
		value: reflect.ValueOf(rawProvider),
	}
}

// DirectProvider todo
type DirectProvider struct {
	value reflect.Value
}

// Provide todo
func (p *DirectProvider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	return p.value, nil
}

// ResultType todo
func (p *DirectProvider) ResultType() reflect.Type {
	return p.value.Type()
}

// Arguments todo
func (p *DirectProvider) Arguments() (args []graph.Key) {
	return args
}
