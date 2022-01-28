package collector

const (
	gaugeValue         float64 = 1
	namespace          string  = "app_operator"
	notInstalledStatus string  = "not-installed"
)

const (
	labelApp              = "app"
	labelAppVersion       = "app_version"
	labelCatalog          = "catalog"
	labelClusterMissing   = "cluster_missing"
	labelDeployedVersion  = "deployed_version"
	labelLatestVersion    = "latest_version"
	labelName             = "name"
	labelNamespace        = "namespace"
	labelStatus           = "status"
	labelTeam             = "team"
	labelUpgradeAvailable = "upgrade_available"
	labelVersion          = "version"
	labelVersionMismatch  = "version_mismatch"
)
