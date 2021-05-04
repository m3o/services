package handler

import (
	"context"
	"strings"
	"sync"
	"time"

	proto "github.com/m3o/services/explore/proto"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	"github.com/micro/micro/v3/service/registry"
	regutil "github.com/micro/micro/v3/service/registry/util"
)

type Explore struct {
	reg  registry.Registry
	meta model.Model

	// store a cache for the index
	sync.RWMutex
	cache       []*proto.Service
	lastUpdated time.Time
}

// NewexploreHandler returns a explore handler configured to report the explore of the given services
func NewExploreHandler(reg registry.Registry) *Explore {
	m := model.NewModel(
		model.WithKey("ServiceName"),
		//model.WithNamespace("meta"),
	)
	m.Register(proto.SaveMetaRequest{})

	// create the handler
	e := &Explore{
		reg:  reg,
		meta: m,
	}

	// create the index
	e.Index(context.TODO(), &proto.IndexRequest{}, &proto.IndexResponse{})

	// return the handler
	return e
}

func (e *Explore) Index(ctx context.Context, req *proto.IndexRequest, rsp *proto.IndexResponse) error {
	// get the existing cache
	e.RLock()
	cache := e.cache
	lastUpdate := e.lastUpdated
	e.RUnlock()

	// do the lookup and cache now if it doesn't exist or is older than a minute
	if len(cache) == 0 || time.Since(lastUpdate) > time.Minute {
		e.Lock()
		defer e.Unlock()

		// get the services
		resp := new(proto.SearchResponse)
		err := e.Search(context.TODO(), &proto.SearchRequest{}, resp)
		if err != nil {
			return err
		}

		// set the cache
		e.cache = resp.Services
		e.lastUpdated = time.Now()

		// update cache val
		cache = resp.Services
	}

	// check the limit
	limit := int(req.Limit)
	if limit <= 0 || limit > len(cache) {
		limit = len(cache)
	}

	// otherwise use the cache
	for i := 0; i < limit; i++ {
		// set and return the response
		rsp.Services = append(rsp.Services, cache[i])
	}

	return nil

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
				Service:      regutil.ToProto(service),
				Readme:       meta.Readme,
				OpenAPIJSON:  meta.OpenAPIJSON,
				ExamplesJSON: meta.ExamplesJSON,
			})
			continue
		}
		if strings.Contains(service.Name, req.SearchTerm) {
			matchedName = append(matchedName, &proto.Service{
				Service:      regutil.ToProto(service),
				Readme:       meta.Readme,
				OpenAPIJSON:  meta.OpenAPIJSON,
				ExamplesJSON: meta.ExamplesJSON,
			})
			continue
		}
		found := false
		for _, ep := range service.Endpoints {
			if strings.Contains(strings.ToLower(ep.Name), req.SearchTerm) {
				matchedEndpointName = append(matchedEndpointName, &proto.Service{
					Service:      regutil.ToProto(service),
					Readme:       meta.Readme,
					OpenAPIJSON:  meta.OpenAPIJSON,
					ExamplesJSON: meta.ExamplesJSON,
				})
				found = true
				break
			}
		}
		if found {
			continue
		}
		if strings.Contains(meta.Readme, req.SearchTerm) {
			matchedOther = append(matchedOther, &proto.Service{
				Service:      regutil.ToProto(service),
				Readme:       meta.Readme,
				OpenAPIJSON:  meta.OpenAPIJSON,
				ExamplesJSON: meta.ExamplesJSON,
			})
		}
	}

	// these will be empty if the search term is empty
	rsp.Services = append(rsp.Services, matchedName...)
	rsp.Services = append(rsp.Services, matchedEndpointName...)
	rsp.Services = append(rsp.Services, matchedOther...)

	return nil
}

func (e *Explore) Service(ctx context.Context, req *proto.ServiceRequest, rsp *proto.ServiceResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("explore.service", "missing name")
	}

	// check the cache for the service
	var service *proto.Service
	e.RLock()
	for _, s := range e.cache {
		if s.Service.Name == req.Name {
			service = s
			break
		}
	}
	e.RUnlock()

	if service != nil {
		rsp.Service = service
		return nil
	}

	// get the service from the registry
	srv, err := e.reg.GetService(req.Name)
	if err != nil {
		return err
	}

	if len(srv) == 0 {
		return errors.NotFound("explore.service", "service not found")
	}

	// get service metadata
	meta := new(proto.SaveMetaRequest)
	err = e.meta.Read(model.QueryEquals("ServiceName", req.Name), &meta)
	if err != nil {
		return err
	}

	// return response
	rsp.Service = &proto.Service{
		Service:      regutil.ToProto(srv[0]),
		Readme:       meta.Readme,
		OpenAPIJSON:  meta.OpenAPIJSON,
		ExamplesJSON: meta.ExamplesJSON,
	}

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
