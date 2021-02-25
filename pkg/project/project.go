package project

const (
	operatorName = "app-operator"
)

var (
	description = "The app-exporter is a Prometheus exporter for the Giant Swarm App Platform."
	gitSHA      = "n/a"
	name        = "app-exporter"
	source      = "https://github.com/giantswarm/app-exporter"
	version     = "0.2.2-dev"
)

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

// WorkloadAppVersion is always 1.0.0 for workload cluster app CRs using Helm 2.
func WorkloadAppVersion() string {
	return "1.0.0"
}
