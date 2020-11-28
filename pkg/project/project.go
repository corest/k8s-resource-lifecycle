package project

var (
	description = "Command line tool to generate history of CR based on audit log files."
	gitSHA      = "n/a"
	name        = "k8s-resource-lifecycle"
	source      = "https://github.com/corest/k8s-resource-lifecycle"
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
