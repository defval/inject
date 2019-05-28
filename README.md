# Inject
[![Build Status](https://img.shields.io/travis/defval/inject.svg?style=for-the-badge&logo=travis)](https://travis-ci.org/defval/inject)
[![Code Coverage](https://img.shields.io/codecov/c/github/defval/inject.svg?style=for-the-badge&logo=codecov)](https://codecov.io/gh/defval/inject)


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

- inject constructor arguments
- inject tagged struct fields
- inject public struct fields
- inject as interface
- inject interface groups
- inject default value of interface group
- inject named definition into structures
- replace interface implementation
- replace provided type
- isolated namespaces

## WIP

- documentation
- inject named definition into constructor
- rework namespaces
