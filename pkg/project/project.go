package project

var (
	description = "The app-collector does something."
	gitSHA      = "n/a"
	name        = "app-collector"
	source      = "https://github.com/giantswarm/app-collector"
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
