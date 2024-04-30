# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.20.0] - 2024-04-30

### Changed

- Ingore patch suffix from the app operator deployment when checking for versions.

### Removed

- Remove legacy monitoring label.

## [0.19.2] - 2024-01-29

### Fixed

- Move PSS values under the global property

## [0.19.1] - 2023-12-05

### Changed

- Configure gsoci.azurecr.io as the registry to use by default

## [0.19.0] - 2023-11-10

### Changed

- Add a switch for PSP CR installation.

## [0.18.0] - 2023-07-04

### Changed

- Updated default `securityContext` values to comply with PSS policies.

## [0.17.6] - 2023-07-04

### Added

- Add service monitor to be scraped by Prometheus Agent.

### Removed

- Stop pushing to `openstack-app-collection`.

## [0.17.5] - 2023-04-25

### Added

- Add runtime default seccomp profile to app-exporter.
- Add app icon.

### Removed

- Remove push to `shared-app-collection` as it is deprecated.

## [0.17.4] - 2022-11-18

### Fix

- Configure `ServiceMonitor` to honor labels and drop Prometheus label.

## [0.17.3] - 2022-10-17

### Added

- ability to configure the following parts of the app:
  - config:
    debug: true
    listenPort: 8000
    alertDefaultTeam: noteam
    appTeamMappings: ""
    retiredTeamsMapping: ""

## [0.17.2] - 2022-10-11

### Changed

- Increase scrape timeout to 45s and add interval greater than the scrape timeout.

## [0.17.1] - 2022-10-06

### Changed

- Reduce scrape timeout to 30s to prevent Prometheus from crashing.

## [0.17.0] - 2022-10-06

### Changed

- Make scrape timeout configurable and set the default to 45s.

## [0.16.2] - 2022-09-06

### Fixes

- Fix skipping app namespaces if the app-operator version was already scraped

### Changed

- Skip collecting app versions for apps in the org-* namespaces as apps in these (typically CAPI cluster) namespaces the `app-operator.giantswarm.io/version` is not mandatory / makes sense

## [0.16.1] - 2022-07-01

### Fixes

- Region value reference.

## [0.16.0] - 2022-07-01

### Changed

- Increase scrape timeout to 30s in China.

## [0.15.0] - 2022-06-17

### Added

- Add Service monitor.

## [0.14.0] - 2022-05-16

### Fixed

- Normalize all the versions the exporter works with.

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

[Unreleased]: https://github.com/giantswarm/app-exporter/compare/v0.20.0...HEAD
[0.20.0]: https://github.com/giantswarm/app-exporter/compare/v0.19.2...v0.20.0
[0.19.2]: https://github.com/giantswarm/app-exporter/compare/v0.19.1...v0.19.2
[0.19.1]: https://github.com/giantswarm/app-exporter/compare/v0.19.0...v0.19.1
[0.19.0]: https://github.com/giantswarm/app-exporter/compare/v0.18.0...v0.19.0
[0.18.0]: https://github.com/giantswarm/app-exporter/compare/v0.17.6...v0.18.0
[0.17.6]: https://github.com/giantswarm/app-exporter/compare/v0.17.5...v0.17.6
[0.17.5]: https://github.com/giantswarm/app-exporter/compare/v0.17.4...v0.17.5
[0.17.4]: https://github.com/giantswarm/app-exporter/compare/v0.17.3...v0.17.4
[0.17.3]: https://github.com/giantswarm/app-exporter/compare/v0.17.2...v0.17.3
[0.17.2]: https://github.com/giantswarm/app-exporter/compare/v0.17.1...v0.17.2
[0.17.1]: https://github.com/giantswarm/app-exporter/compare/v0.17.0...v0.17.1
[0.17.0]: https://github.com/giantswarm/app-exporter/compare/v0.16.2...v0.17.0
[0.16.2]: https://github.com/giantswarm/app-exporter/compare/v0.16.1...v0.16.2
[0.16.1]: https://github.com/giantswarm/app-exporter/compare/v0.16.0...v0.16.1
[0.16.0]: https://github.com/giantswarm/app-exporter/compare/v0.15.0...v0.16.0
[0.15.0]: https://github.com/giantswarm/app-exporter/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/giantswarm/app-exporter/compare/v0.13.0...v0.14.0
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
