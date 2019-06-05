package inject

import (
	"github.com/defval/inject/internal/provider"
	"github.com/defval/inject/internal/provider/combined"
	"github.com/defval/inject/internal/provider/ctor"
	"github.com/defval/inject/internal/provider/direct"
	"github.com/defval/inject/internal/provider/object"
)

// createProvider creates provider
func createProvider(po *providerOptions) (_ provider.Provider, err error) {
	switch provider.DetectType(po.provider) {
	case provider.Constructor:
		return ctor.New(po.provider)
	case provider.Combined:
		var options []object.Option

		if po.includeExported {
			options = append(options, object.Exported())
		}

		return combined.New(po.provider, options...)
	case provider.Object:
		var options []object.Option

		if po.includeExported {
			options = append(options, object.Exported())
		}

		return object.New(po.provider, options...)
	default:
		return direct.New(po.provider), nil
	}
}
