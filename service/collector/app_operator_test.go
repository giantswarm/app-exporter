package collector

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclienttest"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/micrologger/microloggertest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	clientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_collectAppVersions(t *testing.T) {
	var tests = []struct {
		name                string
		apps                []*v1alpha1.App
		expectedAppVersions map[string]map[string]bool
		errorMatcher        func(error) bool
	}{
		{
			name: fmt.Sprintf("case 1: All apps have the %#q label", label.AppOperatorVersion),
			apps: []*v1alpha1.App{
				newTestApp("hello-world-1", "abc01", map[string]string{
					label.AppOperatorVersion: "6.3.0",
				}),
				newTestApp("hello-world-2", "xyz01", map[string]string{
					label.AppOperatorVersion: "5.9.2",
				}),
				newTestApp("hello-world-3", "testing", map[string]string{
					label.AppOperatorVersion: "6.3.0",
				}),
			},
			expectedAppVersions: map[string]map[string]bool{
				"6.3.0": {
					"abc01":   true,
					"testing": true,
				},
				"5.9.2": {
					"xyz01": true,
				},
			},
		},
		{
			name: fmt.Sprintf("case 2: Skip app with invalid %#q label", label.AppOperatorVersion),
			apps: []*v1alpha1.App{
				newTestApp("hello-world-1", "abc01", map[string]string{
					label.AppOperatorVersion: "6.3.0",
				}),
				newTestApp("hello-world-2", "xyz01", map[string]string{
					label.AppOperatorVersion: "this-is-not-a-valid-version",
				}),
				newTestApp("hello-world-3", "testing", map[string]string{
					label.AppOperatorVersion: "6.4.0",
				}),
			},
			expectedAppVersions: map[string]map[string]bool{
				"6.3.0": {
					"abc01": true,
				},
				"6.4.0": {
					"testing": true,
				},
			},
		},
		{
			name: fmt.Sprintf("case 3: Skip app with missing %#q label", label.AppOperatorVersion),
			apps: []*v1alpha1.App{
				newTestApp("hello-world-1", "abc01", map[string]string{
					label.AppOperatorVersion: "6.3.0",
				}),
				newTestApp("hello-world-2", "org-example", map[string]string{}),
				newTestApp("hello-world-3", "testing", map[string]string{
					label.AppOperatorVersion: "6.4.0",
				}),
			},
			expectedAppVersions: map[string]map[string]bool{
				"6.3.0": {
					"abc01": true,
				},
				"6.4.0": {
					"testing": true,
				},
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			var err error

			gsObj := make([]runtime.Object, 0)

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

			appOperator := &AppOperator{
				k8sClient: k8sClientFake,
				logger:    microloggertest.New(),
			}

			appVersions, err := appOperator.collectAppVersions(context.TODO())

			switch {
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case err != nil && !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if !reflect.DeepEqual(appVersions, tc.expectedAppVersions) {
				t.Errorf("collectAppVersions() = %v, want %v", appVersions, tc.expectedAppVersions)
			}
		})
	}
}

func newTestApp(name, namespace string, labels map[string]string) *v1alpha1.App {
	app := v1alpha1.App{
		TypeMeta: metav1.TypeMeta{
			Kind:       "App",
			APIVersion: "application.giantswarm.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.AppSpec{
			Name:      name,
			Namespace: namespace,
			Catalog:   "test-catalog",
			Version:   "1.0.0",
		},
	}

	return &app
}
