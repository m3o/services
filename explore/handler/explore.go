package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	explore "explore/proto/explore"
)

type Explore struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Explore) Call(ctx context.Context, req *explore.Request, rsp *explore.Response) error {
	log.Info("Received Explore.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Explore) Stream(ctx context.Context, req *explore.StreamingRequest, stream explore.Explore_StreamStream) error {
	log.Infof("Received Explore.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&explore.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Explore) PingPong(ctx context.Context, stream explore.Explore_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&explore.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
