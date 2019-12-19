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
	p := provider(newProviderConstructor(params.Name, params.Provider))
	if c.exists(p.Key()) {
		panicf("The `%s` type already exists in container", p.Key())
	}
	if !params.IsPrototype {
		p = asSingleton(p)
	}
	// add provider to graph
	c.add(p)
	// parse embed parameters
	for _, parameter := range p.ParameterList() {
		if parameter.embed {
			c.add(newProviderEmbed(parameter))
		}
	}
	// provide parameter bag
	if len(params.Parameters) != 0 {
		c.add(createParameterBugProvider(p.Key(), params.Parameters))
	}
	// process interfaces
	for _, iface := range params.Interfaces {
		c.processProviderInterface(p, iface)
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
	for _, p := range c.all() {
		c.registerProviderParameters(p)
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
	typ := reflect.TypeOf(params.Target)
	param := parameter{
		name:  params.Name,
		res:   typ.Elem(),
		embed: isEmbedParameter(typ),
	}
	value, err := param.ResolveValue(c)
	if err != nil {
		return err
	}
	targetValue := reflect.ValueOf(params.Target).Elem()
	targetValue.Set(value)
	return nil
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
	for _, p := range c.all() {
		if cleanup, ok := p.(cleanup); ok {
			cleanup.cleanup()
		}
	}
}

// add
func (c *Container) add(p provider) {
	c.graph.AddNode(p.Key())
	c.providers[p.Key()] = p
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
func (c *Container) all() []provider {
	var providers []provider
	for _, k := range c.graph.Nodes() {
		p, _ := c.provider(k.(key))
		providers = append(providers, p)
	}
	return providers
}

// processProviderInterface represents instances as interfaces and groups.
func (c *Container) processProviderInterface(provider provider, as interface{}) {
	// create interface from provider
	iface := newProviderInterface(provider, as)
	if c.graph.NodeExists(iface.Key()) {
		// if iface already exists, restrict interface resolving
		c.providers[iface.Key()] = newProviderStub(iface.Key(), "have several implementations")
	} else {
		// add interface node
		c.graph.AddNode(iface.Key())
		c.providers[iface.Key()] = iface
	}
	// create group
	group := newGroupProvider(iface.Key())
	// check exists
	if c.exists(group.Key()) {
		// if exists use existing group
		group = c.providers[group.Key()].(*interfaceGroup)
	} else {
		// else add new group to graph
		c.graph.AddNode(group.Key())
		c.providers[group.Key()] = group
	}
	// add embedParamProvider ifaceKey into group
	group.Add(provider.Key())
}

// registerProviderParameters registers provider parameters in a dependency graph.
func (c *Container) registerProviderParameters(p provider) {
	for _, param := range p.ParameterList() {
		paramProvider, exists := param.ResolveProvider(c)
		if exists {
			c.graph.AddEdge(paramProvider.Key(), p.Key())
			continue
		}
		if !exists && !param.optional {
			panicf("%s: dependency %s not exists in container", p.Key(), param)
		}
	}
}
