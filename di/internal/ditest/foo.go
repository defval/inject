package ditest

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

// CreateFooConstructor
func CreateFooConstructor(foo *Foo) func() *Foo {
	return func() *Foo {
		return foo
	}
}

// CreateFooConstructorWithError
func CreateFooConstructorWithError(err error) func() (*Foo, error) {
	return func() (foo *Foo, e error) {
		return &Foo{}, err
	}
}

// CreateFooConstructorWithCleanup
func CreateFooConstructorWithCleanup(cleanup func()) func() (*Foo, func()) {
	return func() (foo *Foo, i func()) {
		return &Foo{}, cleanup
	}
}

// CreateFooConstructorWithCleanupAndError
func CreateFooConstructorWithCleanupAndError(cleanup func(), err error) func() (*Foo, func(), error) {
	return func() (foo *Foo, i func(), e error) {
		return &Foo{}, cleanup, err
	}
}
