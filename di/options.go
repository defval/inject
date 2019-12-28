package di

// ExtractOption
type ProvideOption interface {
	apply(params *ProvideParams)
}

type provideOption func(params *ProvideParams)

func (o provideOption) apply(params *ProvideParams) {
	o(params)
}

// ProvideParams is a `Provide()` method options. Name is a unique identifier of type instance. Provider is a constructor
// function. Interfaces is a interface that implements a provider result type.
type ProvideParams struct {
	Name        string
	Interfaces  []interface{}
	Parameters  ParameterBag
	IsPrototype bool
}

func (p ProvideParams) apply(params *ProvideParams) {
	*params = p
}

// As
func As(interfaces ...interface{}) ProvideOption {
	return provideOption(func(params *ProvideParams) {
		params.Interfaces = append(params.Interfaces, interfaces...)
	})
}

// InvokeParams is a invoke parameters.
type InvokeParams struct{}

func (p InvokeParams) apply(params *InvokeParams) {
	*params = p
}

// InvokeOption
type InvokeOption interface {
	apply(params *InvokeParams)
}

// ExtractParams
type ExtractParams struct {
	Name string
}

func (p ExtractParams) apply(params *ExtractParams) {
	*params = p
}

// ExtractOption
type ExtractOption interface {
	apply(params *ExtractParams)
}

type extractOption func(params *ExtractParams)

func (o extractOption) apply(params *ExtractParams) {
	o(params)
}
