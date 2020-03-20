package events

import (
	"net/http"

	ghttp "github.com/micro/go-micro/v2/util/http"
	"github.com/micro/go-micro/v2/web"
	platform "github.com/micro/services/platform/service/proto"
	utils "github.com/micro/services/platform/web/util"
)

// RegisterHandlers adds the GitHub oauth handlers to the servie
func RegisterHandlers(srv web.Service) error {
	h := Handler{
		platform: platform.NewPlatformService("go.micro.service.platform", srv.Options().Service.Client()),
	}

	srv.HandleFunc("/v1/events", h.EventsHandler)
	return nil
}

// Handler encapsulates the events handlers
type Handler struct {
	platform platform.PlatformService
}

// EventsHandler returns all recent events, or the events scoped to the service
// requested using the "service" URL query param.
func (h *Handler) EventsHandler(w http.ResponseWriter, req *http.Request) {
	utils.SetupResponse(&w, req)
	if req.Method == "OPTIONS" {
		return
	}
	ctx := ghttp.RequestToContext(req)
	// construct the request to the platform service
	var eReq platform.ListEventsRequest
	if srvName := req.URL.Query().Get("service"); len(srvName) > 0 {
		eReq.Service = &platform.Service{Name: srvName, Type: "service"}
	}

	// execute the request and handle any errors
	resp, err := h.platform.ListEvents(ctx, &eReq)
	if err != nil {
		utils.Write500(w, err)
		return
	}
	if resp.Events == nil {
		resp.Events = []*platform.Event{}
	}

	// return the events
	utils.WriteJSON(w, resp.Events)
}
