package collector

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8sclient/v8/pkg/k8sclienttest"
	"github.com/giantswarm/k8smetadata/pkg/annotation"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
	prometheustest "github.com/prometheus/client_golang/prometheus/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	clientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// fakeCollector implements prometheus.Collector interface and
// is wrapper for the App that implements exporterkit.Collector
type fakeCollector struct {
	app *App
}

func (fc fakeCollector) Collect(ch chan<- prometheus.Metric) {
	_ = fc.app.Collect(ch)
}

func (fc fakeCollector) Describe(ch chan<- *prometheus.Desc) {
	_ = fc.app.Describe(ch)
}

func Test_convertToTime(t *testing.T) {
	expectedTime, err := time.Parse(time.RFC3339, "2019-12-31T23:59:59Z")
	if err != nil {
		t.Errorf("time.Parse err = %v", err)
	}

	tests := []struct {
		name         string
		datetime     string
		expected     time.Time
		errorMatcher func(error) bool
	}{
		{
			name:     "case 1: normal timestamp parsing",
			datetime: "2019-12-31T23:59:59.000",
			expected: expectedTime,
		},
		{
			name:         "case 2: parsing error since unknown ",
			datetime:     "2019-12-31T23:59:59Z",
			errorMatcher: IsInvalidExecution,
		},
		{
			name:         "case 3: parsing error as wrong date",
			datetime:     "2019-13-31T23:59:59Z",
			errorMatcher: IsInvalidExecution,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convertToTime(tc.datetime)
			switch {
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case err != nil && !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("convertToTime() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func Test_collectAppStatus(t *testing.T) {
	tests := []struct {
		name                 string
		apps                 []*v1alpha1.App
		catalogs             []*v1alpha1.Catalog
		catalogsEntries      []*v1alpha1.AppCatalogEntry
		expectedMetrics      string
		expectedMetricsCount int
	}{
		{
			name: "flawless",
			apps: []*v1alpha1.App{
				newApp("hello-world-app", "giantswarm", "hello-world", "0.3.0", "", "", nil, nil),
				newApp("example", "customer", "default", "1.0.0", "", "", nil, nil),
				newApp("test-app", "default", "test-app", "1.0.0", "", "", nil, map[string]string{
					label.Cluster: "foo",
				}),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
				newCatalog("default", "giantswarm"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE("hello-world-app", "giantswarm", "default", "0.3.0", "", "", true),
				newACE("example", "customer", "default", "1.0.0", "", "", true),
				newACE("test-app", "default", "giantswarm", "1.0.0", "", "", true),
			},
			expectedMetrics:      "testdata/expected.1",
			expectedMetricsCount: 3,
		},
		{
			name: "flawless with v* versions",
			apps: []*v1alpha1.App{
				newApp("hello-world-app", "giantswarm", "hello-world", "v0.3.0", "", "", nil, nil),
				newApp("example", "customer", "default", "v1.0.0", "", "", nil, nil),
				newApp("test-app", "default", "test-app", "v1.0.0", "", "", nil, map[string]string{
					label.Cluster: "foo",
				}),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
				newCatalog("default", "giantswarm"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE("hello-world-app", "giantswarm", "default", "v0.3.0", "", "", true),
				newACE("example", "customer", "default", "v1.0.0", "", "", true),
				newACE("test-app", "default", "giantswarm", "1.0.0", "", "", true),
			},
			expectedMetrics:      "testdata/expected.1",
			expectedMetricsCount: 3,
		},
		// It is RECOMMENDED to use mixed versions in the next tests
		{
			name: "app pending update",
			apps: []*v1alpha1.App{
				newApp("hello-world-app", "giantswarm", "hello-world", "0.3.0", "", "", nil, nil),
				newApp("example", "customer", "default", "v1.0.0", "v0.9.0", "", nil, nil),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE("hello-world-app", "giantswarm", "default", "0.3.0", "", "", true),
				newACE("example", "customer", "default", "v1.0.0", "", "", true),
			},
			expectedMetrics:      "testdata/expected.2",
			expectedMetricsCount: 2,
		},
		{
			name: "non existing ACE, taking team from annotations",
			apps: []*v1alpha1.App{
				newApp("hello-world-app", "giantswarm", "hello-world", "0.3.0", "", "", nil, nil),
				newApp("atlas-app", "giantswarm", "default", "0.9.0", "", "", map[string]string{annotation.AppTeam: "team-atlas"}, nil),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE("hello-world-app", "giantswarm", "default", "0.3.0", "", "", true),
				newACE("atlas-app", "giantswarm", "default", "1.0.0", "", "", true),
			},
			expectedMetrics:      "testdata/expected.3",
			expectedMetricsCount: 2,
		},
		{
			name: "non existing ACE, taking team from labels",
			apps: []*v1alpha1.App{
				newApp("hello-world-app", "giantswarm", "hello-world", "0.3.0", "", "", nil, nil),
				newApp("atlas-app", "giantswarm", "default", "0.9.0", "", "", nil, map[string]string{annotation.AppTeam: "team-atlas"}),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE("hello-world-app", "giantswarm", "default", "0.3.0", "", "", true),
				newACE("atlas-app", "giantswarm", "default", "1.0.0", "", "", true),
			},
			expectedMetrics:      "testdata/expected.3",
			expectedMetricsCount: 2,
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			var err error

			gsObj := make([]runtime.Object, 0)
			for _, ct := range tc.catalogs {
				gsObj = append(gsObj, ct)
			}

			for _, cte := range tc.catalogsEntries {
				gsObj = append(gsObj, cte)
			}

			for _, app := range tc.apps {
				gsObj = append(gsObj, app)
			}

			var k8sClientFake *k8sclienttest.Clients
			{
				schemeBuilder := runtime.SchemeBuilder{
					v1alpha1.AddToScheme,
				}

				err = schemeBuilder.AddToScheme(scheme.Scheme)
				if err != nil {
					t.Fatal(err)
				}

				k8sClientFake = k8sclienttest.NewClients(k8sclienttest.ClientsConfig{
					CtrlClient: clientfake.NewClientBuilder().
						WithScheme(scheme.Scheme).
						WithRuntimeObjects(gsObj...).
						Build(),
				})
			}

			appConfig := AppConfig{
				K8sClient: k8sClientFake,
				Logger:    microloggertest.New(),

				DefaultTeam:         "honeybadger",
				Provider:            "aws",
				RetiredTeamsMapping: map[string]string{},
			}

			app, err := NewApp(appConfig)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			fakeColl := fakeCollector{
				app: app,
			}
			num := prometheustest.CollectAndCount(
				fakeColl,
				prometheus.BuildFQName(namespace, "app", "info"),
			)
			if num != tc.expectedMetricsCount {
				t.Errorf("expected %d metrics to collect, got %d", tc.expectedMetricsCount, num)
			}

			expected, err := os.Open(tc.expectedMetrics)
			if err != nil {
				panic(err)
			}
			defer expected.Close()

			err = prometheustest.CollectAndCompare(
				fakeColl,
				expected,
				prometheus.BuildFQName(namespace, "app", "info"),
			)
			if err != nil {
				t.Errorf("unexpected collecting result:\n %s", err)
			}
		})
	}
}

func Test_getLatestAppVersions(t *testing.T) {
	tests := []struct {
		name             string
		catalogs         []*v1alpha1.Catalog
		catalogsEntries  []*v1alpha1.AppCatalogEntry
		expectedVersions map[string]string
	}{
		{
			name: "flawless",
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE("hello-world-app", "giantswarm", "default", "0.3.0", "", "", true),
				newACE("hello-world-app", "giantswarm", "default", "0.2.0", "", "", false),
				newACE("example", "customer", "default", "1.0.0", "", "", true),
				newACE("example", "customer", "default", "0.2.0", "", "", false),
				newACE("example", "customer", "default", "0.1.0", "", "", false),
			},
			expectedVersions: map[string]string{
				"customer-example":           "1.0.0",
				"giantswarm-hello-world-app": "0.3.0",
			},
		},
		{
			name: "flawless with v* versions",
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE("hello-world-app", "giantswarm", "default", "0.3.0", "", "", true),
				newACE("hello-world-app", "giantswarm", "default", "0.2.0", "", "", false),
				newACE("example", "customer", "default", "v1.0.0", "", "", true),
				newACE("example", "customer", "default", "v0.2.0", "", "", false),
				newACE("example", "customer", "default", "v0.1.0", "", "", false),
			},
			expectedVersions: map[string]string{
				"customer-example":           "1.0.0",
				"giantswarm-hello-world-app": "0.3.0",
			},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			var err error

			gsObj := make([]runtime.Object, 0)
			for _, ct := range tc.catalogs {
				gsObj = append(gsObj, ct)
			}

			for _, cte := range tc.catalogsEntries {
				gsObj = append(gsObj, cte)
			}

			var k8sClientFake *k8sclienttest.Clients
			{
				schemeBuilder := runtime.SchemeBuilder{
					v1alpha1.AddToScheme,
				}

				err = schemeBuilder.AddToScheme(scheme.Scheme)
				if err != nil {
					t.Fatal(err)
				}

				k8sClientFake = k8sclienttest.NewClients(k8sclienttest.ClientsConfig{
					CtrlClient: clientfake.NewClientBuilder().
						WithScheme(scheme.Scheme).
						WithRuntimeObjects(gsObj...).
						Build(),
				})
			}

			appConfig := AppConfig{
				K8sClient: k8sClientFake,
				Logger:    microloggertest.New(),

				DefaultTeam:         "honeybadger",
				Provider:            "aws",
				RetiredTeamsMapping: map[string]string{},
			}

			app, err := NewApp(appConfig)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			latestAppVersions, err := app.getLatestAppVersions(context.TODO())
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			if !reflect.DeepEqual(latestAppVersions, tc.expectedVersions) {
				t.Fatalf("want matching resources \n %s", cmp.Diff(latestAppVersions, tc.expectedVersions))
			}
		})
	}
}

func Test_getTeamMappings(t *testing.T) {
	tests := []struct {
		name                 string
		apps                 []v1alpha1.App
		catalogs             []*v1alpha1.Catalog
		catalogsEntries      []*v1alpha1.AppCatalogEntry
		expectedTeamMappings map[string]string
		retiredTeamsMapping  map[string]string
	}{
		{
			name: "flawless",
			apps: []v1alpha1.App{
				*newApp("hello-world-app", "giantswarm", "hello-world", "0.2.0", "", "", nil, nil),
				*newApp("example", "customer", "default", "1.0.0", "", "", nil, nil),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE(
					"hello-world-app", "giantswarm", "default", "0.3.0",
					"[{'team':'honeybadger','catalog':'giantswarm'}]", "", true,
				),
				newACE(
					"hello-world-app", "giantswarm", "default", "0.2.0",
					"[{'team':'honeybadger','catalog':'giantswarm'}]", "", false,
				),
				newACE("example", "customer", "default", "1.0.0", "", "customer-team", true),
				newACE("example", "customer", "default", "0.2.0", "", "customer-team", false),
				newACE("example", "customer", "default", "0.1.0", "", "customer-team", false),
			},
			expectedTeamMappings: map[string]string{
				"customer-example-1.0.0":           "customer-team",
				"giantswarm-hello-world-app-0.2.0": "honeybadger",
			},
			retiredTeamsMapping: map[string]string{},
		},
		{
			name: "flawless with v* versions",
			apps: []v1alpha1.App{
				*newApp("hello-world-app", "giantswarm", "hello-world", "0.2.0", "", "", nil, nil),
				*newApp("example", "customer", "default", "v1.0.0", "", "", nil, nil),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE(
					"hello-world-app", "giantswarm", "default", "0.3.0",
					"[{'team':'honeybadger','catalog':'giantswarm'}]", "", true,
				),
				newACE(
					"hello-world-app", "giantswarm", "default", "0.2.0",
					"[{'team':'honeybadger','catalog':'giantswarm'}]", "", false,
				),
				newACE("example", "customer", "default", "v1.0.0", "", "customer-team", true),
				newACE("example", "customer", "default", "v0.2.0", "", "customer-team", false),
				newACE("example", "customer", "default", "v0.1.0", "", "customer-team", false),
			},
			expectedTeamMappings: map[string]string{
				"customer-example-1.0.0":           "",
				"giantswarm-hello-world-app-0.2.0": "honeybadger",
			},
			retiredTeamsMapping: map[string]string{},
		},
		{
			name: "flawless with team mappings",
			apps: []v1alpha1.App{
				*newApp("hello-world-app", "giantswarm", "hello-world", "0.3.0", "", "", nil, nil),
				*newApp("example", "customer", "default", "v1.0.0", "", "", nil, nil),
			},
			catalogs: []*v1alpha1.Catalog{
				newCatalog("giantswarm", "default"),
				newCatalog("customer", "default"),
			},
			catalogsEntries: []*v1alpha1.AppCatalogEntry{
				newACE(
					"hello-world-app",
					"giantswarm",
					"default",
					"0.3.0",
					"[{'team':'batman','catalog':'giantswarm'}]",
					"",
					true,
				),
				newACE(
					"example",
					"customer",
					"default",
					"v1.0.0",
					"",
					"customer-team",
					true,
				),
			},
			expectedTeamMappings: map[string]string{
				"customer-example-1.0.0":           "",
				"giantswarm-hello-world-app-0.3.0": "honeybadger",
			},
			retiredTeamsMapping: map[string]string{
				"batman": "honeybadger",
			},
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			var err error

			gsObj := make([]runtime.Object, 0)
			for _, ct := range tc.catalogs {
				gsObj = append(gsObj, ct)
			}

			for _, cte := range tc.catalogsEntries {
				gsObj = append(gsObj, cte)
			}

			var k8sClientFake *k8sclienttest.Clients
			{
				schemeBuilder := runtime.SchemeBuilder{
					v1alpha1.AddToScheme,
				}

				err = schemeBuilder.AddToScheme(scheme.Scheme)
				if err != nil {
					t.Fatal(err)
				}

				k8sClientFake = k8sclienttest.NewClients(k8sclienttest.ClientsConfig{
					CtrlClient: clientfake.NewClientBuilder().
						WithScheme(scheme.Scheme).
						WithRuntimeObjects(gsObj...).
						Build(),
				})
			}

			appConfig := AppConfig{
				K8sClient: k8sClientFake,
				Logger:    microloggertest.New(),

				DefaultTeam:         "honeybadger",
				Provider:            "aws",
				RetiredTeamsMapping: tc.retiredTeamsMapping,
			}

			app, err := NewApp(appConfig)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			teamsMappings, err := app.getTeamMappings(context.TODO(), tc.apps)
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}

			if !reflect.DeepEqual(teamsMappings, tc.expectedTeamMappings) {
				t.Fatalf("want matching resources \n %s", cmp.Diff(teamsMappings, tc.expectedTeamMappings))
			}
		})
	}
}

func newACE(app, catalog, namespace, version, owners, team string, latest bool) *v1alpha1.AppCatalogEntry {
	ace := v1alpha1.AppCatalogEntry{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AppCatalogEntry",
			APIVersion: "application.giantswarm.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"application.giantswarm.io/owners": "[{team:test,catalog:giantswarm}]",
			},
			Labels: map[string]string{
				label.CatalogName:          catalog,
				label.AppKubernetesName:    app,
				label.AppKubernetesVersion: version,
			},
			Name:      fmt.Sprintf("%s-%s-%s", catalog, app, version),
			Namespace: namespace,
		},
		Spec: v1alpha1.AppCatalogEntrySpec{
			AppName: app,
			Catalog: v1alpha1.AppCatalogEntrySpecCatalog{
				Name: catalog,
			},
			Version: version,
		},
	}

	if latest {
		ace.ObjectMeta.Labels["latest"] = "true"
	}

	if owners != "" {
		ace.ObjectMeta.Annotations[annotation.AppOwners] = owners
	}

	if team != "" {
		ace.ObjectMeta.Annotations[annotation.AppTeam] = team
	}

	return &ace
}

func newApp(name, catalog, namespace, version, statusVersion, statusRelease string, annotations, labels map[string]string) *v1alpha1.App {
	if statusVersion == "" {
		statusVersion = version
	}

	if statusRelease == "" {
		statusRelease = "deployed"
	}

	app := v1alpha1.App{
		TypeMeta: metav1.TypeMeta{
			Kind:       "App",
			APIVersion: "application.giantswarm.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Annotations: annotations,
			Labels:      labels,
			Name:        name,
			Namespace:   namespace,
		},
		Spec: v1alpha1.AppSpec{
			Name:      name,
			Namespace: namespace,
			Catalog:   catalog,
			Version:   version,
		},
		Status: v1alpha1.AppStatus{
			Release: v1alpha1.AppStatusRelease{
				Status: statusRelease,
			},
			Version: statusVersion,
		},
	}

	return &app
}

func newCatalog(name, namespace string) *v1alpha1.Catalog {
	catalog := v1alpha1.Catalog{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Catalog",
			APIVersion: "application.giantswarm.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				label.CatalogVisibility: "public",
			},
			Name:      name,
			Namespace: namespace,
		},
	}

	return &catalog
}
