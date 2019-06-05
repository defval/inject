package direct

import (
	"reflect"

	"github.com/defval/inject/internal/provider"
)

// New
func New(rawProvider interface{}) *Provider {
	return &Provider{
		value: reflect.ValueOf(rawProvider),
	}
}

// Provider
type Provider struct {
	value reflect.Value
}

func (p *Provider) Provide(arguments []reflect.Value) (reflect.Value, error) {
	return p.value, nil
}

func (p *Provider) ResultType() reflect.Type {
	return p.value.Type()
}

func (p *Provider) Arguments() (args []provider.Key) {
	return args
}
