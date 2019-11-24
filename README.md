<img width="312"
src="https://github.com/defval/inject/raw/master/logo.png">[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Dependency%20injection%20container%20for%20Golang&url=https://github.com/defval/inject&hashtags=golang,go,di,dependency-injection)

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff)](https://godoc.org/github.com/defval/inject)
![Release](https://img.shields.io/github/tag/defval/inject.svg?label=release&color=24B898&logo=github&style=for-the-badge)
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)
![Contributors](https://img.shields.io/github/contributors/defval/inject.svg?style=for-the-badge)

## How will dependency injection help me?

Dependency injection is one form of the broader technique of inversion
of control. It is used to increase modularity of the program and makes
it extensible.

--------

This container implementation inspired by
[google/wire](https://github.com/google/wire),
[uber-go/fx](https://github.com/uber-go/fx) and
[uber-go/dig](https://github.com/uber-go/dig).


## Installing

```shell
go get -u github.com/defval/inject/v2
```

This library has two using levels. Package `inject` provides a clean and
easy way to build your application components. And package `di` that
provides more advanced techniques.

## Constructors

A constructor must be a function. It must have one or two results. The
first is a result object. Often the type of result object is a pointer.
The second may be an error. Error is optional. A constructor may have an
unlimited number of parameters and a container will provide all of these
automatically.
