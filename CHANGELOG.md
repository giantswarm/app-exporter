# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add `deployed_version`, `version_mismatch` label to app metrics.

### Changed

- Prepare helm values to configuration management.

## [0.4.0] - 2021-03-05

### Added

- Get team from `application.giantswarm.io/team` annotation on the App CR if it
exists. Otherwise check the AppCatalogEntry CR.

## [0.3.0] - 2021-03-02

### Added

- Add app label to app metrics with the name of the app.
- Extend app-operator ready metric to include per workload cluster instances.
- Get team from `application.giantswarm.io/team` or `application.giantswarm.io/owners`
annotations in Chart.yaml. 

### Changed

- Update apiextensions to v3 and replace CAPI with Giant Swarm fork.

### Removed

- App to team mapping configmap.

## [0.2.1] - 2020-10-01
### Fixed

- Update deployment annotation to use checksum instead of helm revision to
reduce how often pods are rolled.

## [0.2.0] - 2020-09-04

### Changed

- Decrease the memory size.
- Add network policy.

## [0.1.0] - 2020-08-25

### Added

- Added initial structures.

[Unreleased]: https://github.com/giantswarm/app-exporter/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/giantswarm/app-exporter/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/app-exporter/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/app-exporter/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/app-exporter/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/app-exporter/releases/tag/v0.1.0
