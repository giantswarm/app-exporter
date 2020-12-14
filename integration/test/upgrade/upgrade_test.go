// +build k8srequired

package upgrade

import (
	"context"
	"testing"

	"github.com/giantswarm/appcatalog"
	"github.com/giantswarm/apptest"

	"github.com/giantswarm/app-exporter/integration/env"
	"github.com/giantswarm/app-exporter/integration/key"
	"github.com/giantswarm/app-exporter/integration/setup"
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

	var latestVersion string
	{
		latestVersion, err = appcatalog.GetLatestVersion(ctx, key.ControlPlaneCatalogURL(), project.Name(), "")
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}
	}

	app := apptest.App{
		CatalogName:   key.ControlPlaneCatalogName(),
		Name:          project.Name(),
		Namespace:     key.Namespace(),
		Version:       latestVersion,
		WaitForDeploy: true,
	}

	// 1. Install the latest version
	{
		apps := []apptest.App{app}
		err = config.AppTest.InstallApps(ctx, apps)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}
	}

	app.CatalogName = key.ControlPlaneTestCatalogName()
	app.Version = ""
	app.SHA = env.CircleSHA()

	// 2. Upgrade app CR to this testing version
	{
		err = config.AppTest.UpgradeApp(ctx, app)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}
	}
}
