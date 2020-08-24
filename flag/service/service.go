package service

import (
	"github.com/giantswarm/operatorkit/v2/pkg/flag/service/kubernetes"

	"github.com/giantswarm/app-exporter/flag/service/collector"
)

// Service is an intermediate data structure for command line configuration flags.
type Service struct {
	Collector  collector.Collector
	Kubernetes kubernetes.Kubernetes
}
