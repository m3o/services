package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
	platform "github.com/micro/services/platform/service/proto"
	utils "github.com/micro/services/platform/web/util"
)

var (
	// DefaultVersion is the default version of the service
	// the assume if none is specified
	DefaultVersion = "latest"
	// DefaultNamespace is the default namespace of the services,
	// this will eventually be loaded from config
	DefaultNamespace = "go.micro"
)

// Handler encapsulates the events handlers
type Handler struct {
	platform platform.PlatformService
}

// RegisterHandlers adds the GitHub webhook handlers to the service
func RegisterHandlers(srv web.Service) error {
	h := Handler{
		platform: platform.NewPlatformService("go.micro.service.platform", srv.Options().Service.Client()),
	}

	srv.HandleFunc("/v1/github/events", h.eventsHandler)
	return nil
}

// processBuildEvent processes build events fired through github actions
func (h *Handler) eventsHandler(w http.ResponseWriter, req *http.Request) {
	// Extract the request body containing the webhook data
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.Write500(w, err)
		log.Errorf("Failed to read webhook data: %v", err)
		return
	}

	// Unmarshal the bytes into a struct
	var data *githubWebhook
	if err := json.Unmarshal(body, &data); err != nil {
		utils.Write500(w, err)
		log.Errorf("Failed to unmarshal webhook data: %v", err)
		return
	}

	// Event type indicates what stage of the build is occurring
	var evType platform.EventType
	switch ev := req.Header.Get("Micro-Event"); ev {
	case "build.started":
		evType = platform.EventType_BuildStarted
	case "build.finished":
		evType = platform.EventType_BuildFinished
	case "build.failed":
		evType = platform.EventType_BuildFailed
	default:
		// unknown event
		log.Errorf("Unknown event type: %s", ev)
		utils.Write500(w, errors.New("unknown event type"))
		return
	}

	// buildID is the github actions build ID, e.g. 46017067
	buildID := req.Header.Get("X-Github-Build")
	log.Infof("Processing %v build event #%v", evType, buildID)

	// repoURL is the source of the code, e.g. "github.com/micro/services"
	repoURL := strings.TrimPrefix(data.Payload.Repository.URL, "https://")

	// commitID is the ID of the last ID. In an ideal world we'd pass
	// all of the IDs in the metadata. e.g. "974b680d403ad7d5594ca812a146b3bd342c089b"
	var commitID string
	if len(data.CommitIDs) > 0 {
		commitID = data.CommitIDs[0]
	}

	// metadata is passed with each event
	metadata := map[string]string{
		"build":  buildID,  // github action number
		"commit": commitID, // commit hash
		"repo":   repoURL,  // github.com/micro/services
	}

	// createEvent is an function which encapsulates the logic to create the event
	// in the platform service, as it uses lots of immutable variables, this function
	// was declared inline so they all didn't need to be passed on every function call.
	createEvent := func(srv string, event platform.EventType) {
		// create the event in the platform service
		_, err := h.platform.CreateEvent(req.Context(), &platform.CreateEventRequest{
			Event: &platform.Event{
				Type:      event,
				Timestamp: time.Now().Unix(),
				Service: &platform.Service{
					Name:    nameForService(srv),
					Version: DefaultVersion,
					Source:  path.Join(repoURL, srv),
				},
				Metadata: metadata,
			},
		})

		// Handle the error which was returned by the platform service, since the error
		// could be a one-off, we don't abort the request, however we do write a 500 code
		// to the response. Note, since this function is called multiple times, the error
		// will be overriden if another one occurs.
		if err != nil {
			log.Errorf("Unable to create event type %v for service %v: %v", event.String(), srv, err)
			utils.Write500(w, err)
		} else {
			log.Infof("Created %v event for service %v", event.String(), srv)
		}
	}

	// We only want to create source changed events once, so we do this on
	// build started, since this only happens as a result of source changing.
	if evType == platform.EventType_BuildStarted {
		for _, srv := range data.Services.Created {
			createEvent(srv, platform.EventType_SourceCreated)
		}

		for _, srv := range data.Services.Updated {
			createEvent(srv, platform.EventType_SourceUpdated)
		}

		for _, srv := range data.Services.Deleted {
			createEvent(srv, platform.EventType_SourceDeleted)
		}
	}

	// Create the build event for all the created and modified services,
	// deleted services don't need this event since they can no longer be
	// built
	for _, srv := range data.Services.Created {
		createEvent(srv, evType)
	}

	for _, srv := range data.Services.Updated {
		createEvent(srv, evType)
	}

	// If the build has finished, create/update/delete all services
	if evType == platform.EventType_BuildFinished {
		for _, srv := range data.Services.Created {
			_, err := h.platform.CreateService(req.Context(), &platform.CreateServiceRequest{
				Service: &platform.Service{
					Name:    srv,
					Source:  path.Join(repoURL, srv),
					Version: DefaultVersion,
					Metadata: map[string]string{
						"build": commitID,
					},
				},
			})

			if err != nil {
				log.Errorf("Unable to create service %v: %v", srv, err)
				utils.Write500(w, err)
			}
		}

		for _, srv := range data.Services.Updated {
			_, err := h.platform.UpdateService(req.Context(), &platform.UpdateServiceRequest{
				Service: &platform.Service{
					Name:    srv,
					Source:  path.Join(repoURL, srv),
					Version: DefaultVersion,
					Metadata: map[string]string{
						"build": commitID,
					},
				},
			})

			if err != nil {
				log.Errorf("Unable to update service %v: %v", srv, err)
				utils.Write500(w, err)
			}
		}

		for _, srv := range data.Services.Deleted {
			_, err := h.platform.DeleteService(req.Context(), &platform.DeleteServiceRequest{
				Service: &platform.Service{
					Name:    srv,
					Source:  path.Join(repoURL, srv),
					Version: DefaultVersion,
				},
			})

			if err != nil {
				log.Errorf("Unable to delete service %v: %v", srv, err)
				utils.Write500(w, err)
			}
		}
	}
}

// nameForService determines the name of the service from the directory path,
// e.g. foo/bar becomes go.micro.srv.foo-bar and foo/api becomes go.micro.api.foo
func nameForService(srv string) string {
	var name string
	if strings.HasSuffix(srv, "web") {
		name = fmt.Sprintf("%v.web.%v", DefaultNamespace, strings.ReplaceAll(srv, "/web", ""))
	} else if strings.HasSuffix(srv, "api") {
		name = fmt.Sprintf("%v.api.%v", DefaultNamespace, strings.ReplaceAll(srv, "/api", ""))
	} else {
		name = fmt.Sprintf("%v.srv.%v", DefaultNamespace, srv)
	}
	return strings.ReplaceAll(name, "/", "-")
}

type githubWebhook struct {
	Services struct {
		Created []string `json:"added"`
		Updated []string `json:"modified"`
		Deleted []string `json:"removed"`
	} `json:"services"`
	CommitIDs []string `json:"commit_ids"`
	Payload   struct {
		Repository struct {
			URL string `json:"url"`
		} `json:"repository"`
	} `json:"payload"`
}
