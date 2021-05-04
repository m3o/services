package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	m3oauth "github.com/m3o/services/pkg/auth"
	pb "github.com/m3o/services/publicapi/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

const (
	prefixName = "apiByName:%s"
	prefixID   = "apiByID:%s"
)

type APIEntry struct {
	ID           string
	Name         string
	Description  string
	OpenAPIJSON  string
	Pricing      map[string]int64 // pricing mapping endpoints to price in 1/100 of a cent
	OwnerID      string
	ExamplesJSON string
}

type Publicapi struct {
}

func NewPublicAPIHandler(srv *service.Service) *Publicapi {
	return &Publicapi{}
}

func (p *Publicapi) Publish(ctx context.Context, request *pb.PublishRequest, response *pb.PublishResponse) error {
	if _, err := m3oauth.VerifyMicroAdmin(ctx, "publicapi.Publish"); err != nil {
		return err
	}
	acc, _ := auth.AccountFromContext(ctx)
	ae := &APIEntry{
		ID:           uuid.New().String(),
		Name:         request.Api.Name,
		Description:  request.Api.Description,
		OpenAPIJSON:  request.Api.OpenApiJson,
		Pricing:      request.Api.Pricing,
		OwnerID:      acc.ID,
		ExamplesJSON: request.Api.ExamplesJson,
	}

	if err := p.updateEntry(ctx, ae); err != nil {
		return err
	}
	// enable auth rules
	// we need two rules, one for the /v1/foo/bar from public internet and one for v1api->foo
	//micro auth create rule --resource="service:v1.helloworld:*" --priority 1 helloworld-v1
	//micro auth create rule --resource="service:helloworld:*" --priority 1 --scope '+' helloworld-internal
	if err := auth.Grant(&auth.Rule{
		ID:    fmt.Sprintf("%s-v1", ae.Name),
		Scope: "",
		Resource: &auth.Resource{
			Name:     fmt.Sprintf("v1.%s", ae.Name),
			Type:     "service",
			Endpoint: "*",
		},
		Access:   auth.AccessGranted,
		Priority: 1,
	}); err != nil {
		log.Errorf("Error adding rule %s", err)
		return errors.InternalServerError("publicapi.Publish", "Error enabling API")
	}

	if err := auth.Grant(&auth.Rule{
		ID:    fmt.Sprintf("%s-internal", ae.Name),
		Scope: "+",
		Resource: &auth.Resource{
			Name:     ae.Name,
			Type:     "service",
			Endpoint: "*",
		},
		Access:   auth.AccessGranted,
		Priority: 1,
	}); err != nil {
		log.Errorf("Error adding rule %s", err)
		return errors.InternalServerError("publicapi.Publish", "Error enabling API")
	}

	// event
	if err := events.Publish("publicapi", pb.Event{Type: "APIPublish",
		ApiPublish: &pb.APIPublishEvent{
			Name: ae.Name,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}
	response.Api = marshal(ae)
	// TODO any other v1api things?
	return nil
}

func (p *Publicapi) updateEntry(ctx context.Context, ae *APIEntry) error {
	b, err := json.Marshal(ae)
	if err != nil {
		return err
	}
	// store it
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf(prefixName, ae.Name),
		Value: b,
	}); err != nil {
		log.Errorf("Error writing to store %s", err)
		return err
	}
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf(prefixID, ae.ID),
		Value: b,
	}); err != nil {
		log.Errorf("Error writing to store %s", err)
		return err
	}

	return nil
}

func (p *Publicapi) Get(ctx context.Context, request *pb.GetRequest, response *pb.GetResponse) error {
	var key string
	if len(request.Id) > 0 {
		key = fmt.Sprintf(prefixID, request.Id)
	} else if len(request.Name) > 0 {
		key = fmt.Sprintf(prefixName, request.Name)
	}
	if len(key) == 0 {
		return errors.BadRequest("publicapi.Get", "ID or name must be specified")
	}
	recs, err := store.Read(key)
	if err != nil {
		log.Errorf("Error reading from store %s", err)
		if err == store.ErrNotFound {
			return errors.NotFound("publicapi.Get", "API not found")
		}
		return errors.InternalServerError("publicapi.Get", "Error retrieving API")
	}
	var ae APIEntry
	if err := json.Unmarshal(recs[0].Value, &ae); err != nil {
		log.Errorf("Error marshalling API %s", err)
		return errors.InternalServerError("publicapi.Get", "Error retrieving API")
	}
	response.Api = marshal(&ae)
	return nil
}

func (p *Publicapi) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	recs, err := store.Read(fmt.Sprintf(prefixName, ""), store.ReadPrefix())
	if err != nil {
		log.Errorf("Error reading from store %s", err)
		return err
	}
	response.Apis = make([]*pb.PublicAPI, len(recs))
	for i, v := range recs {
		var ae APIEntry
		if err := json.Unmarshal(v.Value, &ae); err != nil {
			return err
		}
		response.Apis[i] = marshal(&ae)
	}
	return nil
}

func (p *Publicapi) Remove(ctx context.Context, request *pb.RemoveRequest, response *pb.RemoveResponse) error {
	if _, err := m3oauth.VerifyMicroAdmin(ctx, "publicapi.Remove"); err != nil {
		return err
	}
	var key string
	if len(request.Id) > 0 {
		key = fmt.Sprintf(prefixID, request.Id)
	} else if len(request.Name) > 0 {
		key = fmt.Sprintf(prefixName, request.Name)
	}
	if len(key) == 0 {
		return errors.BadRequest("publicapi.Remove", "ID or name must be specified")
	}
	recs, err := store.Read(key)
	if err != nil {
		return err
	}
	var ae APIEntry
	if err := json.Unmarshal(recs[0].Value, &ae); err != nil {
		return err
	}

	// delete entries
	if err := store.Delete(fmt.Sprintf(prefixID, ae.ID)); err != nil {
		return errors.InternalServerError("publicapi.Remove", "Error removing API")
	}
	if err := store.Delete(fmt.Sprintf(prefixName, ae.Name)); err != nil {
		return errors.InternalServerError("publicapi.Remove", "Error removing API")
	}
	// TODO
	// delete rules

	if err := auth.Revoke(&auth.Rule{
		ID: fmt.Sprintf("%s-v1", ae.Name),
	}); err != nil {
		log.Errorf("Error deleting rule %s", err)
		return errors.InternalServerError("publicapi.Remove", "Error removing API")
	}

	if err := auth.Revoke(&auth.Rule{
		ID: fmt.Sprintf("%s-internal", ae.Name),
	}); err != nil {
		log.Errorf("Error deleting rule %s", err)
		return errors.InternalServerError("publicapi.Remove", "Error removing API")
	}

	// event
	if err := events.Publish("publicapi", pb.Event{Type: "APIRemove",
		ApiRemove: &pb.APIRemoveEvent{
			Name: ae.Name,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}

	return nil
}

func (p *Publicapi) Update(ctx context.Context, request *pb.UpdateRequest, response *pb.UpdateResponse) error {
	if _, err := m3oauth.VerifyMicroAdmin(ctx, "publicapi.Update"); err != nil {
		return err
	}
	// work out which fields have been passed and update
	recs, err := store.Read(fmt.Sprintf(prefixID, request.Api.Id))
	if err != nil {
		if err == store.ErrNotFound {
			return errors.NotFound("publicapi.Update", "API not found with this ID")
		}
		return errors.InternalServerError("publicapi.Update", "Error updating API")
	}
	var ae APIEntry
	if err := json.Unmarshal(recs[0].Value, &ae); err != nil {
		log.Errorf("Error unmarshalling %s", err)
		return errors.InternalServerError("publicapi.Update", "Error updating API")
	}
	if len(request.Api.OpenApiJson) > 0 {
		ae.OpenAPIJSON = request.Api.OpenApiJson
	}
	if len(request.Api.Name) > 0 {
		ae.Name = request.Api.Name
	}
	if len(request.Api.Description) > 0 {
		ae.Description = request.Api.Description
	}
	if len(request.Api.Pricing) > 0 {
		ae.Pricing = request.Api.Pricing
	}
	if len(request.Api.ExamplesJson) > 0 {
		ae.ExamplesJSON = request.Api.ExamplesJson
	}
	if err := p.updateEntry(ctx, &ae); err != nil {
		log.Errorf("Error updating entry %s", err)
		return errors.InternalServerError("publicapi.Update", "Error updating API")
	}
	response.Api = marshal(&ae)
	// TODO update explore service
	return nil
}

func marshal(ae *APIEntry) *pb.PublicAPI {
	return &pb.PublicAPI{
		Id:           ae.ID,
		Name:         ae.Name,
		Description:  ae.Description,
		OpenApiJson:  ae.OpenAPIJSON,
		Pricing:      ae.Pricing,
		OwnerId:      ae.OwnerID,
		ExamplesJson: ae.ExamplesJSON,
	}
}
