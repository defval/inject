package ditest

import "github.com/defval/inject/v2/di"

// NewQuxProvider
func NewQuxProvider() *QuxProvider {
	return &QuxProvider{}
}

// QuxProvider
type QuxProvider struct {
	Fooer Fooer `di:""`
}

// Identity
func (p *QuxProvider) Identity() di.Identity {
	return di.IdentityOf(&Qux{})
}

// Provide
func (p *QuxProvider) Provide() (interface{}, error) {
	return &Qux{fooer: p.Fooer}, nil
}
