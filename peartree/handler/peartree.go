package handler

import (
	"context"

	"github.com/micro/go-micro/v2/util/log"

	peartree "peartree/proto/peartree"
)

type Peartree struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Peartree) Call(ctx context.Context, req *peartree.Request, rsp *peartree.Response) error {
	log.Log("Received Peartree.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Peartree) Stream(ctx context.Context, req *peartree.StreamingRequest, stream peartree.Peartree_StreamStream) error {
	log.Logf("Received Peartree.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&peartree.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Peartree) PingPong(ctx context.Context, stream peartree.Peartree_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&peartree.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
