package project

const (
	operatorName = "app-operator"
)

var (
	description = "The app-collector does something."
	gitSHA      = "n/a"
	name        = "app-collector"
	source      = "https://github.com/giantswarm/app-collector"
	version     = "0.1.0-dev"
)

// AppTenantVersion is always 1.0.0 for tenant cluster app CRs using Helm 2.
// For app CRs using Helm 3 we use project.Version().
func AppTenantVersion() string {
	return "1.0.0"
}

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func OperatorName() string {
	return operatorName
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
