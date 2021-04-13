package handler

import (
	"context"

	nproto "github.com/m3o/services/namespaces/proto"
	cproto "github.com/m3o/services/customers/proto"
	pb "github.com/m3o/services/platform/proto"
	rproto "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
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
	customer     cproto.CustomersService
	namespace    nproto.NamespacesService
}

// New returns an initialised platform handler
func New(service *service.Service) *Platform {

	val, _ := config.Get("micro.platform.resource_limits.cpu")
	defaultResourceLimits.CPU = int32(val.Int(8000))

	val, _ = config.Get("micro.platform.resource_limits.disk")
	defaultResourceLimits.EphemeralStorage = int32(val.Int(8000))

	val, _ = config.Get("micro.platform.resource_limits.memory")
	defaultResourceLimits.Memory = int32(val.Int(8000))

	val, _ = config.Get("micro.platform.resource_requests.cpu")
	defaultResourceRequests.CPU = int32(val.Int(8000))

	val, _ = config.Get("micro.platform.resource_requests.disk")
	defaultResourceRequests.EphemeralStorage = int32(val.Int(8000))

	val, _ = config.Get("micro.platform.resource_requests.memory")
	defaultResourceRequests.Memory = int32(val.Int(8000))

	return &Platform{
		name:    service.Name(),
		runtime: rproto.NewRuntimeService("runtime", client.DefaultClient),
		customer: cproto.NewCustomersService("customers", client.DefaultClient),
		namespace: nproto.NewNamespacesService("namespaces", client.DefaultClient),
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

func (p *Platform) LoginUser(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	if len(req.Username) == 0 {
		return errors.BadRequest("platform.loginUser", "Require username")
	}
	if len(req.Password) == 0 {
		return errors.BadRequest("platform.loginUser", "Require password")
	}

	// if the namespace is specified then just login
	// otherwise we lookup the namespace
	if len(req.Namespace) == 0 {
		namespaces, err := p.lookupNamespace(ctx, req.Username)
		if err != nil {
			return err
		}
		// might have multiple namespaces
		// TODO: check if more than one
		req.Namespace = namespaces[0].Id
	}

	// get the token from the auth service
	token, err := auth.Token(
		auth.WithCredentials(req.Username, req.Password),
		auth.WithTokenIssuer(req.Namespace),
	)
	if err != nil {
		return errors.InternalServerError("platform.loginUser", "Failed to authenticate: %v", err)
	}

	// set the response
	rsp.AccessToken = token.AccessToken
	rsp.RefreshToken = token.RefreshToken
	rsp.Created = token.Created.Unix()
	rsp.Expiry = token.Expiry.Unix()

	return nil
}

func (p *Platform) lookupNamespace(ctx context.Context, email string) ([]*nproto.Namespace, error) {
	custResp, err := p.customer.Read(ctx, &cproto.ReadRequest{Email: email}, client.WithAuthToken())
	if err != nil {
		merr := errors.FromError(err)
		if merr.Code == 404 { // not found
			return nil, errors.NotFound("platform.loginUser", "Could not find an account with that email address")
		}
		return nil, errors.InternalServerError("platform.loginUser", "Error looking up account")
	}

	listRsp, err := p.namespace.List(ctx, &nproto.ListRequest{
		User: custResp.Customer.Id,
	}, client.WithAuthToken())
	if err != nil {
		return nil, errors.InternalServerError("platform.loginUser", "Error calling namespace service: %v", err)
	}

	if len(listRsp.Namespaces) == 0 {
		return nil, errors.BadRequest("platform.loginUser", "We don't recognize this account")
	}

	return listRsp.Namespaces, nil
}
