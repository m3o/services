package manager

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v30/github"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
)

const (
	owner = "micro"
	repo  = "services"
	// DefaultVersion is the default version of the service
	// the assume if none is specified
	DefaultVersion = "latest"
	// DefaultNamespace is the default namespace of the services,
	// this will eventually be loaded from config
	DefaultNamespace = "go.micro"
)

type serviceStatus string

var (
	serviceStatusCreated serviceStatus = "created"
	serviceStatusUpdated serviceStatus = "updated"
	serviceStatusDeleted serviceStatus = "deleted"
)

type githubFileChangeStatus string

// a list of github file status changes.
// not documented in the github API
var (
	githubFileChangeStatusCreated  githubFileChangeStatus = "created"
	githubFileChangeStatusChanged  githubFileChangeStatus = "changed"
	githubFileChangeStatusModified githubFileChangeStatus = "modified"
	githubFileChangeStatusRemoved  githubFileChangeStatus = "removed"
)

type fileToStatus struct {
	fileName string
	status   githubFileChangeStatus
}

type manager struct {
	// workflow file name
	workflow string
	// latest commit
	latest string
	// last updated
	updated time.Time

	// the github client
	// TODO: pluggable source
	client *github.Client
}

func (m *manager) lastCommitForWorkflow() (string, error) {
	log.Debug("Listing workflows")
	w, _, err := m.client.Actions.ListWorkflowRunsByFileName(
		context.Background(),
		owner,
		repo,
		m.workflow,
		&github.ListWorkflowRunsOptions{
			Status: "success",
			Branch: "master",
		},
	)
	// failed to get the workflow run
	if err != nil {
		return "", err
	}

	if len(w.WorkflowRuns) == 0 {
		return "", fmt.Errorf("no workflows")
	}

	// we got one
	wr := w.WorkflowRuns[0]

	// TODO: process different statuses
	if *wr.Status != "completed" {
		return m.latest, nil
	}

	if *wr.Conclusion != "success" {
		return m.latest, nil
	}

	return *wr.HeadSHA, nil
}

// returns a map key -> values of serviceName -> serviceStatus
func (m *manager) getChangedFolders(commit string) (map[string]serviceStatus, error) {
	log.Debugf("Listing files for commit %v", commit)

	repoCommit, _, err := m.client.Repositories.GetCommit(context.Background(), owner, repo, commit)
	if err != nil {
		return nil, err
	}
	filesToStatuses := []fileToStatus{}
	for _, v := range repoCommit.Files {
		filesToStatuses = append(filesToStatuses, fileToStatus{
			fileName: v.GetFilename(),
			status:   githubFileChangeStatus(v.GetStatus()),
		})
	}
	return folderStatuses(filesToStatuses), nil
}

// maps github file change statuses to folders and their deployment status
// ie. "asim/scheduler/main.go" "removed" will become "asim/scheduler" "deleted"
func folderStatuses(statuses []fileToStatus) map[string]serviceStatus {
	folders := map[string]serviceStatus{}
	// Prioritize main.go creates and deletes
	for _, status := range statuses {
		fname := status.fileName
		status := status.status
		if !strings.HasSuffix(fname, "main.go") {
			continue
		}
		fold := path.Dir(fname)

		_, exists := folders[fold]
		if exists {
			continue
		}
		if status == "created" {
			folders[fold] = serviceStatusCreated
		} else if status == "removed" {
			folders[fold] = serviceStatusDeleted
		}

	}
	// continue with normal file changes for service updates
	for _, status := range statuses {
		fname := status.fileName
		folds := topFolders(fname)
		for _, fold := range folds {
			_, exists := folders[fold]
			if exists {
				continue
			}
			folders[fold] = serviceStatusUpdated
		}
	}
	return folders
}

// from path returns the top level dirs to be deployed
// ie.
func topFolders(path string) []string {
	parts := strings.Split(path, "/")
	ret := []string{parts[0]}
	if len(parts) > 2 {
		ret = append(ret, filepath.Join(parts[0], parts[1]))
	}
	return ret
}

func (m *manager) updateService(folderPath, commit string, status serviceStatus) error {
	service := &runtime.Service{
		Name:    folderPath,
		Version: commit,
	}
	typ := typeFromFolder(folderPath)
	image := "micro/services"

	switch status {
	case serviceStatusCreated:
		opts := []runtime.CreateOption{
			// create a specific service type
			runtime.CreateType(typ),
			runtime.CreateImage(image),
		}

		if err := runtime.DefaultRuntime.Create(service, opts...); err != nil {
			return err
		}
	case serviceStatusUpdated:
		if err := runtime.DefaultRuntime.Update(service); err != nil {
			return err
		}
	case serviceStatusDeleted:
		if err := runtime.DefaultRuntime.Delete(service); err != nil {
			return err
		}
	}
	return fmt.Errorf("Unrecognized service status: '%v'", status)
}

func typeFromFolder(folder string) string {
	if strings.Contains(folder, "api") {
		return "api"
	}
	if strings.Contains(folder, "web") {
		return "web"
	}
	return "service"
}

func (m *manager) Run() {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	// every minute we look for changes and apply any updates
	for {
		select {
		case <-t.C:
			latest, err := m.lastCommitForWorkflow()
			if err != nil {
				log.Errorf("Error checking latest: %v", err)
				continue
			}
			if latest == "" {
				log.Error("Can't get commit hash")
				continue
			}

			// same as last time
			if latest == m.latest {
				continue
			}

			folderStatuses, err := m.getChangedFolders(latest)
			if err != nil {
				log.Errorf("Can't get services from commit", err)
			}

			// perform an update
			for folder, status := range folderStatuses {
				if err := m.updateService(folder, latest, status); err != nil {
					log.Errorf("Error updating service '%v': %v", folder, err)
					continue
				}
			}

			// save the latest
			m.latest = latest
			m.updated = time.Now()
		}
	}
}

// Start the scheduler
func Start(workflowFilename string) {
	m := new(manager)
	m.workflow = workflowFilename
	m.client = github.NewClient(nil)

	go m.Run()
}
