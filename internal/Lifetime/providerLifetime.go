package Lifetime

type ProviderLifetime int32

const (
	Singleton ProviderLifetime = 0
	Scoped    ProviderLifetime = 1
	Transient ProviderLifetime = 2
)
