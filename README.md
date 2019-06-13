# Inject
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge&logo=go&logoColor=ffffff)](https://godoc.org/github.com/defval/inject)
![Release](https://img.shields.io/github/tag/defval/inject.svg?label=release&logo=github&style=for-the-badge)
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)


Dependency injection container allows you to inject dependencies
into constructors or structures without the need to have specified
each argument manually.

This container implementation inspired by [google/wire](https://github.com/google/wire),
[uber-go/fx](https://github.com/uber-go/fx) and [uber-go/dig](https://github.com/uber-go/dig).

See [godoc](https://godoc.org/github.com/defval/inject) for feel the difference.

## Installing

```shell
go get -u github.com/defval/inject
```

## Quickstart

### Describe dependencies

##### Controller interface

```go
// Controller interface.
type Controller interface {
	Register(mux *http.ServeMux)
}
```

##### Controller

`AccountController` implements `Controller` interface.

```go
// AccountController contains account related http methods.
type AccountController struct {
}

// Register add routes to mux.
func (c *AccountController) Register(mux *http.ServeMux) {
	mux.HandleFunc("/account", c.Index)
}

func (c *AccountController) Index(writer http.ResponseWriter, request *http.Request) {
	_, _ = io.WriteString(writer, "account")
}
```

##### HTTP Mux

Define `*http.ServeMux` constructor and register all controllers routes.
For this we use a group of interfaces.

```go
// NewMux creates a new http mux.
func NewMux(controllers []Controller) *http.ServeMux {
	mux := &http.ServeMux{}

	for _, ctrl := range controllers {
		ctrl.Register(mux)
	}

	return mux
}
```

##### Server

The server needs a `http.Handler` to work.

```go
// NewServer creates a new http server.
func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
	}
}
```

### Build container

```go
container, err := inject.New(
    inject.Provide(NewMux, inject.As(new(http.Handler))), // provide *http.ServeMux as http.Handler interface
    inject.Provide(NewServer),
    inject.Provide(&AccountController{}, inject.As(new(Controller))), // inject.As automatically creates a group []Controller 
)

if err != nil {
    // handle err
}
```

### Extract type

```go
// extract type
var server *http.Server
if err = container.Extract(&server); err != nil {
    // handle err
}

server.ListenAndServe()
```

## Visualize dependency graph

Container supports `fmt.Stringer` interface. The string is a graph
description via [graphviz dot language](https://www.graphviz.org/).

### Raw output

```
digraph  {
	n4[color="limegreen",fontcolor="white",fontname="Helvetica",label="*http.ServeMux",shape="box",style="filled"];
	n2[color="limegreen",fontcolor="white",fontname="Helvetica",label="*http.Server",shape="box",style="filled"];
	n7[color="limegreen",fontcolor="white",fontname="Helvetica",label="*inject_test.AccountController[account]",shape="box",style="filled"];
	n9[color="limegreen",fontcolor="white",fontname="Helvetica",label="*inject_test.AuthController[auth]",shape="box",style="filled"];
	n1[color="limegreen",fontcolor="white",fontname="Helvetica",label="*log.Logger",shape="box",style="filled"];
	n6[color="orange",fontcolor="white",fontname="Helvetica",label="[]http.Handler",shape="doubleoctagon",style="filled"];
	n5[color="orange",fontcolor="white",fontname="Helvetica",label="[]inject_test.Controller",shape="doubleoctagon",style="filled"];
	n3[color="royalblue",fontcolor="white",fontname="Helvetica",label="http.Handler",style="filled"];
	n8[color="royalblue",fontcolor="white",fontname="Helvetica",label="inject_test.Controller[account]",style="filled"];
	n10[color="royalblue",fontcolor="white",fontname="Helvetica",label="inject_test.Controller[auth]",style="filled"];
	n4->n3;
	n4->n6;
	n7->n8;
	n7->n5;
	n9->n5;
	n9->n10;
	n1->n2;
	n1->n4;
	n1->n7;
	n1->n9;
	n5->n4;
	n3->n2;
}
```

### Image

<img src="https://github.com/defval/inject/raw/12-visualize-graph/graph.png">