// The MIT License (MIT)
//
// Copyright (c) 2019 Maxim Bovtunov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

/*
Package inject make your dependency injection easy. Container allows you to inject dependencies
into constructors or structures without the need to have specified
each argument manually.

Provide dependency

First of all, when creating a new container, you need to describe
how to create each instance of a dependency.To do this, use the container
option inject.Provide(). The first argument in this function is a `provider`.
It determines how to create dependency.

Provider can be a constructor function with optional error:

  	// dependency constructor function
  	func NewDependency(dependency *pkg.AnotherDependency) *pkg.Dependency {
		return &pkg.Dependency{
			dependency: dependency,
		}
	}

  	// and with possible initialization error
  	func NewAnotherDependency() (*pkg.AnotherDependency, error) {
		if dependency, err = initAnotherDependency(); err != nil {
			return nil, err
		}
		return dependency, nil
 	}

	// container initialization codeы
	container, err := New(
		Provide(NewDependency),
		Provide(NewAnotherDependency)
	)

In this case, the container knows how to create `*pkg.AnotherDependency`
and can handle an instance creation error.

Also, a provider can be a structure pointer with public fields:

	// package pkg
	type Dependency struct {
		AnotherDependency *pkg.AnotherDependency `inject:""`
	}

	// container initialization code
	container, err := New(
		// pointer to structure
		Provide(&pkg.Dependency{}),
		// or structure value
		Provide(pkg.Dependency{})
	)

In this case, the necessity of implementing specific fields are defined
with the tag `inject`.
*/
package inject // import "github.com/defval/inject"
