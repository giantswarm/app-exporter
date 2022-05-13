# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Normalize all the versions the exporter works with.

### Changed

- Change how exporter fetches the ACE for a given App CR.

## [0.13.0] - 2022-01-28

### Added

- Add `cluster_missing` label to the metric for detecting org-namespace App CRs without `giantswarm.io/cluster` Kubernetes label.

## [0.12.1] - 2022-01-24

### Fixed

- Set release status to `not-installed` when it is empty string.

## [0.12.0] - 2022-01-05

### Added

- Add manual mapping of apps to teams for alert routing when team annotation
in Chart.yaml is missing.

## [0.11.0] - 2021-12-15

## Changed

- Use `key.Version()` to get App CR version. This is to support `v`-prefixed App CR versions.

## [0.10.1] - 2021-12-01

### Changed

- Use `apiextensions-application` instead of `apiextensions` for CRDs to remove
CAPI dependency.

### Fixed

- Add `halo` to `cabbage` to retired teams mapping.

## [0.10.0] - 2021-11-12

### Changed

- Change default team and route retired teams to new teams for app metrics.

### Removed

- Remove helm 2 app-operator ready logic as migration is complete.

## [0.9.0] - 2021-09-08

### Added

- Add `app_version` label to app metrics.

## [0.8.0] - 2021-08-31

### Added

- Add latest_version and upgrade_available labels to show if App CRs
in public catalogs have an upgrade available.

## [0.7.1] - 2021-08-20

### Fixed

- Fix finding AppCatalogEntry CRs and reduce number of CRs fetched.

## [0.7.0] - 2021-08-19

### Fixed

- Find the relevant AppCatalogEntry CRs by searching all namespaces.

## [0.6.1] - 2021-08-06

### Fixed

- Revert VPA as resources are too low on management clusters with high number of apps.

## [0.6.0] - 2021-08-06

### Added

- Use VPA to manage deployment resources.

## [0.5.0] - 2021-06-17

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

[Unreleased]: https://github.com/giantswarm/app-exporter/compare/v0.13.0...HEAD
[0.13.0]: https://github.com/giantswarm/app-exporter/compare/v0.12.1...v0.13.0
[0.12.1]: https://github.com/giantswarm/app-exporter/compare/v0.12.0...v0.12.1
[0.12.0]: https://github.com/giantswarm/app-exporter/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/giantswarm/app-exporter/compare/v0.10.1...v0.11.0
[0.10.1]: https://github.com/giantswarm/app-exporter/compare/v0.10.0...v0.10.1
[0.10.0]: https://github.com/giantswarm/app-exporter/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/giantswarm/app-exporter/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/giantswarm/app-exporter/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/giantswarm/app-exporter/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/giantswarm/app-exporter/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/giantswarm/app-exporter/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/giantswarm/app-exporter/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/app-exporter/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/app-exporter/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/app-exporter/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/app-exporter/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/app-exporter/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/app-exporter/releases/tag/v0.1.0
