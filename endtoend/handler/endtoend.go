package handler

import (
	"context"
	"encoding/json"

	log "github.com/micro/micro/v3/service/logger"
)

type Endtoend struct{}

type MailinResponse struct{}

type mailinMessage struct {
	Headers  map[string]interface{} `json:"headers"`
	Envelope map[string]interface{} `json:"envelope"`
	Plain    string                 `json:"plain"`
	Html     string                 `json:"html"`
}

func (e *Endtoend) Mailin(ctx context.Context, req *json.RawMessage, rsp *MailinResponse) error {
	log.Info("Received Endtoend.Mailin request %d", len(*req))
	var inbound mailinMessage

	if err := json.Unmarshal(*req, &inbound); err != nil {
		log.Errorf("Error unmarshalling request %s", err)
	}

	log.Infof("Request %+v", inbound)
	return nil
}

//func (e *Endtoend) Check(ctx context.Context, request *endtoend.Request, response *endtoend.Response) error {
//	log.Info("Received Endtoend.Check request")
//	response.StatusCode = 200
//	return nil
//}
