package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/m3o/services/explore/proto/explore"
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
	ID          string
	Name        string
	Description string
	OpenAPIJSON string
	UnitPrice   float64
}

type Publicapi struct {
	explSvc explore.ExploreService
}

func NewHandler(srv *service.Service) *Publicapi {
	return &Publicapi{
		explSvc: explore.NewExploreService("explore", srv.Client()),
	}
}

func (p *Publicapi) Publish(ctx context.Context, request *pb.PublishRequest, response *pb.PublishResponse) error {
	if err := verifyAdmin(ctx, "publicapi.Remove"); err != nil {
		return err
	}

	ae := &APIEntry{
		ID:          uuid.New().String(),
		Name:        request.Api.Name,
		Description: request.Api.Description,
		OpenAPIJSON: request.Api.OpenApiJson,
		UnitPrice:   request.Api.PricePerRequest,
	}
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

	// publish to explore API
	if _, err := p.explSvc.SaveMeta(ctx, &explore.SaveMetaRequest{
		ServiceName: ae.Name,
		Readme:      ae.Description,
		OpenAPIJSON: ae.OpenAPIJSON,
	}); err != nil {
		log.Errorf("Error publishing to explore service %s", err)
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
		return errors.InternalServerError("v1api.EnableAPI", "Error enabling API")
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
		return errors.InternalServerError("v1api.EnableAPI", "Error enabling API")
	}

	// event
	if err := events.Publish("publicapi", pb.Event{Type: "APIPublish",
		ApiEnable: &pb.APIEnableEvent{
			Name: ae.Name,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}

	// TODO any other v1api things?
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
	response.Api = &pb.PublicAPI{
		Id:              ae.ID,
		Name:            ae.Name,
		Description:     ae.Description,
		OpenApiJson:     ae.OpenAPIJSON,
		PricePerRequest: ae.UnitPrice,
	}
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
		response.Apis[i] = &pb.PublicAPI{
			Id:              ae.ID,
			Name:            ae.Name,
			Description:     ae.Description,
			OpenApiJson:     ae.OpenAPIJSON,
			PricePerRequest: ae.UnitPrice,
		}
	}
	return nil
}

func (p *Publicapi) Remove(ctx context.Context, request *pb.RemoveRequest, response *pb.RemoveResponse) error {
	if err := verifyAdmin(ctx, "publicapi.Remove"); err != nil {
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
	if err := store.Delete(key); err != nil {
		return errors.InternalServerError("publicapi.Remove", "Error removing API")
	}
	return nil
}

func verifyAdmin(ctx context.Context, method string) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != "micro" {
		return errors.Forbidden(method, "Forbidden")
	}
	for _, s := range acc.Scopes {
		if s == "admin" || s == "service" {
			return nil
		}
	}
	return errors.Forbidden(method, "Forbidden")
}
