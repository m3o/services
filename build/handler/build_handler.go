package handler

// BuildHandler implements the build service interface:
type BuildHandler struct {
	baseImageURL string
}

// New returns an initialised BuildHandler:
func New(baseImageURL string) (*BuildHandler, error) {

	// Eventually we will log in to the docker registry here (and return an error if that doesn't work):

	return &BuildHandler{
		baseImageURL: baseImageURL,
	}, nil
}
