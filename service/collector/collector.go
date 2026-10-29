package collector

const (
	gaugeValue         float64 = 1
	namespace          string  = "app_operator"
	notInstalledStatus string  = "not-installed"

	// ociAppTeam is the team annotation of the OCI-compliant chart metadata
	// format, see https://docs.giantswarm.io/reference/platform-api/chart-metadata/.
	// Charts migrated to it carry it instead of annotation.AppTeam, which
	// k8smetadata does not define.
	ociAppTeam string = "io.giantswarm.application.team"
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
	labelClusterId        = "cluster_id"
)
