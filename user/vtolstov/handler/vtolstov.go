package handler

import (
	"context"

	"github.com/micro/go-micro/v2/util/log"

	vtolstov "vtolstov/proto/vtolstov"
)

type Vtolstov struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Vtolstov) Call(ctx context.Context, req *vtolstov.Request, rsp *vtolstov.Response) error {
	log.Log("Received Vtolstov.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Vtolstov) Stream(ctx context.Context, req *vtolstov.StreamingRequest, stream vtolstov.Vtolstov_StreamStream) error {
	log.Logf("Received Vtolstov.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&vtolstov.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Vtolstov) PingPong(ctx context.Context, stream vtolstov.Vtolstov_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&vtolstov.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
