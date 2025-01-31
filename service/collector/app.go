package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/app/v6/pkg/key"
	"github.com/giantswarm/k8sclient/v6/pkg/k8sclient"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	expkey "github.com/giantswarm/app-exporter/internal/key"
)

var (
	appDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "app", "info"),
		"Managed apps status.",
		[]string{
			labelApp,
			labelAppVersion,
			labelCatalog,
			labelClusterMissing,
			labelDeployedVersion,
			labelLatestVersion,
			labelName,
			labelNamespace,
			labelStatus,
			labelTeam,
			labelUpgradeAvailable,
			labelVersion,
			labelVersionMismatch,
		},
		nil,
	)

	appCordonExpireTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "app", "cordon_expire_time_seconds"),
		"A metric of the expire time of cordoned apps unix seconds.",
		[]string{
			labelName,
			labelNamespace,
		},
		nil,
	)
)

// AppConfig is this collector's configuration struct.
type AppConfig struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	AppTeamMappings     map[string]string
	DefaultTeam         string
	Provider            string
	RetiredTeamsMapping map[string]string
}

// App is the main struct for this collector.
type App struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	appTeamMappings     map[string]string
	defaultTeam         string
	provider            string
	retiredTeamsMapping map[string]string
}

// NewApp creates a new App metrics collector
func NewApp(config AppConfig) (*App, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.DefaultTeam == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.DefaultTeam must not be empty", config)
	}
	if config.Provider == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Provider must not be empty", config)
	}
	if config.RetiredTeamsMapping == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.RetiredTeamsMapping must not be empty", config)
	}

	a := &App{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		appTeamMappings:     config.AppTeamMappings,
		defaultTeam:         config.DefaultTeam,
		provider:            config.Provider,
		retiredTeamsMapping: config.RetiredTeamsMapping,
	}

	return a, nil
}

// Collect is the main metrics collection function.
func (a *App) Collect(ch chan<- prometheus.Metric) error {
	ctx := context.Background()

	err := a.collectAppStatus(ctx, ch)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// Describe emits the description for the metrics collected here.
func (a *App) Describe(ch chan<- *prometheus.Desc) error {
	ch <- appDesc
	ch <- appCordonExpireTimeDesc
	return nil
}

func (a *App) collectAppStatus(ctx context.Context, ch chan<- prometheus.Metric) error {
	apps := &v1alpha1.AppList{}
	err := a.k8sClient.CtrlClient().List(ctx, apps)
	if err != nil {
		return microerror.Mask(err)
	}

	latestAppVersions, err := a.getLatestAppVersions(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	teamMappings, err := a.getTeamMappings(ctx, apps.Items)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, app := range apps.Items {
		appCatalogEntryName := key.AppCatalogEntryName(key.CatalogName(app), key.AppName(app), key.Version(app))
		team := teamMappings[appCatalogEntryName]
		if team == "" {
			// Set the default team if there is no mapping.
			team = a.defaultTeam
		}

		// Team annotation on the App CR overrides if it exists.
		if key.AppTeam(app) != "" {
			team = formatTeamName(key.AppTeam(app))
		}

		if v, ok := app.Labels[annotation.AppTeam]; ok {
			team = formatTeamName(v)
		}

		// Trim `v` prefix from App CR version if there is any
		// TODO once Flux supports more sophisticated regexes
		// we can get rid of this trimming.
		appSpecVersion := key.Version(app)
		appStatusVersion := expkey.FormatVersion(app.Status.Version)

		// clusterMissing is true if `giantswarm.io/cluster` label is missing
		// on the org-namespaced app. Otherwise it's false.
		clusterMissing := key.IsInOrgNamespace(app) && key.ClusterLabel(app) == ""

		// For optional apps in public catalogs we check if an upgrade
		// is available.
		latestVersion := latestAppVersions[fmt.Sprintf("%s-%s", key.CatalogName(app), key.AppName(app))]
		upgradeAvailable := latestVersion != "" && latestVersion != appSpecVersion

		releaseStatus := app.Status.Release.Status
		if releaseStatus == "" {
			releaseStatus = notInstalledStatus
		}

		ch <- prometheus.MustNewConstMetric(
			appDesc,
			prometheus.GaugeValue,
			gaugeValue,
			app.Spec.Name,
			appVersion(app),
			app.Spec.Catalog,
			strconv.FormatBool(clusterMissing),
			appStatusVersion,
			latestVersion,
			app.Name,
			app.Namespace,
			releaseStatus,
			team,
			strconv.FormatBool(upgradeAvailable),
			// Getting version from spec, not status since the version in the spec is the desired version.
			appSpecVersion,
			strconv.FormatBool(appSpecVersion != appStatusVersion),
		)

		if !key.IsAppCordoned(app) {
			continue
		}

		t, err := convertToTime(key.CordonUntil(app))
		if err != nil {
			a.logger.Errorf(ctx, err, "could not convert cordon-until for app %q", key.AppName(app))
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			appCordonExpireTimeDesc,
			prometheus.GaugeValue,
			float64(t.Unix()),
			key.AppName(app),
			key.Namespace(app),
		)
	}
	return nil
}

// getLatestAppVersions checks for the latest version of each app in public catalogs.
// There will be an AppCatalogEntry CR with the label latest=true for the latest
// entry according to semantic versioning.
func (a *App) getLatestAppVersions(ctx context.Context) (map[string]string, error) {
	latestAppVersions := map[string]string{}

	// TODO: Remove community once helm-stable catalog is removed.
	// https://github.com/giantswarm/giantswarm/issues/17490
	catalogLabels, err := labels.Parse(fmt.Sprintf("%s=public,%s!=community", label.CatalogVisibility, label.CatalogType))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	catalogs := &v1alpha1.CatalogList{}
	err = a.k8sClient.CtrlClient().List(ctx, catalogs, &client.ListOptions{LabelSelector: catalogLabels})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	for _, catalog := range catalogs.Items {
		aces := &v1alpha1.AppCatalogEntryList{}
		err = a.k8sClient.CtrlClient().List(ctx, aces, client.InNamespace(catalog.Namespace), client.MatchingLabels{
			label.CatalogName: catalog.Name,
			"latest":          "true",
		})
		if err != nil {
			return nil, microerror.Mask(err)
		}

		for _, ace := range aces.Items {
			latestAppVersions[fmt.Sprintf("%s-%s", ace.Spec.Catalog.Name, ace.Spec.AppName)] = expkey.FormatVersion(ace.Spec.Version)
		}
	}

	return latestAppVersions, nil
}

func (a *App) getOwningTeam(ctx context.Context, app v1alpha1.App, owners []owner) (string, error) {
	for _, o := range owners {
		if key.CatalogName(app) == o.Catalog && a.provider == o.Provider {
			return o.Team, nil
		} else if key.CatalogName(app) == o.Catalog && o.Provider == "" {
			return o.Team, nil
		} else if o.Catalog == "" && a.provider == o.Provider {
			return o.Team, nil
		}
	}

	// If no owning team is found we fall back to the team annotation.
	return "", nil
}

// getTeam returns the team to assign for this app CR. It checks the
// AppCatalogEntry CR to see if it has owners or team annotations.
func (a *App) getTeam(ctx context.Context, app v1alpha1.App) (string, error) {
	var team string
	var err error

	// Team has been configured manually via the configmap. This can be used
	// if the team annotation is missing in Chart.yaml. Make sure the
	// annotation is added and once its present for all deployments of the
	// app the mapping can be removed.
	team = a.appTeamMappings[key.AppName(app)]
	if team != "" {
		return team, nil
	}

	// Note, for custom catalogs carrying the `v`-prefixed app versions, constructing
	// the name like this will give incorrect result without the prefix, resulting in
	// `IsNotFound` error. On the `IsNotFound` error we could fall back to checking
	// next name constructed with just the `app.Spec.Version` (without stripping), however
	// it does not seem necessary at the moment as we do not control the `owners` nor `team`
	// annotations on apps outside our catalogs, hence trying to get it here at all
	// cost might be pointless after all.
	appCatalogEntryName := key.AppCatalogEntryName(key.CatalogName(app), key.AppName(app), key.Version(app))

	ace := &v1alpha1.AppCatalogEntry{}
	{
		// Check giantswarm namespace first as it has more CRs.
		namespaces := []string{"giantswarm", metav1.NamespaceDefault}
		for _, ns := range namespaces {
			err = a.k8sClient.CtrlClient().Get(ctx, types.NamespacedName{Namespace: ns, Name: appCatalogEntryName}, ace)
			if apierrors.IsNotFound(err) {
				// Check next namespace.
				continue
			} else if err != nil {
				return "", microerror.Mask(err)
			} else {
				break
			}
		}
	}

	// Owners annotation takes precedence if it exists.
	ownersYAML := key.AppCatalogEntryOwners(*ace)

	if ownersYAML != "" {
		owners := []owner{}

		err = yaml.Unmarshal([]byte(ownersYAML), &owners)
		if err != nil {
			// If the YAML in the owners annotation is invalid log the error and
			// fall back to trying the team annotation.
			a.logger.Errorf(ctx, err, "could not parse owners YAML for app %q", key.AppName(app))
		}

		if len(owners) > 0 {
			team, err = a.getOwningTeam(ctx, app, owners)
			if err != nil {
				return "", microerror.Mask(err)
			}
		}
	}

	// Use team annotation if there are no owners.
	if team == "" {
		team = key.AppCatalogEntryTeam(*ace)
	}

	// Normalize team name e.g. remove team- prefix if its present.
	team = formatTeamName(team)

	// Team may have been retired. If so we try to map it to the new team.
	newTeam := a.retiredTeamsMapping[team]
	if newTeam != "" {
		return newTeam, nil
	}

	return team, nil
}

// getTeamMappings returns a map of AppCatalogEntry CR names to teams. This
// reduces the number of API calls we need to make to fetch the teams metadata.
func (a *App) getTeamMappings(ctx context.Context, apps []v1alpha1.App) (map[string]string, error) {
	teamMappings := map[string]string{}

	for _, app := range apps {
		appCatalogEntryName := key.AppCatalogEntryName(key.CatalogName(app), key.AppName(app), key.Version(app))

		_, ok := teamMappings[appCatalogEntryName]
		if !ok {
			team, err := a.getTeam(ctx, app)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			teamMappings[appCatalogEntryName] = team
		}
	}

	return teamMappings, nil
}

// appVersion returns the AppVersion if it differs from the Version. This is so
// we can show the upstream chart version packaged by the app.
func appVersion(app v1alpha1.App) string {
	statusVersion := expkey.FormatVersion(app.Status.Version)
	statusAppVersion := expkey.FormatVersion(app.Status.AppVersion)
	if statusAppVersion != statusVersion {
		return statusAppVersion
	}

	return ""
}

func convertToTime(input string) (time.Time, error) {
	layout := "2006-01-02T15:04:05"

	split := strings.Split(input, ".")
	if len(split) == 0 {
		return time.Time{}, microerror.Maskf(invalidExecutionError, "%#q must have at least one item in order to collect metrics for the cordon expiration", input)
	}

	t, err := time.Parse(layout, split[0])
	if err != nil {
		return time.Time{}, microerror.Maskf(invalidExecutionError, "parsing timestamp %#q failed: %#v", split[0], err.Error())
	}

	return t, nil
}

// formatTeamName allows for normalizing the team name. This is needed as our
// GitHub team names use the prefix team but in Prometheus this isn't present.
func formatTeamName(input string) string {
	return strings.TrimPrefix(input, "team-")
}
