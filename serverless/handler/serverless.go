package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	serverless "serverless/proto/serverless"
)

type Serverless struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Serverless) Call(ctx context.Context, req *serverless.Request, rsp *serverless.Response) error {
	log.Info("Received Serverless.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Serverless) Stream(ctx context.Context, req *serverless.StreamingRequest, stream serverless.Serverless_StreamStream) error {
	log.Infof("Received Serverless.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&serverless.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Serverless) PingPong(ctx context.Context, stream serverless.Serverless_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&serverless.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
