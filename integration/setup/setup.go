// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/app-exporter/integration/key"
	"github.com/giantswarm/app-exporter/pkg/appsetup"
	"github.com/giantswarm/app-exporter/pkg/project"
)

func Setup(m *testing.M, config Config) {
	var v int

	var err error

	ctx := context.Background()

	err = installResources(ctx, config)
	if err != nil {
		config.Logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("failed to install %#q", project.Name()), "stack", fmt.Sprintf("%#v", err))
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	os.Exit(v)
}

func installResources(ctx context.Context, config Config) error {
	var err error

	{
		apps := []appsetup.App{
			{
				CatalogName:   key.ControlPlaneTestCatalogName(),
				CatalogURL:    key.ControlPlaneTestCatalogStorageURL(),
				Name:          project.Name(),
				Namespace:     key.Namespace(),
				Version:       project.Version(),
				WaitForDeploy: true,
			},
		}
		err = config.AppSetup.InstallApps(ctx, apps)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
