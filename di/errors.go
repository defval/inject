package di

import "fmt"

// ErrProviderNotFound
type ErrProviderNotFound struct {
	k key
}

func (e ErrProviderNotFound) Error() string {
	return fmt.Sprintf("not exists in container")
}
