package handler

import (
	"context"
	"encoding/json"
	"strings"

	proto "github.com/m3o/services/explore/proto/explore"
	regproto "github.com/micro/micro/v3/proto/registry"
	"github.com/micro/micro/v3/service/model"
	"github.com/micro/micro/v3/service/registry"
	regutil "github.com/micro/micro/v3/service/registry/util"
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
	services, err := e.reg.ListServices(registry.ListDomain("backend"))
	if err != nil {
		return err
	}

	rsp.Services = []*regproto.Service{}

	// Very rudimentary search result ranking
	// prioritize name and endoint name matches
	matchedName := []*regproto.Service{}
	matchedEndpointName := []*regproto.Service{}
	matchedOther := []*regproto.Service{}

	metas := []*proto.SaveMetaRequest{}
	err = e.meta.Read(model.QueryAll(), &metas)
	if err != nil {
		return err
	}

	for _, service := range services {
		if req.SearchTerm != "" {
			rsp.Services = append(rsp.Services, regutil.ToProto(service))
			continue
		}
		if strings.Contains(service.Name, req.SearchTerm) {
			matchedName = append(matchedName, regutil.ToProto(service))
			continue
		}
		for _, ep := range service.Endpoints {
			if strings.Contains(ep.Name, req.SearchTerm) {
				matchedEndpointName = append(matchedEndpointName, regutil.ToProto(service))
				continue
			}
		}
		js, _ := json.Marshal(service)
		meta := &proto.SaveMetaRequest{}
		for _, m := range metas {
			if m.ServiceName == service.Name {
				meta = m
			}
		}
		if strings.Contains(string(js), req.SearchTerm) ||
			strings.Contains(meta.OpenAPIJSON, req.SearchTerm) ||
			strings.Contains(meta.Readme, req.SearchTerm) {
			matchedOther = append(matchedOther, regutil.ToProto(service))
		}
	}

	if len(rsp.Services) == 0 {
		rsp.Services = append(rsp.Services, matchedName...)
		rsp.Services = append(rsp.Services, matchedEndpointName...)
		rsp.Services = append(rsp.Services, matchedOther...)
	}
	return nil
}

// make this endpoint only callable by admins
func (e *Explore) SaveMeta(ctx context.Context, req *proto.SaveMetaRequest, rsp *proto.SaveMetaResponse) error {
	return e.meta.Create(req)
}
