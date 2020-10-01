// +build k8srequired

package key

func AppExporterAppName() string {
	return "app-exporter"
}

func ControlPlaneTestCatalogName() string {
	// TODO Put back test catalog.
	return "control-plane-catalog"
}

func ControlPlaneTestCatalogStorageURL() string {
	// TODO Put back test catalog.
	return "https://giantswarm.github.io/control-plane-catalog/"
}

func Namespace() string {
	return "giantswarm"
}

func ServerPort() int {
	return 8000
}
