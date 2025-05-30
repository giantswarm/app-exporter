//go:build functional
// +build functional

package ats

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/giantswarm/k8smetadata/pkg/label"
	v1 "k8s.io/api/core/v1"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/k8sportforward/v2"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/app-exporter/pkg/project"

	expkey "github.com/giantswarm/app-exporter/internal/key"
)

const (
	namespace  = metav1.NamespaceDefault
	serverPort = 8000
)

// TestMetrics checks the exporter emits app info metrics for its own app CR.
//
// - Waits for the pod to start and creates a port forwarding connection
// to the metrics endpoint.
// - Scrapes the metrics and checks the expected app info metric is present.
func TestMetrics(t *testing.T) {
	var err error

	ctx := context.Background()

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			t.Fatalf("could not create logger %v", err)
		}
	}

	var k8sClients *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			Logger: logger,
			SchemeBuilder: k8sclient.SchemeBuilder{
				v1alpha1.AddToScheme,
			},

			KubeConfigPath: KubeConfigPath(),
		}

		k8sClients, err = k8sclient.NewClients(c)
		if err != nil {
			t.Fatalf("could not create k8sclients %v", err)
		}
	}

	catalogCR := &v1alpha1.Catalog{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "giantswarm",
			Labels: map[string]string{
				label.AppOperatorVersion: "0.0.0",
			},
		},
		Spec: v1alpha1.CatalogSpec{
			Description: "default",
			Title:       "default",
			Repositories: []v1alpha1.CatalogSpecRepository{
				{
					Type: "helm",
					URL:  "https://giantswarm.github.io/default-catalog",
				},
			},
			Storage: v1alpha1.CatalogSpecStorage{
				Type: "helm",
				URL:  "https://giantswarm.github.io/default-catalog",
			},
		},
	}
	err = k8sClients.CtrlClient().Create(ctx, catalogCR)
	if err != nil {
		t.Fatalf("failed to create default catalog: %#v", err)
	}

	{
		err := k8sClients.CtrlClient().Create(ctx, &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-app",
			},
		})
		if err != nil {
			t.Fatalf("failed to create test-app namespace: %#v", err)
		}
	}

	{
		testAppUserValues := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app-user-values",
				Namespace: "test-app",
			},
			Data: map[string]string{
				"values": "namespace: test-app",
			},
		}

		_, err := k8sClients.K8sClient().CoreV1().ConfigMaps("test-app").Create(ctx, testAppUserValues, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed to create test-app-user-valies config map: %#v", err)
		}
	}

	{
		err := k8sClients.CtrlClient().Create(ctx, &v1alpha1.App{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-app",
				Labels: map[string]string{
					label.AppOperatorVersion: "0.0.0",
					label.AppKubernetesName:  "test-app",
					label.App:                "test-app",
					label.Cluster:            "kind",
					"foo":                    "bar",
				},
			},
			Spec: v1alpha1.AppSpec{
				Catalog: "default",
				KubeConfig: v1alpha1.AppSpecKubeConfig{
					InCluster: true,
				},
				Name:      "test-app",
				Namespace: "test-app",
				Version:   "1.0.0",
				UserConfig: v1alpha1.AppSpecUserConfig{
					ConfigMap: v1alpha1.AppSpecUserConfigConfigMap{
						Name:      "test-app-user-values",
						Namespace: "test-app",
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("failed to create test-app app cr: %#v", err)
		}
	}

	t.Logf("Waiting for test-app to come up...")

	testAppPodName, err := waitForPod(ctx, k8sClients, "test-app", "test-app")
	if err != nil {
		t.Fatalf("could not get test-app pod %#v", err)
	}

	t.Logf("Waited for test-app to come up: %s", testAppPodName)

	var fw *k8sportforward.Forwarder
	{
		c := k8sportforward.ForwarderConfig{
			RestConfig: k8sClients.RESTConfig(),
		}

		fw, err = k8sportforward.NewForwarder(c)
		if err != nil {
			t.Fatalf("could not create forwarder %v", err)
		}
	}

	var podName string
	{
		t.Logf("waiting for %#q pod", project.Name())

		podName, err = waitForPod(ctx, k8sClients, namespace, project.Name())
		if err != nil {
			t.Fatalf("could not get %#q pod %#v", project.Name(), err)
		}

		t.Logf("waited for %#q pod: %s", project.Name(), podName)
	}

	var tunnel *k8sportforward.Tunnel
	{
		t.Logf("creating tunnel to pod %#q on port %d", podName, serverPort)

		tunnel, err = fw.ForwardPort(namespace, podName, serverPort)
		if err != nil {
			t.Fatalf("could not create tunnel %v", err)
		}

		t.Logf("created tunnel to pod %#q on port %d", podName, serverPort)
	}

	var metricsResp *http.Response
	{
		metricsURL := fmt.Sprintf("http://%s/metrics", tunnel.LocalAddress())

		t.Logf("getting metrics from %#q", metricsURL)

		metricsResp, err = waitForServer(metricsURL)
		if err != nil {
			t.Fatalf("server didn't come up on time")
		}

		if metricsResp.StatusCode != http.StatusOK {
			t.Fatalf("expected http status %#q got %#q", http.StatusOK, metricsResp.StatusCode)
		}

		t.Logf("got metrics from %#q", metricsURL)
	}

	var app *v1alpha1.App
	var testApp *v1alpha1.App
	{
		app = &v1alpha1.App{}
		err = k8sClients.CtrlClient().Get(ctx, types.NamespacedName{Namespace: namespace, Name: project.Name()}, app)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}

		var appVersion string

		if app.Status.AppVersion != app.Status.Version {
			appVersion = expkey.FormatVersion(app.Status.AppVersion)
		}

		expectedAppExporterMetric := fmt.Sprintf("app_operator_app_info{app=\"%s\",app_version=\"%s\",catalog=\"%s\",cluster_id=\"%s\",cluster_missing=\"%s\",deployed_version=\"%s\",latest_version=\"%s\",name=\"%s\",namespace=\"%s\",status=\"%s\",team=\"noteam\",upgrade_available=\"%s\",version=\"%s\",version_mismatch=\"%s\"} 1",
			app.Spec.Name,
			appVersion,
			app.Spec.Catalog,
			"",
			"false",
			expkey.FormatVersion(app.Status.Version), // deployed_version
			"",                                       // latest_version is empty
			app.Name,
			app.Namespace,
			app.Status.Release.Status,
			"false",                                // upgrade_available is false
			expkey.FormatVersion(app.Spec.Version), // version is the desired version
			strconv.FormatBool(app.Spec.Version != app.Status.Version))

		t.Logf("Expected app-exporter metrics:\n%s", expectedAppExporterMetric)

		testApp = &v1alpha1.App{}
		err = k8sClients.CtrlClient().Get(ctx, types.NamespacedName{Namespace: "test-app", Name: "test-app"}, testApp)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}

		expectedTestAppMetric := fmt.Sprintf("app_operator_app_info{app=\"%s\",app_version=\"%s\",catalog=\"%s\",cluster_id=\"%s\",cluster_missing=\"%s\",deployed_version=\"%s\",latest_version=\"%s\",name=\"%s\",namespace=\"%s\",status=\"%s\",team=\"honeybadger\",upgrade_available=\"%s\",version=\"%s\",version_mismatch=\"%s\"} 1",
			testApp.Spec.Name,
			"2.13.0",
			testApp.Spec.Catalog,
			"kind",
			"false",
			expkey.FormatVersion(testApp.Status.Version), // deployed_version
			"", // latest_version is empty
			testApp.Name,
			testApp.Namespace,
			testApp.Status.Release.Status,
			"false", // upgrade_available is false
			expkey.FormatVersion(testApp.Spec.Version), // version is the desired version
			strconv.FormatBool(testApp.Spec.Version != testApp.Status.Version))

		t.Logf("Expected test-app metrics:\n%s", expectedTestAppMetric)

		respBytes, err := ioutil.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}

		metrics := string(respBytes)

		t.Logf("METRICS RESPONSE START")
		t.Logf("%s", metrics)
		t.Logf("METROCS RESPONSE END")

		if !strings.Contains(metrics, expectedAppExporterMetric) {
			t.Fatalf("expected app (app-exporter) metric\n\n%s\n\nnot found in response\n\n%s", expectedAppExporterMetric, metrics)
		}

		if !strings.Contains(metrics, expectedTestAppMetric) {
			t.Fatalf("expected app (test-app) metric\n\n%s\n\nnot found in response\n\n%s", expectedTestAppMetric, metrics)
		}

		t.Logf("found expected app metric")

		expectedAppOperatorMetric := "app_operator_ready_total{namespace=\"giantswarm\",version=\"0.0.0\"} 1"

		t.Logf("checking for expected app-operator metric\n%s", expectedAppOperatorMetric)

		if !strings.Contains(metrics, expectedAppOperatorMetric) {
			t.Fatalf("expected app metric\n\n%s\n\nnot found in response\n\n%s", expectedAppOperatorMetric, metrics)
		}

		t.Logf("found expected app-operator metric")
	}
}

func waitForPod(ctx context.Context, k8sClients *k8sclient.Clients, namespace string, appLabelValue string) (string, error) {
	var err error
	var podName string

	o := func() error {
		lo := metav1.ListOptions{
			FieldSelector: "status.phase=Running",
			LabelSelector: fmt.Sprintf("app=%s", appLabelValue),
		}
		pods, err := k8sClients.K8sClient().CoreV1().Pods(namespace).List(ctx, lo)
		if err != nil {
			return microerror.Mask(err)
		}
		if len(pods.Items) != 1 {
			return microerror.Maskf(executionFailedError, "expected 1 pod got %d", len(pods.Items))
		}

		pod := pods.Items[0]
		podName = pod.Name

		return nil
	}

	n := func(err error, t time.Duration) {
		log.Printf("waiting for pod for %s: %#v", t, err)
	}

	err = backoff.RetryNotify(o, backoff.NewConstant(5*time.Minute, 15*time.Second), n)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return podName, nil
}

func waitForServer(url string) (*http.Response, error) {
	var err error
	var resp *http.Response

	o := func() error {
		resp, err = http.Get(url)
		if err != nil {
			return microerror.Maskf(executionFailedError, "could not retrieve %s: %v", url, err)
		}

		return nil
	}

	n := func(err error, t time.Duration) {
		log.Printf("waiting for server at %s: %v", t, err)
	}

	err = backoff.RetryNotify(o, backoff.NewConstant(5*time.Minute, 15*time.Second), n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return resp, nil
}
