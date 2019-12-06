package di

import "fmt"

// ErrProviderNotFound
type ErrProviderNotFound struct {
	k key
}

func (e ErrProviderNotFound) Error() string {
	return fmt.Sprintf("type `%s` not exists in container", e.k)
}
