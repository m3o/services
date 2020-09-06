package builder

type Config struct {
	BaseImageURL     string
	BuildImageURL    string
	BuildRegistryURL string
	RegistryUsername string
	RegistryPassword string
}

// build has what we need to render the Dockerfile template:
type build struct {
	BaseImage     string
	BuildImage    string
	SourceGitRepo string
}
