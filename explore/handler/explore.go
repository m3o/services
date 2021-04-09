package handler

import (
	"context"
	"encoding/json"
	"strings"

	proto "github.com/m3o/services/explore/proto/explore"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/model"
	"github.com/micro/micro/v3/service/registry"
	regutil "github.com/micro/micro/v3/service/registry/util"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

type Explore struct {
	reg  registry.Registry
	meta model.Model
}

// NewexploreHandler returns a explore handler configured to report the explore of the given services
func NewExploreHandler(reg registry.Registry) *Explore {
	m := model.NewModel(
		model.WithKey("ServiceName"),
		model.WithNamespace("meta"),
	)
	m.Register(proto.SaveMetaRequest{})

	return &Explore{
		reg:  reg,
		meta: m,
	}
}

func (e *Explore) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	// This endpoint likely won't scale further
	// than a few thousand services on the platform.
	// Let's worry about that later

	// @todo Do some cachin here
	services, err := e.reg.ListServices()
	if err != nil {
		return err
	}

	rsp.Services = []*proto.Service{}

	// Very rudimentary search result ranking
	// prioritize name and endoint name matches
	matchedName := []*proto.Service{}
	matchedEndpointName := []*proto.Service{}
	matchedOther := []*proto.Service{}

	metas := []*proto.SaveMetaRequest{}
	err = e.meta.Read(model.QueryAll(), &metas)
	if err != nil {
		return err
	}
	logger.Infof("Found %v metas", len(metas))

	for _, service := range services {
		meta := &proto.SaveMetaRequest{}
		for _, m := range metas {
			if m.ServiceName == service.Name {
				meta = m
			}
		}

		if req.SearchTerm == "" {
			rsp.Services = append(rsp.Services, &proto.Service{
				Service:     regutil.ToProto(service),
				Readme:      meta.Readme,
				OpenAPIJSON: meta.OpenAPIJSON,
			})
			continue
		}
		if strings.Contains(service.Name, req.SearchTerm) {
			matchedName = append(matchedName, &proto.Service{
				Service:     regutil.ToProto(service),
				Readme:      meta.Readme,
				OpenAPIJSON: meta.OpenAPIJSON,
			})
			continue
		}
		for _, ep := range service.Endpoints {
			if strings.Contains(ep.Name, req.SearchTerm) {
				matchedEndpointName = append(matchedEndpointName, &proto.Service{
					Service:     regutil.ToProto(service),
					Readme:      meta.Readme,
					OpenAPIJSON: meta.OpenAPIJSON,
				})
				continue
			}
		}
		js, _ := json.Marshal(service)
		if strings.Contains(string(js), req.SearchTerm) ||
			strings.Contains(meta.OpenAPIJSON, req.SearchTerm) ||
			strings.Contains(meta.Readme, req.SearchTerm) {
			matchedOther = append(matchedOther, &proto.Service{
				Service:     regutil.ToProto(service),
				Readme:      meta.Readme,
				OpenAPIJSON: meta.OpenAPIJSON,
			})
		}
	}

	// these will be empty if the search term is empty
	rsp.Services = append(rsp.Services, matchedName...)
	rsp.Services = append(rsp.Services, matchedEndpointName...)
	rsp.Services = append(rsp.Services, matchedOther...)

	return nil
}

func (e *Explore) SaveMeta(ctx context.Context, req *proto.SaveMetaRequest, rsp *proto.SaveMetaResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	isAdmin := func(ss []string) bool {
		for _, s := range ss {
			if s == "admin" {
				return true
			}
		}
		return false
	}
	if !ok || !isAdmin(acc.Scopes) {
		return errors.BadRequest("explore.SaveMeta", "Unauthorized")
	}
	return e.meta.Create(req)
}
