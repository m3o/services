package manager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v30/github"
	log "github.com/micro/go-micro/v2/logger"
)

const (
	owner = "micro"
	repo  = "services"
)

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

func (m *manager) getChangedServices(commit string) ([]string, error) {
	log.Debugf("Listing files for commit %v", commit)

	repoCommit, _, err := m.client.Repositories.GetCommit(context.Background(), owner, repo, commit)
	if err != nil {
		return nil, err
	}
	folders := map[string]struct{}{}
	for _, file := range repoCommit.Files {
		fname := file.GetFilename()

		// do not care about files or hidden folders
		if strings.Contains(fname, "/") && !strings.HasPrefix(fname, ".") {
			folders[fname] = struct{}{}
		}
	}
	ret := []string{}
	for k := range folders {
		ret = append(ret, k)
	}
	return ret, nil
}

func (m *manager) updateService(serviceName, commit string) error {
	return nil
}

func (m *manager) Run() {
	t := time.NewTicker(time.Second * 3)
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

			serviceNames, err := m.getChangedServices(latest)
			if err != nil {
				log.Errorf("Can't get services from commit", err)
			}

			// perform an update
			for _, serviceName := range serviceNames {
				if err := m.updateService(serviceName, latest); err != nil {
					log.Errorf("Error updating latest: %v", err)
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
