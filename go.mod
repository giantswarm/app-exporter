module github.com/giantswarm/app-collector

go 1.14

require (
	github.com/Masterminds/semver/v3 v3.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/giantswarm/apiextensions/v2 v2.1.0
	github.com/giantswarm/exporterkit v0.2.0
	github.com/giantswarm/k8sclient/v4 v4.0.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/microkit v0.2.1
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit/v2 v2.0.0
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/viper v1.6.2
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
)
