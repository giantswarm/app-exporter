module github.com/giantswarm/app-exporter

go 1.15

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/giantswarm/apiextensions/v3 v3.18.1
	github.com/giantswarm/app/v4 v4.4.0
	github.com/giantswarm/apptest v0.10.2
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/exporterkit v0.2.1
	github.com/giantswarm/k8sclient/v5 v5.10.0
	github.com/giantswarm/k8sportforward/v2 v2.0.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/operatorkit/v2 v2.0.2
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/viper v1.7.1
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v0.18.9
)

replace (
	// Use v0.8.10 of hcsshim to fix nancy alert.
	github.com/Microsoft/hcsshim v0.8.7 => github.com/Microsoft/hcsshim v0.8.10
	// Apply fix for CVE-2020-15114 not yet released in github.com/spf13/viper.
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.3
	// Use v1.0.0-rc7 of runc to fix nancy alert.
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc7
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
