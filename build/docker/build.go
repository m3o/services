package docker

// Build has what we need to render the Dockerfile template:
type Build struct {
	BaseImage  string
	BuildImage string
	Source     string
}
