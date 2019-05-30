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

Provide

First of all, when creating a new container, you need to describe
how to create each instance of a dependency. To do this, use the container
option inject.Provide().

	container, err := New(
		Provide(NewDependency),
		Provide(NewAnotherDependency)
	)

	func NewDependency(dependency *pkg.AnotherDependency) *pkg.Dependency {
		return &pkg.Dependency{
			dependency: dependency,
		}
	}

	func NewAnotherDependency() (*pkg.AnotherDependency, error) {
		if dependency, err = initAnotherDependency(); err != nil {
			return nil, err
		}

		return dependency, nil
	}

Now, container knows how to create *pkg.Dependency and *pkg.AnotherDependency.
For advanced providing see inject.Provide() and inject.ProvideOption documentation.

Extract

After building a container, it is easy to get any previously provided type.
To do this, use the container's Extract() method.

	var anotherDependency *pkg.AnotherDependency
	if err = container.Extract(&anotherDependency); err != nil {
		// handle error
	}

The container collects a dependencies of *pkg.AnotherDependency, creates its instance and
places it in a target pointer.
For advanced extraction see Extract() and inject.ExtractOption documentation.
*/
package inject // import "github.com/defval/inject"
