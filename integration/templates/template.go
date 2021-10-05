//go:build k8srequired
// +build k8srequired

package templates

// AppExporterValues values required by app-exporter chart.
const AppExporterValues = `Installation:
  V1:
    Provider:
      Kind: aws
    Registry:
      Domain: quay.io`
