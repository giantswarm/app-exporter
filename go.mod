module github.com/giantswarm/app-exporter

go 1.16

require (
	github.com/Masterminds/semver/v3 v3.2.1
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
	github.com/google/go-cmp v0.6.0
	github.com/prometheus/client_golang v1.17.0
	github.com/spf13/viper v1.14.0
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	sigs.k8s.io/controller-runtime v0.9.7
	sigs.k8s.io/yaml v1.4.0
)

replace (
	github.com/Microsoft/hcsshim v0.8.7 => github.com/Microsoft/hcsshim v0.8.10
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.5
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.9.1
	// Use go-logr/logr v0.1.0 due to breaking changes in v0.2.0 that can't be applied.
	github.com/go-logr/logr v0.2.0 => github.com/go-logr/logr v0.1.0
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	github.com/kataras/iris/v12 => github.com/kataras/iris/v12 v12.1.8
	github.com/labstack/echo/v4 => github.com/labstack/echo/v4 v4.11.2
	github.com/microcosm-cc/bluemonday => github.com/microcosm-cc/bluemonday v1.0.21
	github.com/nats-io/jwt/v2 => github.com/nats-io/jwt/v2 v2.5.2
	github.com/nats-io/nats-server/v2 => github.com/nats-io/nats-server/v2 v2.9.10
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc7
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.5
	github.com/urfave/negroni/v3 => github.com/urfave/negroni/v3 v3.0.0
	github.com/valyala/fasthttp => github.com/valyala/fasthttp v1.43.0
	// Same as go-logr/logr, klog/v2 is using logr v0.2.0
	k8s.io/klog/v2 v2.4.0 => k8s.io/klog/v2 v2.0.0
)
