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
		providers: make(map[identity]dependencyProvider),
	}
}

// Container is a dependency injection container.
type Container struct {
	graph     *dag.DirectedGraph
	providers map[identity]dependencyProvider
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
	key := provider.Identity()

	if c.graph.NodeExists(key) {
		panicf("The `%s` type already exists in container", provider.Identity())
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
	// provide extractor
	c.Provide(ProvideParams{
		Provider: func() Extractor {
			return c
		},
	})

	for _, key := range c.graph.Nodes() {
		// load provider parameters
		plist := c.providers[key.(identity)].Parameters()
		plist.Register(c, key.(identity))
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

	key := identity{
		name: params.Name,
		typ:  reflect.TypeOf(params.Target).Elem(),
	}

	return key.Extract(c, params.Target)
}

func (c *Container) provideAs(provider dependencyProvider, as interface{}) {
	// create interface from provider
	iface := createInterfaceProvider(provider, as)
	ifaceKey := iface.Identity()

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
	groupKey := group.Identity()

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
	group.Add(provider.Identity())
}
