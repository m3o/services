package handler

import (
	"context"
	"encoding/json"
	"io"

	pb "github.com/m3o/services/api/client/proto/client"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
	log "github.com/micro/micro/v3/service/logger"
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

	// TODO: we will whitelist in auth
	request := c.Client.NewRequest(
		req.Service,
		req.Endpoint,
		&payload,
		client.WithContentType(ct),
	)

	// create request/response
	var response json.RawMessage

	// make the call
	if err := c.Client.Call(ctx, request, &response); err != nil {
		return err
	}

	// marshall response
	// TODO implement errors
	rsp.Body, _ = response.MarshalJSON()
	return nil
}

// Client.Stream is a bidirectional stream called by the API at /client/stream
func (c *Client) Stream(ctx context.Context, stream pb.Client_StreamStream) error {
	// get the first request
	req, err := stream.Recv()
	if err != nil {
		log.Errorf("Failed to get Client.Stream request: %v", err)
		return err
	}

	log.Infof("Received Client.Stream request service %s endpoint %s", req.Service, req.Endpoint)

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

	// TODO: we will whitelist in auth
	request := c.Client.NewRequest(
		req.Service,
		req.Endpoint,
		&payload,
		client.WithContentType(ct),
	)

	// make the call
	serviceStream, err := c.Client.Stream(ctx, request)
	if err != nil {
		return err
	}

	// send the first request to backend
	if err := serviceStream.Send(payload); err != nil {
		return err
	}

	go writeLoop(stream, serviceStream)

	// get raw response
	resp := serviceStream.Response()

	// create server response write loop
	for {
		// read backend response body
		body, err := resp.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		var payload json.RawMessage

		if len(body) > 0 {
			payload = json.RawMessage(body)
		}

		// stream to frontend
		rsp := new(pb.Response)
		// marshall response
		// TODO implement errors
		rsp.Body = payload

		if err := stream.Send(rsp); err != nil {
			return err
		}
	}

	return nil
}

func writeLoop(frontend pb.Client_StreamStream, backend client.Stream) {
	defer backend.Close()

	for {
		var payload json.RawMessage

		// get request from frontend
		rsp, err := frontend.Recv()
		if err != nil {
			return
		}

		// check if there's data
		if len(rsp.Body) > 0 {
			payload = json.RawMessage(rsp.Body)
		}

		// send request to backend
		if err := backend.Send(payload); err != nil {
			return
		}
	}
}
