//go:build smoke
// +build smoke

package ats

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/k8sportforward/v2"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/app-exporter/pkg/project"
)

const (
	namespace  string = metav1.NamespaceDefault
	serverPort int    = 8000
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

			KubeConfigPath: os.Getenv("ATS_KUBE_CONFIG_PATH"),
		}

		k8sClients, err = k8sclient.NewClients(c)
		if err != nil {
			t.Fatalf("could not create k8sclients %v", err)
		}
	}

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
		logger.Debugf(ctx, "waiting for %#q pod", project.Name())

		podName, err = waitForPod(ctx, k8sClients)
		if err != nil {
			t.Fatalf("could not get %#q pod %#v", project.Name(), err)
		}

		logger.Debugf(ctx, "waited for %#q pod", project.Name())
	}

	var tunnel *k8sportforward.Tunnel
	{
		logger.Debugf(ctx, "creating tunnel to pod %#q on port %d", podName, serverPort)

		tunnel, err = fw.ForwardPort(namespace, podName, serverPort)
		if err != nil {
			t.Fatalf("could not create tunnel %v", err)
		}

		logger.Debugf(ctx, "created tunnel to pod %#q on port %d", podName, serverPort)
	}

	var metricsResp *http.Response
	{
		metricsURL := fmt.Sprintf("http://%s/metrics", tunnel.LocalAddress())

		logger.Debugf(ctx, "getting metrics from %#q", metricsURL)

		metricsResp, err = waitForServer(metricsURL)
		if err != nil {
			t.Fatalf("server didn't come up on time")
		}

		if metricsResp.StatusCode != http.StatusOK {
			t.Fatalf("expected http status %#q got %#q", http.StatusOK, metricsResp.StatusCode)
		}

		logger.Debugf(ctx, "got metrics from %#q", metricsURL)
	}

	var app *v1alpha1.App
	{
		app, err = k8sClients.G8sClient().ApplicationV1alpha1().Apps(namespace).Get(ctx, project.Name(), metav1.GetOptions{})
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}

		var appVersion string

		if app.Status.AppVersion != app.Status.Version {
			appVersion = app.Status.AppVersion
		}

		expectedAppMetric := fmt.Sprintf("app_operator_app_info{app=\"%s\",app_version=\"%s\",catalog=\"%s\",deployed_version=\"%s\",latest_version=\"%s\",name=\"%s\",namespace=\"%s\",status=\"%s\",team=\"batman\",upgrade_available=\"%s\",version=\"%s\",version_mismatch=\"%s\"} 1",
			app.Spec.Name,
			appVersion,
			app.Spec.Catalog,
			app.Status.Version, // deployed_version
			"",                 // latest_version is empty
			app.Name,
			app.Namespace,
			app.Status.Release.Status,
			"false",          // upgrade_avaiable is false
			app.Spec.Version, // version is the desired version
			strconv.FormatBool(app.Spec.Version != app.Status.Version))

		logger.Debugf(ctx, "checking for expected app metric\n%s", expectedAppMetric)

		respBytes, err := ioutil.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}

		metrics := string(respBytes)
		if !strings.Contains(metrics, expectedAppMetric) {
			t.Fatalf("expected app metric\n\n%s\n\nnot found in response\n\n%s", expectedAppMetric, metrics)
		}

		logger.Debugf(ctx, "found expected app metric")

		expectedAppOperatorMetric := "app_operator_ready_total{namespace=\"giantswarm\",version=\"0.0.0\"} 1"

		logger.Debugf(ctx, "checking for expected app-operator metric\n%s", expectedAppOperatorMetric)

		if !strings.Contains(metrics, expectedAppOperatorMetric) {
			t.Fatalf("expected app metric\n\n%s\n\nnot found in response\n\n%s", expectedAppOperatorMetric, metrics)
		}

		logger.Debugf(ctx, "found expected app-operator metric")
	}
}

func waitForPod(ctx context.Context, k8sClients *k8sclient.Clients) (string, error) {
	var err error
	var podName string

	o := func() error {
		lo := metav1.ListOptions{
			FieldSelector: "status.phase=Running",
			LabelSelector: fmt.Sprintf("app=%s", project.Name()),
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
