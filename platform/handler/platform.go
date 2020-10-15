package handler

import (
	"context"

	pb "github.com/m3o/services/platform/proto"
	rproto "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
)

var (
	defaultNetworkPolicyName = "ingress"
	defaultAllowedLabels     = map[string]string{"owner": "micro"}
)

// Platform implements the platform service interface
type Platform struct {
	name    string
	runtime rproto.RuntimeService
}

// New returns an initialised platform handler
func New(service *service.Service) *Platform {
	return &Platform{
		name:    service.Name(),
		runtime: rproto.NewRuntimeService("runtime", client.DefaultClient),
	}
}

// CreateNamespace creates a new namespace, as well as a default network policy
func (k *Platform) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest, rsp *pb.CreateNamespaceResponse) error {

	// namespace
	if _, err := k.runtime.Create(ctx, &rproto.CreateRequest{
		Resource: &rproto.Resource{
			Namespace: &rproto.Namespace{
				Name: req.Name,
			},
		},
		Options: &rproto.CreateOptions{
			Namespace: req.Name,
		},
	}); err != nil {
		return err
	}

	// networkpolicy
	_, err := k.runtime.Create(ctx, &rproto.CreateRequest{
		Resource: &rproto.Resource{
			Networkpolicy: &rproto.NetworkPolicy{
				Allowedlabels: defaultAllowedLabels,
				Name:          defaultNetworkPolicyName,
				Namespace:     req.Name,
			},
		},
		Options: &rproto.CreateOptions{
			Namespace: req.Name,
		},
	})

	return err
}

// DeleteNamespace deletes a namespace, as well as its default network policy
func (k *Platform) DeleteNamespace(ctx context.Context, req *pb.DeleteNamespaceRequest, rsp *pb.DeleteNamespaceResponse) error {

	// networkpolicy (ignoring any error)
	k.runtime.Delete(ctx, &rproto.DeleteRequest{
		Resource: &rproto.Resource{
			Networkpolicy: &rproto.NetworkPolicy{
				Allowedlabels: defaultAllowedLabels,
				Name:          defaultNetworkPolicyName,
				Namespace:     req.Name,
			},
		},
		Options: &rproto.DeleteOptions{
			Namespace: req.Name,
		},
	})

	// namespace
	_, err := k.runtime.Delete(ctx, &rproto.DeleteRequest{
		Resource: &rproto.Resource{
			Namespace: &rproto.Namespace{
				Name: req.Name,
			},
		},
		Options: &rproto.DeleteOptions{
			Namespace: req.Name,
		},
	})
	return err
}
