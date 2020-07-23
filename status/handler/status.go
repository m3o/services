package handler

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2/client"

	proto "github.com/micro/go-micro/v2/debug/service/proto"

	api "github.com/micro/go-micro/v2/api/proto"
)

var (
	services = []string{
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

type Status struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

// Call is called by the API as /status/call with post body {"name": "foo"}
func (e *Status) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	response := map[string]string{}
	overallOK := true

	// Are the services up?
	for _, serverName := range services {
		req := client.NewRequest(serverName, "Debug.Health", &proto.HealthRequest{})
		rsp := &proto.HealthResponse{}

		err := client.Call(context.TODO(), req, rsp)
		status := "OK"
		if err != nil || rsp.Status != "ok" {
			status = "NOT_HEALTHY"
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
