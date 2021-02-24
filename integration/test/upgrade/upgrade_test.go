// +build k8srequired

package upgrade

import (
	"context"
	"testing"

	"github.com/giantswarm/apptest"

	"github.com/giantswarm/app-exporter/integration/env"
	"github.com/giantswarm/app-exporter/integration/key"
	"github.com/giantswarm/app-exporter/integration/setup"
	"github.com/giantswarm/app-exporter/integration/templates"
	"github.com/giantswarm/app-exporter/pkg/project"
)

var (
	config setup.Config
)

func init() {
	var err error

	{
		config, err = setup.NewConfig()
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestUpgrade(t *testing.T) {
	var err error
	ctx := context.Background()

	currentApp := apptest.App{
		CatalogName:   key.ControlPlaneCatalogName(),
		Name:          project.Name(),
		Namespace:     key.Namespace(),
		ValuesYAML:    templates.AppExporterValues,
		WaitForDeploy: true,
	}

	desiredApp := apptest.App{
		CatalogName:   key.ControlPlaneTestCatalogName(),
		Name:          project.Name(),
		Namespace:     key.Namespace(),
		SHA:           env.CircleSHA(),
		ValuesYAML:    templates.AppExporterValues,
		WaitForDeploy: true,
	}

	{
		err = config.AppTest.UpgradeApp(ctx, currentApp, desiredApp)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}
	}
}
