package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/platform/service/proto"
)

// ListEvents returns recent events, if a resource is provided then this is scoped to their events
func (h *Handler) ListEvents(ctx context.Context, req *pb.ListEventsRequest, rsp *pb.ListEventsResponse) error {
	var records []*store.Record
	var err error
	var prefix string

	// Use a prefix to scope to the resource (if one was provided)
	if req.Service != nil && len(req.Service.Name) > 0 {
		prefix = req.Service.Name
	}

	records, err = h.Store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.service.platform", "unable to read from store: %v", err)
	}

	// Decode the records
	rsp.Events = make([]*pb.Event, len(records))
	for i, r := range records {
		if err := json.Unmarshal(r.Value, &rsp.Events[i]); err != nil {
			return errors.InternalServerError("go.micro.service.platform", "unable to decode records")
		}
	}

	// sort the events
	sort.Slice(rsp.Events, func(i, j int) bool { return rsp.Events[i].Timestamp > rsp.Events[j].Timestamp })

	return nil
}

// GetEvent finds an event by the ID
func (h *Handler) GetEvent(ctx context.Context, req *pb.GetEventRequest, rsp *pb.GetEventResponse) error {
	// Perform the validations
	if req.Event == nil {
		return errors.BadRequest("go.micro.service.platform", "missing event")
	}
	if len(req.Event.Id) == 0 {
		return errors.BadRequest("go.micro.service.platform", "missing event ID")
	}

	// Find the record
	records, err := h.Store.Read(req.Event.Id, store.ReadSuffix())
	if err != nil {
		return errors.InternalServerError("go.micro.service.platform", "unable to read from store")
	}
	if len(records) == 0 {
		return errors.NotFound("go.micro.service.platform", "event not found")
	}

	// Decode the event
	if err := json.Unmarshal(records[0].Value, &rsp.Event); err != nil {
		return errors.InternalServerError("go.micro.service.platform", "unable to decode record")
	}

	return nil
}

// CreateEvent records a new event for a resource
func (h *Handler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest, rsp *pb.CreateEventResponse) error {
	// Perform the validations
	if req.Event == nil {
		return errors.BadRequest("go.micro.service.platform", "missing event")
	}
	if req.Event.Type == pb.EventType_Unknown {
		return errors.BadRequest("go.micro.service.platform", "invalid event type")
	}
	if req.Event.Service == nil || req.Event.Service.Name == "" {
		return errors.BadRequest("go.micro.service.platform", "invalid event service")
	}

	// Construct the event
	event := &Event{req.Event}
	if len(event.Id) == 0 {
		event.Id = uuid.New().String()
	}
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().Unix()
	}

	// Write to the store
	err := h.Store.Write(&store.Record{
		Key:   event.Key(),
		Value: event.Bytes(),
	})
	if err != nil {
		return errors.InternalServerError("go.micro.service.platform", "unable to write to store")
	}

	return nil
}

// Event is the store representation of an event
type Event struct {
	*pb.Event
}

// Key to be used in the store
func (e *Event) Key() string {
	return fmt.Sprintf("%v:%v:%v", e.Service.Name, e.Event.Type.String(), e.Id)
}

// Bytes is the JSON encoded event
func (e *Event) Bytes() []byte {
	b, _ := json.Marshal(e)
	return b
}
