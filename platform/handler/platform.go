package handler

import (
	"context"

	"github.com/micro/micro/v3/service/client"

	pb "github.com/m3o/services/platform/proto"
	k8s "github.com/micro/go-micro/v3/util/kubernetes/client"
	"github.com/micro/micro/v3/service"
	rproto "github.com/micro/micro/v3/service/runtime/proto"
)

// Platform implements the platform service interface
type Platform struct {
	name    string
	client  k8s.Client
	runtime rproto.RuntimeService
}

// New returns an initialised platform handler
func New(service *service.Service) *Platform {
	return &Platform{
		name:    service.Name(),
		client:  k8s.NewClusterClient(),
		runtime: rproto.NewRuntimeService("runtime", client.DefaultClient),
	}
}

// CreateNamespace
func (k *Platform) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest, rsp *pb.CreateNamespaceResponse) error {
	_, err := k.runtime.CreateNamespace(ctx, &rproto.CreateNamespaceRequest{
		Namespace: req.Name,
	})
	return err
}

// DeleteNamespace
func (k *Platform) DeleteNamespace(ctx context.Context, req *pb.DeleteNamespaceRequest, rsp *pb.DeleteNamespaceResponse) error {
	_, err := k.runtime.DeleteNamespace(ctx, &rproto.DeleteNamespaceRequest{
		Namespace: req.Name,
	})
	return err
}
