module github.com/giantswarm/app-exporter

go 1.15

require (
	github.com/Masterminds/semver/v3 v3.1.0
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/giantswarm/apiextensions/v3 v3.4.0
	github.com/giantswarm/apptest v0.4.0
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/exporterkit v0.2.0
	github.com/giantswarm/k8sclient/v5 v5.0.0
	github.com/giantswarm/k8sportforward/v2 v2.0.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/microkit v0.2.1
	github.com/giantswarm/micrologger v0.3.4
	github.com/giantswarm/operatorkit/v2 v2.0.0
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/afero v1.3.4 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v0.18.9
)

replace (
	// Apply fix for CVE-2020-15114 not yet released in github.com/spf13/viper.
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.3
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
