// +build k8srequired

package key

func ControlPlaneTestCatalogName() string {
	return "control-plane-test-catalog"
}

func ControlPlaneTestCatalogStorageURL() string {
	return "https://giantswarm.github.io/control-plane-test-catalog/"
}

func Namespace() string {
	return "giantswarm"
}

func ServerPort() int {
	return 8000
}