package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	foobar "foobar/proto/foobar"
)

type Foobar struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Foobar) Call(ctx context.Context, req *foobar.Request, rsp *foobar.Response) error {
	log.Info("Received Foobar.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Foobar) Stream(ctx context.Context, req *foobar.StreamingRequest, stream foobar.Foobar_StreamStream) error {
	log.Infof("Received Foobar.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&foobar.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Foobar) PingPong(ctx context.Context, stream foobar.Foobar_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&foobar.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
