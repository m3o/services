package handler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/runtime"
	runtimepb "github.com/micro/go-micro/v2/runtime/service/proto"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/platform/service/proto"
)

var (
	// Topic aysnc messages are published to
	Topic = "go.micro.runtime.events"

	// The base image for our builds
	Image = "docker.pkg.github.com/micro/services"
)

type status struct {
	*pb.Status
}

func (s *status) Ok() bool {
	if len(s.Status.Error) > 0 {
		return false
	}
	return true
}

// Handler implements the platform service interface
type Handler struct {
	// current status of the platform
	sync.RWMutex
	status map[string]*status

	// all the things
	Auth     auth.Auth
	Broker   broker.Broker
	Config   config.Config
	Registry registry.Registry
	Store    store.Store
	Runtime  runtime.Runtime
	Event    micro.Event
}

// NewHandler returns an initialized Handler
func NewHandler(srv micro.Service) *Handler {
	h := &Handler{
		status:   make(map[string]*status),
		Auth:     srv.Options().Auth,
		Broker:   srv.Options().Broker,
		Config:   srv.Options().Config,
		Event:    micro.NewEvent(Topic, srv.Client()),
		Registry: srv.Options().Registry,
		Store:    srv.Options().Store,
		Runtime:  runtime.DefaultRuntime,
	}

	err := micro.RegisterSubscriber(
		Topic,
		srv.Server(),
		h.HandleEvent,
		server.SubscriberQueue("queue.platform"),
	)

	if err != nil {
		log.Errorf("Error subscribing to registry: %v", err)
	}

	// Check the status of the platform every 30 seconds
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		// check the status once
		h.checkStatus()

		for _ = range ticker.C {
			h.checkStatus()
		}
	}()

	return h
}

// HandleEvent such as service created, updated or deleted. It reformats
// the request to match the proto and then passes it off to the handler to process
// as it would any other request, ensuring there is no duplicate logic.
func (h *Handler) HandleEvent(ctx context.Context, event *runtimepb.Event) error {
	req := &pb.CreateEventRequest{Event: &pb.Event{
		Type:      RuntimeEventsMap[event.Type],
		Timestamp: event.Timestamp,
		Service: &pb.Service{
			Name:    nameForService(event.Service),
			Version: event.Version,
		},
	}}

	return h.CreateEvent(ctx, req, &pb.CreateEventResponse{})
}

var (
	// RuntimeEventsMap stores the map of runtime event types
	// and the service event type.
	RuntimeEventsMap = map[string]pb.EventType{
		"create": pb.EventType_ServiceCreated,
		"update": pb.EventType_ServiceUpdated,
		"delete": pb.EventType_ServiceDeleted,
	}
	// DefaultNamespace is the default namespace of the services,
	// this will eventually be loaded from config
	DefaultNamespace = "go.micro"
)

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
