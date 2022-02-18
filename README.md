[![CircleCI](https://circleci.com/gh/giantswarm/app-exporter.svg?&style=shield)](https://circleci.com/gh/giantswarm/app-exporter)

# app-exporter

A Prometheus exporter for [App] CRs (custom resources). Part of the Giant Swarm [App Platform].

## Related Components

- [app-admission-controller] provides mutating and validating webhooks for App CRs.
- [app-operator] manages App CRs and generates Chart CRs.
- [chart-operator] manages Chart CRs and automates Helm deployments.

## Getting Project

Clone the git repository: https://github.com/giantswarm/app-exporter.git

### How to build

Build it using the standard `go build` command.

```
go build
```

## Changelog

See [CHANGELOG](CHANGELOG.md)

## Contact

- Bugs: [issues](https://github.com/giantswarm/app-exporter/issues)
- Please visit https://www.giantswarm.io/responsible-disclosure for information on reporting security issues.

## License

app-exporter is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for
details.

[App]: https://docs.giantswarm.io/ui-api/management-api/crd/apps.application.giantswarm.io/
[App Platform]: https://docs.giantswarm.io/app-platform/overview/
[app-admission-controller]: https://github.com/giantswarm/app-admission-controller
[app-operator]: https://github.com/giantswarm/app-operator
[chart-operator]: https://github.com/giantswarm/chart-operator

