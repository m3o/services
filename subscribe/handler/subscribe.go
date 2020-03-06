package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	subscribe "subscribe/proto/subscribe"
)

type Subscribe struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Subscribe) Call(ctx context.Context, req *subscribe.Request, rsp *subscribe.Response) error {
	log.Info("Received Subscribe.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Subscribe) Stream(ctx context.Context, req *subscribe.StreamingRequest, stream subscribe.Subscribe_StreamStream) error {
	log.Infof("Received Subscribe.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&subscribe.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Subscribe) PingPong(ctx context.Context, stream subscribe.Subscribe_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&subscribe.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
