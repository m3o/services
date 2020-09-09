package handler

type changeType int

const (
	created changeType = iota
	deleted
)

// WebhookResponse is the response type of the webhook request
type WebhookResponse struct{}

// WebhookRequest is the payload struct sent from GitHub
type WebhookRequest struct {
	After      string     `json:"after"`  // e.g. 4c4eee3fad645d165817ecbec597be6d24685d54
	Before     string     `json:"before"` // e.g. 1cb75bed2ae11fe6c860e4ec2b73ba70f22210de
	Reference  string     `json:"ref"`    // e.g. refs/heads/master
	Repository repository `json:"repository"`
	Commits    []commit   `json:"commits"`
}

// commit object sent from GitHub
type commit struct {
	ID       string   `json:"id"`       // e.g. fadc31277ec9137a3605f51ebc97d6802e796000
	Added    []string `json:"added"`    // e.g. ["test/main.go", "go.sum"]
	Modified []string `json:"modified"` // e.g. ["test/foo.go"]
	Deleted  []string `json:"deleted"`  // e.g. ["foo/handler/handler.go"]
}

// repository object sent from GitHub
type repository struct {
	Name string `json:"full_name"` // e.g. m3o/services
}
