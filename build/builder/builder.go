package builder

// Builder does the actual building of docker images:
type Builder interface {
	Build(sourceGitRepo, targetImageTag string) error
}

// build has what we need to render the Dockerfile template:
type build struct {
	BaseImage     string
	BuildImage    string
	SourceGitRepo string
}
