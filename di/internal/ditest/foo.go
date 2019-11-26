package ditest

import (
	"errors"
)

// Foo test struct
type Foo struct{}

// NewFoo create new foo
func NewFoo() *Foo {
	return &Foo{}
}

// NewCycleFooBar
func NewCycleFooBar(bar *Bar) *Foo {
	return &Foo{}
}

// NewFooError
func NewFooError() (*Foo, error) {
	return nil, errors.New("internal error")
}

// NewFooCleanup
func NewFooCleanup() (*Foo, func()) {
	return &Foo{}, func() {}
}

// CreateFooConstructor
func CreateFooConstructor(foo *Foo) func() *Foo {
	return func() *Foo {
		return foo
	}
}

// CreateFooConstructorWithCleanup
func CreateFooConstructorWithCleanup(cleanup func()) func() (*Foo, func()) {
	return func() (foo *Foo, i func()) {
		return &Foo{}, cleanup
	}
}
