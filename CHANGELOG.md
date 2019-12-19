# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## Fixed

- Removed duplicate function `resolveParameterProvider()`

## v2.2.2

Internal refactoring

### Added

- Visualize parameter bag

### Fixed

- Visualize type detection

## v2.2.1

### Fixed

- Invoke: interface is nil, not error

## v2.2.0

### Added

- `container.Invoke()` for invocations

## v2.1.1

### Fixed

- Incorrect di.Parameter resolving

## v2.1.0

### Added

- Visualization

## v2.0.1

### Added

- Helper methods to ParameterBag

## v2.0.0

Massive refactoring and rethinking features.

### Changed

- Graph implementation
- Simplify injection code
- Documentation

### Added

- Prototypes
- Cleanup
- Parameter bag
- Optional parameters
- Low-level container interface

### Removed

- Replacing (investigating)
- Non constructor providers
- Combined providers

## v1.5.2

### Fixed

- Checksum problem

## v1.5.1

### Added

- Ability to extract `github.com/emicklei/*dot.Graph` from container

## v1.5.0

### Added

- Error `inject.ErrTypeNotProvided`

## v1.4.4

### Changed

- Internal refactoring of adding nodes

## v1.4.3

### Fixed

- Improve test coverage

### Changed

- Internal refactoring of groups creation

## v1.4.2

### Fixed

- Replace: check that provider implement interface

## v1.4.1

### Fixed

- Lint

## v1.4.0

### Change

- `Container.WriteTo()` signature to `io.WriterTo`

## v1.3.1

### Added

- Documentation

## v1.3.0

### Added

- Graph visualization

## v1.2.1

### Fixed

- inject.As() allows provide same interface without name

## v1.2.0

### Changed

- Combined provider declaration
- Some refactoring

## v1.1.1

### Changed

- Refactor graph storage

## v1.1.0

### Added

- Combined provider

### Changed

- Provider refactoring

## v1.0.0

- Initial release