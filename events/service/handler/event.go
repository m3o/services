package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service"
	mstore "github.com/micro/micro/v3/service/store"

	pb "github.com/m3o/services/events/service/proto"
)

// Events implement the event service interface
type Events struct {
	name string
}

// New returns an initialised event handler
func New(service *service.Service) *Events {
	return &Events{
		name: service.Name(),
	}
}

// Create an event
func (e *Events) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Type == pb.EventType_Unknown {
		return errors.BadRequest(e.name, "Missing type")
	}
	if len(req.EnvironmentId) == 0 {
		return errors.BadRequest(e.name, "Missing environment id")
	}

	// construct the event
	event := &pb.Event{
		Id:            uuid.New().String(),
		Type:          req.Type,
		EnvironmentId: req.EnvironmentId,
		Created:       time.Now().Unix(),
		Metadata:      req.Metadata,
	}

	// write the event to the store
	bytes, err := json.Marshal(event)
	if err != nil {
		return errors.InternalServerError(e.name, "Error marshaling json: %v", err)
	}
	key := event.EnvironmentId + "/" + event.Id
	if err := mstore.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(e.name, "Error writing to the store: %v", err)
	}

	return nil
}

// Read events
func (e *Events) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if len(req.EnvironmentId) == 0 {
		return errors.BadRequest(e.name, "Missing environment id")
	}

	// lookup the environments matching this prefix, if the event
	// id is blank all the environments events will be returned
	prefix := req.EnvironmentId + "/" + req.EventId
	recs, err := mstore.Read(prefix, store.ReadPrefix())
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
