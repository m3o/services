package handler

import (
	custpb "github.com/m3o/services/customers/proto"
)

type Endtoend struct {
	custSvc custpb.CustomersService
	email   string
}

type mailinMessage struct {
	Headers  map[string]interface{} `json:"headers"`
	Envelope map[string]interface{} `json:"envelope"`
	Plain    string                 `json:"plain"`
	Html     string                 `json:"html"`
}

type otp struct {
	Token string `json:"token"`
	Time  int64  `json:"time"`
}

type checkResult struct {
	Time   int64  `json:"time"`
	Passed bool   `json:"passed"`
	Error  string `json:"error"`
}

type MailinResponse struct{}
