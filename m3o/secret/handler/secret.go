package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	k8s "github.com/micro/go-micro/v2/util/kubernetes/client"

	pb "github.com/micro/services/m3o/secret/proto"
)

// NewSecret returns an initialised handler
func NewSecret(service micro.Service) *Secret {
	return &Secret{
		name: service.Name(),
		k8s:  k8s.NewClusterClient(),
	}
}

// Secret implements the proto secret service interface
type Secret struct {
	name string
	k8s  k8s.Client
}

// Create an image pull secret in kubernetes
func (s *Secret) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// ensure the token is valid before writing to k8s
	if err := s.validateGitHubToken(req.Token); err != nil {
		return err
	}

	// the secret structure required for img pull secrets
	secret := map[string]interface{}{
		"auths": map[string]interface{}{
			"docker.pkg.github.com": map[string]string{
				"auth": req.Token,
			},
		},
	}

	// encode the secret to json and then base64 encode
	bytes, _ := json.Marshal(secret)
	str := base64.StdEncoding.EncodeToString(bytes)

	// create the secret in k8s
	return s.k8s.Create(&k8s.Resource{
		Name: req.Namespace,
		Kind: "secret",
		Value: &k8s.Secret{
			Metadata: &k8s.Metadata{
				Name: req.Namespace,
			},
			Type: "kubernetes.io/dockerconfigjson",
			Data: map[string]string{
				".dockerconfigjson": str,
			},
		},
	})
}

func (s *Secret) validateGitHubToken(token string) error {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return errors.InternalServerError(s.name, "Unable to connect to GitHub API: %v", err)
	}
	if rsp.StatusCode != 200 {
		return errors.BadRequest(s.name, "Invalid credentials, status: %v", rsp.StatusCode)
	}

	return nil
}
