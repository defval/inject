# Inject ![Release](https://img.shields.io/github/tag/defval/inject.svg?label=release&logo=github&style=flat-square) [![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square&logo=go&logoColor=ffffff)](https://godoc.org/github.com/defval/inject) [![Build Status](https://img.shields.io/travis/defval/inject.svg?style=flat-square&logo=travis)](https://travis-ci.org/defval/inject) [![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=flat-square&logo=codecov)](https://codecov.io/gh/defval/inject)


Dependency injection container allows you to inject dependencies
into constructors or structures without the need to have specified
each argument manually.

This container implementation inspired by [google/wire](https://github.com/google/wire),
[uber-go/fx](https://github.com/uber-go/fx) and [uber-go/dig](https://github.com/uber-go/dig).

## Installing

```shell
go get -u github.com/defval/inject
```

## Features

- Construct arguments injection
- Tagged and public struct fields injection
- Set interface for implementations
- Named definitions
- Replacing
