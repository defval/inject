package di

import (
	"fmt"
	"reflect"

	"github.com/defval/inject/di/internal/dag"
	"github.com/defval/inject/di/internal/reflection"
)

// New create new container.
func New() *Container {
	return &Container{
		graph:     dag.NewDirectedGraph(),
		providers: make(map[providerKey]dependencyProvider),
	}
}

// Container is a dependency injection container.
type Container struct {
	graph     *dag.DirectedGraph
	providers map[providerKey]dependencyProvider
	compiled  bool
}

// ProvideParams params for Provide method.
type ProvideParams struct {
	Name        string
	Provider    interface{}
	Interfaces  []interface{}
	IsPrototype bool
}

// Provide provides given provider into container.
func (c *Container) Provide(params ProvideParams) {
	var provider dependencyProvider = createConstructor(params.Name, params.Provider)
	key := provider.Result()

	if c.graph.NodeExists(key) {
		panicf("The `%s` type already exists in container", provider.Result())
	}

	if !params.IsPrototype {
		provider = asSingleton(provider)
	}

	c.graph.AddNode(key)
	c.providers[key] = provider

	for _, iface := range params.Interfaces {
		c.provideAs(provider, iface)
	}
}

// Compile compiles the container. It iterates over all nodes
// in graph and register their parameters.
func (c *Container) Compile() {
	for _, key := range c.graph.Nodes() {
		// load provider parameters
		plist := c.providers[key.(providerKey)].Parameters()
		plist.Register(c, key.(providerKey))
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

	key := providerKey{
		Name: params.Name,
		Type: reflect.TypeOf(params.Target).Elem(),
	}

	return key.Extract(c, params.Target)
}

func (c *Container) provideAs(provider dependencyProvider, as interface{}) {
	// create interface from provider
	iface := createInterfaceProvider(provider, as)
	ifaceKey := iface.Result()

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
	groupKey := group.Result()

	// check exists
	if c.graph.NodeExists(groupKey) {
		// if exists use existing group
		group = c.providers[groupKey].(*interfaceGroup)
	} else {
		// else add new group to graph
		c.graph.AddNode(groupKey)
		c.providers[groupKey] = group
	}

	// add provider ifaceKey into group
	group.Add(provider.Result())
}
