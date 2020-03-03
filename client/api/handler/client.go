package handler

import (
	"context"
	"encoding/json"

	pb "api/proto/client"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/metadata"
)

type Client struct {
	// micro client
	Client client.Client
}

// Client.Call is called by the API as /client/call with post body {"name": "foo"}
func (c *Client) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	log.Infof("Received Client.Call request service %s endpoint %s", req.Service, req.Endpoint)

	ct, ok := metadata.Get(ctx, "Content-Type")
	if !ok || len(ct) == 0 {
		ct = req.ContentType
	}

	// assume json until otherwise
	if ct != "application/json" {
		ct = "application/json"
	}

	// forward the request
	var payload json.RawMessage
	// if the extracted payload isn't empty lets use it
	if len(req.Body) > 0 {
		payload = json.RawMessage(req.Body)
	}

	// create request/response
	var response json.RawMessage

	// TODO: we will whitelist in auth
	request := c.Client.NewRequest(
		req.Service,
		req.Endpoint,
		&payload,
		client.WithContentType(ct),
	)

	// make the call
	if err := c.Client.Call(ctx, request, &response); err != nil {
		return err
	}

	// marshall response
	// TODO implement errors
	rsp.Body, _ = response.MarshalJSON()
	return nil
}
