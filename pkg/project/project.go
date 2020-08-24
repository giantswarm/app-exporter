package project

var (
	description = "The app-exporter does something."
	gitSHA      = "n/a"
	name        = "app-exporter"
	source      = "https://github.com/giantswarm/app-exporter"
	version     = "0.1.0-dev"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
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
