// +build k8srequired

package metrics

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sportforward/v2"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/app-exporter/integration/key"
	"github.com/giantswarm/app-exporter/pkg/project"
)

// TestMetrics TODO2
//
func TestMetrics(t *testing.T) {
	var err error

	ctx := context.Background()

	var fw *k8sportforward.Forwarder
	{
		c := k8sportforward.ForwarderConfig{
			RestConfig: config.K8sClients.RESTConfig(),
		}

		fw, err = k8sportforward.NewForwarder(c)
		if err != nil {
			t.Fatalf("could not create forwarder %v", err)
		}
	}

	var podName string
	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waiting for %#q pod", project.Name()))

		podName, err = waitForPod(ctx)
		if err != nil {
			t.Fatalf("could not get %#q pod %#v", project.Name(), err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waited for %#q pod", project.Name()))
	}

	var tunnel *k8sportforward.Tunnel
	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating tunnel to pod %#q on port %d", podName, key.ServerPort()))

		tunnel, err = fw.ForwardPort(key.Namespace(), podName, key.ServerPort())
		if err != nil {
			t.Fatalf("could not create tunnel %v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created tunnel to pod %#q on port %d", podName, key.ServerPort()))
	}

	var metricsResp *http.Response
	{
		metricsURL := fmt.Sprintf("http://%s/metrics", tunnel.LocalAddress())

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("getting metrics from %#q", metricsURL))

		metricsResp, err = waitForServer(metricsURL)
		if err != nil {
			t.Fatalf("server didn't come up on time")
		}

		if metricsResp.StatusCode != http.StatusOK {
			t.Fatalf("expected http status %#q got %#q", http.StatusOK, metricsResp.StatusCode)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("got metrics from %#q", metricsURL))
	}

	var app *v1alpha1.App
	{
		app, err = config.K8sClients.G8sClient().ApplicationV1alpha1().Apps(key.Namespace()).Get(ctx, project.Name(), metav1.GetOptions{})
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}

		expectedMetric := fmt.Sprintf("app_operator_app_info{catalog=\"%s\",name=\"%s\",namespace=\"%s\",status=\"%s\",team=\"batman\",version=\"%s\"} 1",
			app.Spec.Catalog,
			app.Name,
			app.Namespace,
			app.Status.Release.Status,
			app.Spec.Version)

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("checking for expected metric\n%s", expectedMetric))

		respBytes, err := ioutil.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}

		metrics := string(respBytes)
		if !strings.Contains(metrics, expectedMetric) {
			t.Fatalf("expected metric\n\n%s\n\nnot found in response\n\n%s", expectedMetric, metrics)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", "found expected metric")
	}
}

func waitForPod(ctx context.Context) (string, error) {
	var err error
	var podName string

	o := func() error {
		lo := metav1.ListOptions{
			FieldSelector: "status.phase=Running",
			LabelSelector: fmt.Sprintf("app=%s", project.Name()),
		}
		pods, err := config.K8sClients.K8sClient().CoreV1().Pods(key.Namespace()).List(ctx, lo)
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
