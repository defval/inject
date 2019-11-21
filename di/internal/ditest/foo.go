package ditest

import (
	"errors"
	"fmt"
)

// Foo test struct
type Foo struct{}

// NewFoo create new foo
func NewFoo() *Foo {
	fmt.Println("asd")
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

// CreateFooConstructor
func CreateFooConstructor(foo *Foo) func() *Foo {
	return func() *Foo {
		return foo
	}
}
