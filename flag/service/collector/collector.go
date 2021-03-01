package collector

import (
	"github.com/giantswarm/app-exporter/flag/service/collector/apps"
	"github.com/giantswarm/app-exporter/flag/service/collector/provider"
)

type Collector struct {
	Apps     apps.Apps
	Provider provider.Provider
}
