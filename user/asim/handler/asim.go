package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	asim "asim/proto/asim"
)

type Asim struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Asim) Call(ctx context.Context, req *asim.Request, rsp *asim.Response) error {
	log.Info("Received Asim.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Asim) Stream(ctx context.Context, req *asim.StreamingRequest, stream asim.Asim_StreamStream) error {
	log.Infof("Received Asim.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&asim.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Asim) PingPong(ctx context.Context, stream asim.Asim_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&asim.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
