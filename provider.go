package inject

import (
	"github.com/defval/inject/internal/graph"
	"github.com/defval/inject/internal/provider"
)

// Provider helper struct that indicates that structure is injection provider.
type Provider struct {
}

// IsInjectProvider ...
func (p *Provider) IsInjectProvider() {
	panic("hopefully, never be called")
}

// determineInstanceProvider creates instanceProvider
func determineInstanceProvider(po *providerOptions) (_ graph.InstanceProvider, err error) {
	if provider.IsConstructor(po.provider) {
		return provider.NewConstructorProvider(po.provider)
	}

	if provider.IsCombinedProvider(po.provider) {
		return provider.NewCombinedProvider(po.provider, "inject", po.includeExported)
	}

	if provider.IsObjectProvider(po.provider) {
		return provider.NewObjectProvider(po.provider, "inject", po.includeExported)
	}

	return provider.NewDirectProvider(po.provider), nil
}
