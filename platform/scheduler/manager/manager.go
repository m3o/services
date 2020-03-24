package manager

import (
	"fmt"

	"github.com/google/go-github/v30/github"
)

type manager struct {
	// workflow id
	workflow string
	// latest commit
	latest string
	// last updated
	updated time.Time

	// the github client
	// TODO: pluggable source
	client *github.Client
}

// check immediately checks if theres a new build
// it returns the latest commit
func (m *manager) check() (string, error) {
	id, _ := strconv.Atoi(m.workflow)

	w, r, err := m.client.Actions.ListWorkflowRunsByID(
		context.Background(),
		"micro",
		"services",
		int64(id),
		opts,
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

func (m *manager) Run() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	// every minute we look for changes and apply any updates
	for {
		select {
		case <-t.C:
			latest, err := m.check()
			if err != nil {
				fmt.Println("Error checking latest: %v", err)
				continue
			}

			// same as last time
			if latest == m.latest {
				continue
			}

			// perform an update
			if err := m.update(latest); err != nil {
				fmt.Println("Error updating latest: %v", err)
				continue
			}

			// save the latest
			m.latest = latest
			m.updated = time.Now()
		}
	}
}

// Start the scheduler
func Start(id string) {
	m := new(manager)
	m.workflow = id
	m.client = github.NewClient(nil)
	m.latest, _ = m.check()
	m.updated = time.Now()

	go m.Run()
}
