package handler

import (
	"context"

	"github.com/micro/go-micro/v2/util/log"

	cruft "cruft/proto/cruft"
)

type Cruft struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Cruft) Call(ctx context.Context, req *cruft.Request, rsp *cruft.Response) error {
	log.Log("Received Cruft.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Cruft) Stream(ctx context.Context, req *cruft.StreamingRequest, stream cruft.Cruft_StreamStream) error {
	log.Logf("Received Cruft.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&cruft.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Cruft) PingPong(ctx context.Context, stream cruft.Cruft_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&cruft.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
