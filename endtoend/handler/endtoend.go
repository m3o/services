package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	endtoend "github.com/m3o/services/endtoend/proto"
)

type Endtoend struct{}

func (e *Endtoend) Mailin(ctx context.Context, request *endtoend.Request, response *endtoend.Response) error {
	log.Info("Received Endtoend.Mailin request")
	for hdr, pr := range request.Header {
		log.Infof("Received header %s %s %+v", hdr, pr.Key, pr.Values)
	}
	log.Infof("Request %+v", *request)
	response.StatusCode = 200
	return nil
}

func (e *Endtoend) Check(ctx context.Context, request *endtoend.Request, response *endtoend.Response) error {
	log.Info("Received Endtoend.Check request")
	response.StatusCode = 200
	return nil
}
