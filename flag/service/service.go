package service

import (
	"github.com/giantswarm/app-exporter/flag/service/collector"
)

// TLS is a data structure for Kubernetes TLS configuration with command line
// flags.
type TLS struct {
	CAFile  string
	CrtFile string
	KeyFile string
}

// Watch is a data structure to hold Kubernetes specific configuration
// for watching for Kubernetes resources.
type Watch struct {
	Namespace string
}

// Kubernetes is a data structure to hold Kubernetes specific command line
// configuration flags.
type Kubernetes struct {
	Address        string
	InCluster      string
	KubeConfig     string
	KubeConfigPath string
	TLS            TLS
	Watch          Watch
}

// Service is an intermediate data structure for command line configuration flags.
type Service struct {
	Collector  collector.Collector
	Kubernetes Kubernetes
}
