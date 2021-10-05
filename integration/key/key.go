//go:build k8srequired
// +build k8srequired

package key

func ControlPlaneCatalogName() string {
	return "control-plane-catalog"
}

func ControlPlaneTestCatalogName() string {
	return "control-plane-test-catalog"
}

func Namespace() string {
	return "giantswarm"
}

func ServerPort() int {
	return 8000
}
