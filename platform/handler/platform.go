package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"

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

// CreateImagePullSecret in k8s
func (k *Platform) CreateImagePullSecret(ctx context.Context, req *pb.CreateImagePullSecretRequest, rsp *pb.CreateImagePullSecretResponse) error {
	// the secret structure required for img pull secrets
	secret := map[string]interface{}{
		"auths": map[string]interface{}{
			"docker.pkg.github.com": map[string]string{
				"Username": "username",
				"Password": req.Token,
			},
		},
	}

	// encode the secret to json and then base64 encode
	bytes, _ := json.Marshal(secret)
	str := base64.StdEncoding.EncodeToString(bytes)

	// create the secret in k8s
	return k.client.Create(&k8s.Resource{
		Name: req.Name,
		Kind: "secret",
		Value: &k8s.Secret{
			Metadata: &k8s.Metadata{
				Name: req.Name,
			},
			Type: "kubernetes.io/dockerconfigjson",
			Data: map[string]string{
				".dockerconfigjson": str,
			},
		},
	}, k8s.CreateNamespace(req.Namespace))
}

// DeleteImagePullSecret in k8s
func (k *Platform) DeleteImagePullSecret(ctx context.Context, req *pb.DeleteImagePullSecretRequest, rsp *pb.DeleteImagePullSecretResponse) error {
	return k.client.Delete(&k8s.Resource{
		Name: req.Name,
		Kind: "secret",
		Value: &k8s.Secret{
			Metadata: &k8s.Metadata{
				Name: req.Name,
			},
		},
	}, k8s.DeleteNamespace(req.Namespace))
}

// CreateServiceAccount in k8s. Note, the service accounts are always named the same
// as the namespace so there is no name attribute in the request.
func (k *Platform) CreateServiceAccount(ctx context.Context, req *pb.CreateServiceAccountRequest, rsp *pb.CreateServiceAccountResponse) error {
	var secrets []k8s.ImagePullSecret
	if req.ImagePullSecrets != nil {
		for _, s := range req.ImagePullSecrets {
			secrets = append(secrets, k8s.ImagePullSecret{
				Name: s,
			})
		}
	}

	return k.client.Create(&k8s.Resource{
		Name: req.Namespace,
		Kind: "serviceaccount",
		Value: &k8s.ServiceAccount{
			Metadata: &k8s.Metadata{
				Name: req.Namespace,
			},
			ImagePullSecrets: secrets,
		},
	}, k8s.CreateNamespace(req.Namespace))
}

// DeleteServiceAccount in k8s. Note, the service accounts are always named the same
// as the namespace so there is no name attribute in the request.
func (k *Platform) DeleteServiceAccount(ctx context.Context, req *pb.DeleteServiceAccountRequest, rsp *pb.DeleteServiceAccountResponse) error {
	return k.client.Delete(&k8s.Resource{
		Name: req.Namespace,
		Kind: "serviceaccount",
		Value: &k8s.ServiceAccount{
			Metadata: &k8s.Metadata{
				Name: req.Namespace,
			},
		},
	}, k8s.DeleteNamespace(req.Namespace))
}
