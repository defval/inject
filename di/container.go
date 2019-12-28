package di

import (
	"fmt"
	"reflect"

	"github.com/defval/inject/v2/di/internal/graphkv"
	"github.com/defval/inject/v2/di/internal/reflection"
)

// Interactor is a helper interface.
type Interactor interface {
	Extract(target interface{}, options ...ExtractOption) error
	Invoke(fn interface{}, options ...InvokeOption) error
}

// Builder is helper interface.
type Builder interface {
	Provide(provider interface{}, options ...ProvideOption)
}

// New create new container.
func New() *Container {
	return &Container{
		graph: graphkv.New(),
	}
}

// Container is a dependency injection container.
type Container struct {
	compiled bool
	graph    *graphkv.Graph
	cleanups []func()
}

// Provide adds constructor into container with parameters.
func (c *Container) Provide(constructor interface{}, options ...ProvideOption) {
	params := ProvideParams{}
	for _, opt := range options {
		opt.apply(&params)
	}
	provider := internalProvider(newProviderConstructor(params.Name, constructor))
	key := provider.Key()
	if c.graph.Exists(key) {
		panicf("The `%s` type already exists in container", provider.Key())
	}
	if !params.IsPrototype {
		provider = asSingleton(provider)
	}
	// add provider to graph
	c.graph.Add(key, provider)
	// parse embed parameters
	for _, param := range provider.ParameterList() {
		if param.embed {
			embed := newProviderEmbed(param)
			c.graph.Add(embed.Key(), embed)
		}
	}
	// provide parameter bag
	if len(params.Parameters) != 0 {
		parameterBugProvider := createParameterBugProvider(provider.Key(), params.Parameters)
		c.graph.Add(parameterBugProvider.Key(), parameterBugProvider)
	}
	// process interfaces
	for _, iface := range params.Interfaces {
		c.processProviderInterface(provider, iface)
	}
}

// Compile compiles the container. It iterates over all nodes
// in graph and register their parameters.
func (c *Container) Compile() {
	graphProvider := func() *Graph { return &Graph{graph: c.graph.DOTGraph()} }
	interactorProvider := func() Interactor { return c }
	c.Provide(graphProvider)
	c.Provide(interactorProvider)
	for _, node := range c.graph.Nodes() {
		c.registerProviderParameters(node.Value.(internalProvider))
	}
	if err := c.graph.CheckCycles(); err != nil {
		panic(err.Error())
	}
	c.compiled = true
}

// Extract builds instance of target type and fills target pointer.
func (c *Container) Extract(target interface{}, options ...ExtractOption) error {
	params := ExtractParams{}
	for _, opt := range options {
		opt.apply(&params)
	}
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}
	if target == nil {
		return fmt.Errorf("extract target must be a pointer, got `nil`")
	}
	if !reflection.IsPtr(target) {
		return fmt.Errorf("extract target must be a pointer, got `%s`", reflect.TypeOf(target))
	}
	typ := reflect.TypeOf(target)
	param := parameter{
		name:  params.Name,
		res:   typ.Elem(),
		embed: isEmbedParameter(typ),
	}
	value, err := param.ResolveValue(c)
	if err != nil {
		return err
	}
	targetValue := reflect.ValueOf(target).Elem()
	targetValue.Set(value)
	return nil
}

// Invoke calls provided function.
func (c *Container) Invoke(fn interface{}, options ...InvokeOption) error {
	params := InvokeParams{}
	for _, opt := range options {
		opt.apply(&params)
	}
	if !c.compiled {
		return fmt.Errorf("container not compiled")
	}
	invoker, err := newInvoker(fn)
	if err != nil {
		return err
	}
	return invoker.Invoke(c)
}

// Cleanup runs destructors in order that was been created.
func (c *Container) Cleanup() {
	for _, cleanup := range c.cleanups {
		cleanup()
	}
}

// processProviderInterface represents instances as interfaces and groups.
func (c *Container) processProviderInterface(provider internalProvider, as interface{}) {
	// create interface from provider
	iface := newProviderInterface(provider, as)
	key := iface.Key()
	if c.graph.Exists(key) {
		stub := newProviderStub(key, "have several implementations")
		c.graph.Replace(key, stub)
	} else {
		// add interface node
		c.graph.Add(key, iface)
	}
	// create group
	group := newProviderGroup(key)
	groupKey := group.Key()
	// check exists
	if c.graph.Exists(groupKey) {
		// if exists use existing group
		node := c.graph.Get(groupKey)
		group = node.Value.(*providerGroup)
	} else {
		// else add new group to graph
		c.graph.Add(groupKey, group)
	}
	// add provider reference into group
	providerKey := provider.Key()
	group.Add(providerKey)
}

// registerProviderParameters registers provider parameters in a dependency graph.
func (c *Container) registerProviderParameters(p internalProvider) {
	for _, param := range p.ParameterList() {
		provider, exists := param.ResolveProvider(c)
		if exists {
			c.graph.Edge(provider.Key(), p.Key())
			continue
		}
		if !exists && !param.optional {
			panicf("%s: dependency %s not exists in container", p.Key(), param)
		}
	}
}
