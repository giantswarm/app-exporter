module github.com/giantswarm/app-exporter

go 1.16

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/CloudyKit/jet/v3 v3.0.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/Shopify/goreferrer v0.0.0-20210630161223-536fa16abd6f // indirect
	github.com/giantswarm/apiextensions-application v0.3.0
	github.com/giantswarm/app/v6 v6.6.1
	github.com/giantswarm/backoff v1.0.0
	github.com/giantswarm/exporterkit v1.0.0
	github.com/giantswarm/k8sclient/v6 v6.1.0
	github.com/giantswarm/k8smetadata v0.9.2
	github.com/giantswarm/k8sportforward/v2 v2.0.0
	github.com/giantswarm/microendpoint v1.0.0
	github.com/giantswarm/microerror v0.4.0
	github.com/giantswarm/microkit v1.0.0
	github.com/giantswarm/micrologger v0.6.0
	github.com/giantswarm/operatorkit/v6 v6.1.0
	github.com/gin-gonic/gin v1.8.1 // indirect
	github.com/go-playground/validator/v10 v10.11.0 // indirect
	github.com/google/go-cmp v0.5.7
	github.com/hashicorp/consul/sdk v0.9.0 // indirect
	github.com/iris-contrib/jade v1.1.4 // indirect
	github.com/iris-contrib/schema v0.0.6 // indirect
	github.com/kataras/golog v0.1.7 // indirect
	github.com/kataras/iris/v12 v12.1.8 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/labstack/echo/v4 v4.7.2 // indirect
	github.com/microcosm-cc/bluemonday v1.0.18 // indirect
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-server/v2 v2.8.4 // indirect
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/pkg/sftp v1.13.5 // indirect
	github.com/prometheus/client_golang v1.12.1
	github.com/ryanuber/columnize v2.1.2+incompatible // indirect
	github.com/spf13/viper v1.10.1
	github.com/urfave/negroni v1.0.0 // indirect
	github.com/valyala/fasthttp v1.37.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220614195744-fb05da6f9022 // indirect
	golang.org/x/sys v0.0.0-20220614162138-6c1b26c55098 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	sigs.k8s.io/controller-runtime v0.9.7
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/Microsoft/hcsshim v0.8.7 => github.com/Microsoft/hcsshim v0.8.10
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.3
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	// Use go-logr/logr v0.1.0 due to breaking changes in v0.2.0 that can't be applied.
	github.com/go-logr/logr v0.2.0 => github.com/go-logr/logr v0.1.0
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc7
	// Same as go-logr/logr, klog/v2 is using logr v0.2.0
	k8s.io/klog/v2 v2.4.0 => k8s.io/klog/v2 v2.0.0
)
