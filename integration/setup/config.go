// +build k8srequired

package setup

import (
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/app-exporter/integration/env"
	"github.com/giantswarm/app-exporter/pkg/appsetup"
)

type Config struct {
	AppSetup   appsetup.Interface
	K8sClients k8sclient.Interface
	Logger     micrologger.Logger
}

func NewConfig() (Config, error) {
	var err error

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var k8sClients *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			Logger: logger,

			KubeConfigPath: env.KubeConfigPath(),
		}

		k8sClients, err = k8sclient.NewClients(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var appSetup appsetup.Interface
	{
		c := appsetup.Config{
			K8sClient: k8sClients,
			Logger:    logger,
		}

		appSetup, err = appsetup.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		AppSetup:   appSetup,
		K8sClients: k8sClients,
		Logger:     logger,
	}

	return c, nil
}
