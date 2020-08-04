package handler

import (
	"context"
	"encoding/json"
	"fmt"

	status "github.com/m3o/services/status/proto/status"
	api "github.com/micro/go-micro/v3/api/proto"
	proto "github.com/micro/go-micro/v3/debug/service/proto"
	"github.com/micro/micro/v3/service/client"
)

var (
	defaultServices = []string{
		"go.micro.api", // If this is down then this wouldn't even get routed to...
		"go.micro.auth",
		"go.micro.broker",
		"go.micro.config",
		"go.micro.debug",
		"go.micro.network",
		"go.micro.proxy",
		"go.micro.registry",
		"go.micro.runtime",
		"go.micro.store",
	}
)

type Status struct {
	monitoredServices []string
}

// NewStatusHandler returns a status handler configured to report the status of the given services
func NewStatusHandler(services []string) status.StatusHandler {
	svcs := defaultServices
	if len(services) > 0 {
		svcs = services
	}
	return &Status{monitoredServices: svcs}
}

// Call is called by the API as /status/call with post body {"name": "foo"}
func (e *Status) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	response := map[string]string{}
	overallOK := true

	// Are the services up?
	for _, serverName := range e.monitoredServices {
		req := client.NewRequest(serverName, "Debug.Health", &proto.HealthRequest{})
		rsp := &proto.HealthResponse{}

		err := client.Call(context.TODO(), req, rsp)
		status := "OK"
		if err != nil || rsp.Status != "ok" {
			status = "NOT_HEALTHY"
			if rsp != nil && rsp.Status != "" {
				status = rsp.Status
			}
			if err != nil {
				status = fmt.Sprintf("%s %s", status, err.Error())
			}
			overallOK = false
		}
		response[serverName] = status
	}

	b, _ := json.Marshal(response)
	statusCode := 200
	if !overallOK {
		statusCode = 500
	}
	rsp.StatusCode = int32(statusCode)
	rsp.Body = string(b)

	return nil
}
