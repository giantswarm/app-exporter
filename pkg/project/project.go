package project

const (
	operatorName = "app-operator"
)

var (
	description = "The app-exporter is a Prometheus exporter for the Giant Swarm App Platform."
	gitSHA      = "n/a"
	name        = "app-exporter"
	source      = "https://github.com/giantswarm/app-exporter"
	version     = "0.19.1"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

// Helm2AppVersion is always 1.0.0 for workload cluster app CRs using Helm 2.
func Helm2AppVersion() string {
	return "1.0.0"
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
