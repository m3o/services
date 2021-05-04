package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	pb "github.com/m3o/services/publicapi/proto"
	publicapi "github.com/m3o/services/publicapi/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	regutil "github.com/micro/micro/v3/service/registry/util"
	"github.com/micro/micro/v3/service/store"
)

type Explore struct {
	sync.RWMutex
	cache       []*publicapi.ExploreAPI
	regCache    map[string]*registry.Service
	lastUpdated time.Time
}

func NewExploreAPIHandler(srv *service.Service) *Explore {
	return &Explore{}
}

// Index returns all the APIs
func (e *Explore) Index(ctx context.Context, request *publicapi.IndexRequest, response *publicapi.IndexResponse) error {
	ret, err := e.exploreList()
	if err != nil {
		return err
	}
	response.Apis = ret
	return nil
}

func (e *Explore) exploreList() ([]*publicapi.ExploreAPI, error) {
	e.RLock()
	cache := e.cache
	lastUpdate := e.lastUpdated
	e.RUnlock()
	if len(cache) > 0 && time.Since(lastUpdate) < time.Minute {
		return cache, nil
	}
	e.Lock()
	defer e.Unlock()
	svcs, err := registry.ListServices()
	if err != nil {
		return nil, err
	}
	e.regCache = map[string]*registry.Service{}
	for _, s := range svcs {
		e.regCache[s.Name] = s
	}
	log.Infof("Registry cache is %+v", e.cache)
	e.lastUpdated = time.Now()

	recs, err := store.Read(fmt.Sprintf(prefixName, ""), store.ReadPrefix())
	if err != nil {
		log.Errorf("Error reading from store %s", err)
		return nil, err
	}
	ret := []*pb.ExploreAPI{}
	for _, v := range recs {
		var ae APIEntry
		if err := json.Unmarshal(v.Value, &ae); err != nil {
			return nil, err
		}
		log.Infof("Marshalling %+v", ae)
		reg, ok := e.regCache[ae.Name]
		if !ok {
			// entry without matching registry entry probably means an old entry, ignore
			continue
		}
		ret = append(ret, marshalExploreAPI(&ae, reg))
	}

	e.cache = ret
	return ret, nil
}

func marshalExploreAPI(ae *APIEntry, svc *registry.Service) *publicapi.ExploreAPI {
	return &publicapi.ExploreAPI{
		Api:     marshal(ae),
		Service: regutil.ToProto(svc),
	}

}

// Search returns APIs based on the given search term
func (e *Explore) Search(ctx context.Context, request *publicapi.SearchRequest, response *publicapi.SearchResponse) error {
	list, err := e.exploreList()
	if err != nil {
		return err
	}

	if request.SearchTerm == "" {
		response.Apis = list
		return nil
	}

	// Very rudimentary search result ranking
	// prioritize name and endoint name matches
	matchedName := []*publicapi.ExploreAPI{}
	matchedEndpointName := []*publicapi.ExploreAPI{}
	matchedOther := []*publicapi.ExploreAPI{}

	for _, api := range list {
		if strings.Contains(strings.ToLower(api.Api.Name), request.SearchTerm) {
			matchedName = append(matchedName, api)
			continue
		}
		found := false
		for _, ep := range api.Service.Endpoints {
			if strings.Contains(strings.ToLower(ep.Name), request.SearchTerm) {
				matchedEndpointName = append(matchedEndpointName, api)
				found = true
				break
			}
		}
		if found {
			continue
		}
		if strings.Contains(strings.ToLower(api.Api.Description), request.SearchTerm) {
			matchedOther = append(matchedOther, api)
			continue
		}

	}
	response.Apis = append(response.Apis, matchedName...)
	response.Apis = append(response.Apis, matchedEndpointName...)
	response.Apis = append(response.Apis, matchedOther...)
	return nil
}

// API returns a single named API entry
func (e *Explore) API(ctx context.Context, request *publicapi.APIRequest, response *publicapi.APIResponse) error {
	if len(request.Name) == 0 {
		return errors.BadRequest("explore.API", "missing name")
	}
	list, err := e.exploreList()
	if err != nil {
		return err
	}
	// check the cache
	for _, api := range list {
		if api.Api.Name == request.Name {
			response.Api = api
			return nil
		}
	}
	return errors.NotFound("explore.API", "not found")

}
