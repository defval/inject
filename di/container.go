package di

import (
	"fmt"
	"reflect"

	"github.com/defval/inject/v2/di/internal/dag"
	"github.com/defval/inject/v2/di/internal/reflection"
)

// New create new container.
func New() *Container {
	return &Container{
		graph:     dag.NewDirectedGraph(),
		providers: make(map[key]provider),
	}
}

// Container is a dependency injection container.
type Container struct {
	graph     *dag.DirectedGraph
	providers map[key]provider
	compiled  bool
}

// ProvideParams is a `Provide()` method options. Name is a unique identifier of type instance. Provider is a constructor
// function. Interfaces is a interface that implements a provider result type.
type ProvideParams struct {
	Name        string
	Provider    interface{}
	Interfaces  []interface{}
	Parameters  ParameterBag
	IsPrototype bool
}

// Provide adds constructor into container with parameters.
func (c *Container) Provide(params ProvideParams) {
	prov := provider(newConstructorProvider(params.Name, params.Provider))
	k := prov.resultKey()

	if c.exists(k) {
		panicf("The `%s` type already exists in container", prov.resultKey())
	}

	if !params.IsPrototype {
		prov = asSingleton(prov)
	}

	c.addProvider(prov)
	c.provideEmbedParameters(prov)

	if len(params.Parameters) != 0 {
		parameterBugProvider := createParameterBugProvider(k, params.Parameters)
		c.addProvider(parameterBugProvider)
	}

	for _, iface := range params.Interfaces {
		c.processProviderInterface(prov, iface)
	}
}

// Compile compiles the container. It iterates over all nodes
// in graph and register their parameters.
func (c *Container) Compile() {
	// provide extractor
	c.Provide(ProvideParams{
		Provider: func() Extractor {
			return c
		},
	})

	c.Provide(ProvideParams{
		Provider: func() *Graph {
			return &Graph{graph: c.graph.DOTGraph()}
		},
	})

	for _, key := range c.all() {
		// register provider parameters
		provider, _ := c.provider(key)
		c.registerParameters(provider)
	}

	_, err := c.graph.DFSSort()
	if err != nil {
		switch err {
		case dag.ErrCyclicGraph:
			panicf("Cycle detected") // todo: add nodes to message
		default:
			panic(err.Error())
		}
	}

	c.compiled = true
}

// ExtractParams
type ExtractParams struct {
	Name   string
	Target interface{}
}

// Extract
func (c *Container) Extract(params ExtractParams) error {
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}

	if params.Target == nil {
		return fmt.Errorf("extract target must be a pointer, got `nil`")
	}

	if !reflection.IsPtr(params.Target) {
		return fmt.Errorf("extract target must be a pointer, got `%s`", reflect.TypeOf(params.Target))
	}

	key := key{
		name: params.Name,
		typ:  reflect.TypeOf(params.Target).Elem(),
	}

	return key.extract(c, params.Target)
}

// InvokeParams
type InvokeParams struct {
	// The function
	Fn interface{}
}

// Invoke calls provided function.
func (c *Container) Invoke(params InvokeParams) error {
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}

	invoker, err := newInvoker(params.Fn)
	if err != nil {
		return err
	}

	return invoker.Invoke(c)
}

// Cleanup
func (c *Container) Cleanup() {
	for _, key := range c.all() {
		provider, _ := c.provider(key)
		if cleanup, ok := provider.(cleanup); ok {
			cleanup.cleanup()
		}
	}
}

// addProvider
func (c *Container) addProvider(p provider) {
	c.graph.AddNode(p.resultKey())
	c.providers[p.resultKey()] = p
}

// provideEmbedParameters
func (c *Container) provideEmbedParameters(p provider) {
	for _, parameter := range p.parameters() {
		if parameter.embed {
			c.addProvider(newEmbedProvider(parameter))
		}
	}
}

// exists checks that key registered in container graph.
func (c *Container) exists(k key) bool {
	return c.graph.NodeExists(k)
}

// provider checks that provider exists and return it.
func (c *Container) provider(k key) (provider, bool) {
	if !c.exists(k) {
		return nil, false
	}

	return c.providers[k], true
}

// all return all container keys.
func (c *Container) all() []key {
	var keys []key

	for _, node := range c.graph.Nodes() {
		keys = append(keys, node.(key))
	}

	return keys
}

// processProviderInterface represents instances as interfaces and groups.
func (c *Container) processProviderInterface(provider provider, as interface{}) {
	// create interface from embedParamProvider
	iface := newInterfaceProvider(provider, as)
	ifaceKey := iface.resultKey()

	if c.graph.NodeExists(ifaceKey) {
		// if iface already exists, restrict interface resolving
		c.providers[ifaceKey] = iface.Multiple()
	} else {
		// add interface node
		c.graph.AddNode(ifaceKey)
		c.providers[ifaceKey] = iface
	}

	// create group
	group := newGroupProvider(ifaceKey)
	groupKey := group.resultKey()

	// check exists
	if c.graph.NodeExists(groupKey) {
		// if exists use existing group
		group = c.providers[groupKey].(*interfaceGroup)
	} else {
		// else add new group to graph
		c.graph.AddNode(groupKey)
		c.providers[groupKey] = group
	}

	// add embedParamProvider ifaceKey into group
	group.Add(provider.resultKey())
}

// registerParameters registers provider parameters in a dependency graph.
func (c *Container) registerParameters(provider provider) {
	for _, parameter := range provider.parameters() {
		_, exists := c.provider(parameter.resultKey())
		if exists {
			c.graph.AddEdge(parameter.resultKey(), provider.resultKey())
		}

		if !exists && !parameter.optional {
			panicf("%s: dependency %s not exists in container", provider.resultKey(), parameter.resultKey())
		}
	}
}
