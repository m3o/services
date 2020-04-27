package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/event/service/proto"
)

// Event implement the event service interface
type Event struct {
	name  string
	store store.Store
}

// New returns an initialised event handler
func New(service micro.Service) *Event {
	return &Event{
		name:  service.Name(),
		store: store.DefaultStore,
	}
}

// Create an event
func (e *Event) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Type == pb.EventType_Unknown {
		return errors.BadRequest(e.name, "Missing type")
	}
	if len(req.ProjectId) == 0 {
		return errors.BadRequest(e.name, "Missing project id")
	}

	// construct the event
	event := &pb.Event{
		Id:        uuid.New().String(),
		Type:      req.Type,
		ProjectId: req.ProjectId,
		Created:   time.Now().Unix(),
		Metadata:  req.Metadata,
	}

	// write the event to the store
	bytes, err := json.Marshal(event)
	if err != nil {
		return errors.InternalServerError(e.name, "Error marshaling json: %v", err)
	}
	key := event.ProjectId + "/" + event.Id
	if err := e.store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(e.name, "Error writing to the store: %v", err)
	}

	return nil
}

// Read events
func (e *Event) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if len(req.ProjectId) == 0 {
		return errors.BadRequest(e.name, "Missing project id")
	}

	// lookup the projects matching this prefix, if the event
	// id is blank all the projects events will be returned
	prefix := req.ProjectId + "/" + req.EventId
	recs, err := e.store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError(e.name, "Error reading from the store: %v", err)
	}

	// unmarshal the records
	rsp.Events = make([]*pb.Event, 0, len(recs))
	for _, r := range recs {
		var event *pb.Event
		if err := json.Unmarshal(r.Value, &event); err != nil {
			return errors.InternalServerError(e.name, "Error unmarshaling json: %v", err)
		}
		rsp.Events = append(rsp.Events, event)
	}

	return nil
}
