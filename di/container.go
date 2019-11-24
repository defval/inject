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

// ProvideParams parameterList for Provide method.
type ProvideParams struct {
	Name        string
	Provider    interface{}
	Interfaces  []interface{}
	IsPrototype bool
}

// Provide adds constructor into container.
func (c *Container) Provide(params ProvideParams) {
	var provider provider = createConstructor(params.Name, params.Provider)
	id := provider.resultKey()

	if c.graph.NodeExists(id) {
		panicf("The `%s` type already exists in container", provider.resultKey())
	}

	if !params.IsPrototype {
		provider = asSingleton(provider)
	}

	c.graph.AddNode(id)
	c.providers[id] = provider

	for _, iface := range params.Interfaces {
		c.provideAs(provider, iface)
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

	for _, key := range c.all() {
		// register provider parameters
		provider, _ := c.provider(key)
		provider.parameters().register(c)
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
		return fmt.Errorf("extractInto target must be a pointer, got `nil`")
	}

	if !reflection.IsPtr(params.Target) {
		return fmt.Errorf("extractInto target must be a pointer, got `%s`", reflect.TypeOf(params.Target))
	}

	key := key{
		name: params.Name,
		typ:  reflect.TypeOf(params.Target).Elem(),
	}

	return key.extractInto(c, params.Target)
}

// exists checks that resultKey exists
func (c *Container) exists(k key) bool {
	return c.graph.NodeExists(k)
}

// provider
func (c *Container) provider(k key) (provider, error) {
	if !c.exists(k) {
		return nil, errProviderNotFound{k: k}
	}

	return c.providers[k], nil
}

// all
func (c *Container) all() []key {
	var keys []key

	for _, node := range c.graph.Nodes() {
		keys = append(keys, node.(key))
	}

	return keys
}

// registerDependency registers dependency
func (c *Container) registerDependency(dependency key, dependant key) {
	c.graph.AddEdge(dependency, dependant)
}

// provideAs
func (c *Container) provideAs(provider provider, as interface{}) {
	// create interface from structProvider
	iface := createInterfaceProvider(provider, as)
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
	group := createInterfaceGroup(ifaceKey)
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

	// add structProvider ifaceKey into group
	group.Add(provider.resultKey())
}
