package ditest

// Baz
type Baz struct {
	foo *Foo
	bar *Bar
}

// NewBaz
func NewBaz(foo *Foo, bar *Bar) *Baz {
	return &Baz{
		foo: foo,
		bar: bar,
	}
}

func (b *Baz) Foo() *Foo { return b.foo }
func (b *Baz) Bar() *Bar { return b.bar }
