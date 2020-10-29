package handler

import (
	"context"

	pb "github.com/m3o/services/platform/proto"
	rproto "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
)

var (
	defaultNetworkPolicyName = "ingress"
	defaultResourceQuotaName = "quota"
	defaultAllowedLabels     = map[string]string{"owner": "micro"}
	defaultResourceLimits    = &rproto.Resources{}
	defaultResourceRequests  = &rproto.Resources{}
)

// Platform implements the platform service interface
type Platform struct {
	name    string
	runtime rproto.RuntimeService
}

// New returns an initialised platform handler
func New(service *service.Service) *Platform {

	if val, err := config.Get("micro.platform.resource_limits.cpu"); err != nil {
		defaultResourceLimits.CPU = int32(val.Int(8000))
	}

	if val, err := config.Get("micro.platform.resource_limits.disk"); err != nil {
		defaultResourceLimits.EphemeralStorage = int32(val.Int(8000))
	}

	if val, err := config.Get("micro.platform.resource_limits.memory"); err != nil {
		defaultResourceLimits.Memory = int32(val.Int(8000))
	}

	if val, err := config.Get("micro.platform.resource_requests.cpu"); err != nil {
		defaultResourceRequests.CPU = int32(val.Int(8000))
	}

	if val, err := config.Get("micro.platform.resource_requests.disk"); err != nil {
		defaultResourceRequests.EphemeralStorage = int32(val.Int(8000))
	}

	if val, err := config.Get("micro.platform.resource_requests.memory"); err != nil {
		defaultResourceRequests.Memory = int32(val.Int(8000))
	}

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
	if _, err := k.runtime.Create(ctx, &rproto.CreateRequest{
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
	}); err != nil {
		return err
	}

	// resourcequota
	_, err := k.runtime.Create(ctx, &rproto.CreateRequest{
		Resource: &rproto.Resource{
			Resourcequota: &rproto.ResourceQuota{
				Name:      defaultResourceQuotaName,
				Namespace: req.Name,
				Requests:  defaultResourceRequests,
				Limits:    defaultResourceLimits,
			},
		},
		Options: &rproto.CreateOptions{
			Namespace: req.Name,
		},
	})

	return err
}

// DeleteNamespace deletes a namespace, as well as anything inside it (services, network policies, etc)
func (k *Platform) DeleteNamespace(ctx context.Context, req *pb.DeleteNamespaceRequest, rsp *pb.DeleteNamespaceResponse) error {
	// kill all the services
	rrsp, err := k.runtime.Read(ctx, &rproto.ReadRequest{Options: &rproto.ReadOptions{Namespace: req.Name}})
	if err != nil {
		return err
	}
	for _, s := range rrsp.Services {
		k.runtime.Delete(ctx, &rproto.DeleteRequest{
			Resource: &rproto.Resource{
				Service: s,
			},
			Options: &rproto.DeleteOptions{Namespace: req.Name},
		})

	}

	// resourcequota (ignoring any error)
	k.runtime.Delete(ctx, &rproto.DeleteRequest{
		Resource: &rproto.Resource{
			Resourcequota: &rproto.ResourceQuota{
				Name:      defaultResourceQuotaName,
				Namespace: req.Name,
			},
		},
		Options: &rproto.DeleteOptions{
			Namespace: req.Name,
		},
	})

	// networkpolicy (ignoring any error)
	k.runtime.Delete(ctx, &rproto.DeleteRequest{
		Resource: &rproto.Resource{
			Networkpolicy: &rproto.NetworkPolicy{
				Name:      defaultNetworkPolicyName,
				Namespace: req.Name,
			},
		},
		Options: &rproto.DeleteOptions{
			Namespace: req.Name,
		},
	})

	// namespace
	_, err = k.runtime.Delete(ctx, &rproto.DeleteRequest{
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
