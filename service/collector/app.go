package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/app/v5/pkg/key"
	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	appDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "app", "info"),
		"Managed apps status.",
		[]string{
			labelName,
			labelNamespace,
			labelDeployedVersion,
			labelStatus,
			labelTeam,
			labelVersion,
			labelVersionMismatch,
			labelCatalog,
			labelApp,
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

	DefaultTeam string
	Provider    string
}

// App is the main struct for this collector.
type App struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	defaultTeam string
	provider    string
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

	a := &App{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		defaultTeam: config.DefaultTeam,
		provider:    config.Provider,
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
	r, err := a.k8sClient.G8sClient().ApplicationV1alpha1().Apps("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	for _, app := range r.Items {
		team, err := a.getTeam(ctx, app)
		if err != nil {
			a.logger.Errorf(ctx, err, "could not get team for app %q", key.AppName(app))
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			appDesc,
			prometheus.GaugeValue,
			gaugeValue,
			app.Name,
			app.Namespace,
			app.Status.Version,
			app.Status.Release.Status,
			team,
			// Getting version from spec, not status since the version in the spec is the desired version.
			app.Spec.Version,
			strconv.FormatBool(app.Spec.Version != app.Status.Version),
			app.Spec.Catalog,
			app.Spec.Name,
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

// getTeam returns the team to assign for this app CR. It first checks the App CR.
// If the team annotation does not exist it checks the AppCatalogEntry CR. Finally
// it returns the default team so metrics always have a team.
func (a *App) getTeam(ctx context.Context, app v1alpha1.App) (string, error) {
	var team string

	// Team annotation on the App CR takes precedence if it exists.
	if key.AppTeam(app) != "" {
		return key.AppTeam(app), nil
	}

	lo := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,%s=%s,%s=%s",
			label.CatalogName,
			key.CatalogName(app),
			label.AppKubernetesName,
			key.AppName(app),
			label.AppKubernetesVersion,
			key.Version(app)),
	}
	aces, err := a.k8sClient.G8sClient().ApplicationV1alpha1().AppCatalogEntries(metav1.NamespaceAll).List(ctx, lo)
	if len(aces.Items) == 0 {
		// No result so use default team.
		return a.defaultTeam, nil
	} else if len(aces.Items) > 1 {
		// Something is wrong so use default team but log the dupe result.
		a.logger.Debugf(ctx, "expected 1 result for catalog %#q app %#q version %#q but found %d",
			key.CatalogName(app),
			key.AppName(app),
			key.Version(app),
			len(aces.Items))
		return a.defaultTeam, nil
	} else if err != nil {
		return "", microerror.Mask(err)
	}

	ace := aces.Items[0]

	// Owners annotation takes precedence if it exists.
	ownersYAML := key.AppCatalogEntryOwners(ace)

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

			if team != "" {
				return team, nil
			}
		}
	}

	team = key.AppCatalogEntryTeam(ace)
	if team == "" {
		// If there is no team annotation we use the default.
		team = a.defaultTeam
	}

	return team, nil
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
