package project

var (
	description = "The pss-operator does something."
	gitSHA      = "n/a"
	name        = "pss-operator"
	source      = "https://github.com/giantswarm/pss-operator"
	version     = "0.2.1"
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
