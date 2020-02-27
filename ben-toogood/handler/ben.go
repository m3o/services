package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	ben "ben-toogood/proto/ben"
)

type Ben struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Ben) Call(ctx context.Context, req *ben.Request, rsp *ben.Response) error {
	log.Info("Received Ben.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Ben) Stream(ctx context.Context, req *ben.StreamingRequest, stream ben.Ben_StreamStream) error {
	log.Infof("Received Ben.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&ben.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Ben) PingPong(ctx context.Context, stream ben.Ben_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&ben.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
