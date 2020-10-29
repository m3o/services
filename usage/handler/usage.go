package handler

import (
	"context"

	usage "github.com/m3o/services/usage/proto"

	nsproto "github.com/m3o/services/namespaces/proto"
	pb "github.com/micro/micro/v3/proto/auth"
	rproto "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/client"
	log "github.com/micro/micro/v3/service/logger"
)

const (
	defaultNamespace = "micro"
)

type Usage struct {
	ns      nsproto.NamespacesService
	as      pb.AccountsService
	runtime rproto.RuntimeService
}

func NewUsage(ns nsproto.NamespacesService, as pb.AccountsService, runtime rproto.RuntimeService) *Usage {
	u := &Usage{
		ns:      ns,
		as:      as,
		runtime: runtime,
	}
	return u
}

// Read account history by namespace, or lists latest values for each namespace if history is not provided.
func (e *Usage) Read(ctx context.Context, req *usage.ReadRequest, rsp *usage.ReadResponse) error {
	log.Infof("Received Usage.Read request, reading namespace '%v'", req.Namespace)

	u, err := e.usageForNamespace(req.Namespace)
	if err != nil {
		return err
	}
	rsp.Accounts = []*usage.Account{
		{
			Namespace: req.Namespace,
			Users:     u.Users,
			Services:  u.Services,
		},
	}
	return nil
}

type usg struct {
	Users     int64
	Services  int64
	Namespace string
}

func (e *Usage) usageForNamespace(namespace string) (*usg, error) {
	arsp, err := e.as.List(context.TODO(), &pb.ListAccountsRequest{
		Options: &pb.Options{
			Namespace: namespace,
		},
	}, client.WithAuthToken())
	if err != nil {
		return nil, err
	}
	userCount := 0
	for _, account := range arsp.Accounts {
		if account.Type == "user" {
			userCount++
		}
	}
	rrsp, err := e.runtime.Read(context.TODO(), &rproto.ReadRequest{
		Options: &rproto.ReadOptions{
			Namespace: namespace,
		},
	}, client.WithAuthToken())
	if err != nil {
		return nil, err
	}
	serviceCount := len(rrsp.Services)
	return &usg{
		Users:     int64(userCount),
		Services:  int64(serviceCount),
		Namespace: namespace,
	}, nil
}
