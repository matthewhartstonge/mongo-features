# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## [v0.4.0] - 2020-09-02
### Changed
- CHANGELOG.md: adds dates to release versions.
- features: makes `canSession()` checking smarter.

### Fixed
- deps: reverts to `semver@1.5.0` for `dep` and `go mod` compatibility.

## [v0.3.0] - 2020-09-01
### Added
- CHANGELOG.md: adds a changelog.

### Changed
- features: makes features API more explicit.
    - `store.Sessions` => `store.HasSessions`
    - `store.Transactions` => `store.HasTransactions`
    - `store.Version` => `store.MongoVersion`
- features: documents `Features` struct.
- README.md: adds badges.

### Fixes
- _examples: fixes go mod name.

## [v0.2.0] - 2020-09-01
### Fixes
- deps: fixes go mod name


## [v0.1.0] - 2020-09-01

Mongo (.feat Features)!

[Unreleased]: https://github.com/matthewhartstonge/mongo-features/compare/v0.4.0...HEAD
[v0.4.0]: https://github.com/matthewhartstonge/mongo-features/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/matthewhartstonge/mongo-features/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com//matthewhartstonge/mongo-features/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com//matthewhartstonge/mongo-features/releases/tag/v0.1.0
