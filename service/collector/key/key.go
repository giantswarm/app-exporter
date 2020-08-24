package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/application/v1alpha1"

	"github.com/giantswarm/app-exporter/pkg/annotation"
)

func AppName(customResource v1alpha1.App) string {
	return customResource.Spec.Name
}

func CordonUntil(customResource v1alpha1.App) string {
	return customResource.GetAnnotations()[fmt.Sprintf("%s/%s", annotation.ChartOperatorPrefix, annotation.CordonUntil)]
}

func IsAppCordoned(customResource v1alpha1.App) bool {
	_, reasonOk := customResource.Annotations[fmt.Sprintf("%s/%s", annotation.AppOperatorPrefix, annotation.CordonReason)]
	_, untilOk := customResource.Annotations[fmt.Sprintf("%s/%s", annotation.AppOperatorPrefix, annotation.CordonUntil)]

	if reasonOk && untilOk {
		return true
	} else {
		return false
	}
}

func Namespace(customResource v1alpha1.App) string {
	return customResource.Spec.Namespace
}
