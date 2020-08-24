package handler

import (
	"context"

	pb "github.com/m3o/services/platform/proto"
	k8s "github.com/micro/go-micro/v3/util/kubernetes/client"
	"github.com/micro/micro/v3/service"
)

// Platform implements the platform service interface
type Platform struct {
	name   string
	client k8s.Client
}

// New returns an initialised platform handler
func New(service *service.Service) *Platform {
	return &Platform{
		name:   service.Name(),
		client: k8s.NewClusterClient(),
	}
}

// CreateNamespace in k8s
func (k *Platform) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest, rsp *pb.CreateNamespaceResponse) error {
	return k.client.Create(&k8s.Resource{
		Kind: "namespace",
		Value: k8s.Namespace{
			Metadata: &k8s.Metadata{
				Name: req.Name,
			},
		},
	})
}

// DeleteNamespace in k8s
func (k *Platform) DeleteNamespace(ctx context.Context, req *pb.DeleteNamespaceRequest, rsp *pb.DeleteNamespaceResponse) error {
	return k.client.Delete(&k8s.Resource{
		Kind: "namespace",
		Name: req.Name,
		Value: k8s.Namespace{
			Metadata: &k8s.Metadata{
				Name: req.Name,
			},
		},
	})
}
