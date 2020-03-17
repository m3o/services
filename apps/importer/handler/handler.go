package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"gopkg.in/yaml.v2"

	pb "github.com/micro/services/apps/importer/proto/importer"
	apps "github.com/micro/services/apps/service/proto/apps"
)

// File is a document in GitHub
type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	URL  string `json:"download_url"`
}

// Handler implements the importer service interface
type Handler struct {
	apps apps.AppsService
}

// NewHandler returns an initialised handdler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		apps: apps.NewAppsService("go.micro.service.apps", srv.Client()),
	}
}

// Import updates the apps service to use the latest apps as define in the micro/apps github repo
func (h *Handler) Import(ctx context.Context, req *pb.ImportRequest, rsp *pb.ImportResponse) error {
	// Fetch all the existing apps from the apps service
	aRsp, err := h.apps.List(context.TODO(), &apps.ListRequest{})
	if err != nil {
		return err
	}
	existingApps := aRsp.GetApps()

	// Get the files from the GitHub API
	newApps, err := h.getAppsFromGitHub()
	if err != nil {
		return err
	}

	// Determine if any existing apps need to be deactivated, since
	// their ID is no longer present in the repo
	if err := h.deactivateDeletedApps(context.TODO(), existingApps, newApps); err != nil {
		return err
	}

	// Loop through all new apps and either update it if it exists in the
	// existingApps slice, or create it in the apps service
	if err := h.updateAllApps(context.TODO(), existingApps, newApps); err != nil {
		return err
	}

	return nil
}

// getAppsFromGitHub fetches the .yaml files in the Micro/Apps GitHub repo
// and decodes them into Apps Service App objects.
func (h *Handler) getAppsFromGitHub() ([]*apps.App, error) {
	resp, err := http.Get("https://api.github.com/repos/micro/apps/contents")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var allFiles []File
	if err := json.NewDecoder(resp.Body).Decode(&allFiles); err != nil {
		return nil, err
	}

	// Get the download URLs for yaml files
	var appFiles []File
	for _, f := range allFiles {
		if strings.HasSuffix(f.Path, ".yaml") {
			appFiles = append(appFiles, f)
		}
	}

	// Loop over the urls and construct the apps from the YAML
	newApps := make([]*apps.App, 0, len(appFiles))
	for _, f := range appFiles {
		resp, err := http.Get(f.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var app apps.App
		if err := yaml.Unmarshal(body, &app); err != nil {
			return nil, err
		}

		newApps = append(newApps, &app)
	}

	return newApps, nil
}

func (h *Handler) deactivateDeletedApps(ctx context.Context, old []*apps.App, new []*apps.App) error {
	for _, app := range old {
		shouldBeDeactived := func() bool {
			for _, a := range new {
				if a.Id == app.Id {
					return false
				}
			}

			return app.Active
		}()

		if shouldBeDeactived {
			log.Infof("Deactivating App: %v", app.Id)

			app.Active = false
			_, err := h.apps.Update(ctx, &apps.UpdateRequest{App: app})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *Handler) updateAllApps(ctx context.Context, old []*apps.App, new []*apps.App) error {
	// Create or update all the apps listed in the repos
	for _, app := range new {
		var oldApp *apps.App
		for _, a := range old {
			if app.Id == a.Id {
				oldApp = a
				break
			}
		}

		if oldApp == nil {
			// Create a new app
			log.Infof("Creating App: %v", app.Id)
			_, err := h.apps.Create(context.TODO(), &apps.CreateRequest{App: app})
			if err != nil {
				return err
			}
		} else {
			// Update an existing app
			log.Infof("Updating App: %v", app.Id)
			_, err := h.apps.Update(context.TODO(), &apps.UpdateRequest{App: app})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
