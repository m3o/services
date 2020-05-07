package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"

	kubernetes "github.com/micro/services/kubernetes/service/proto"
	pb "github.com/micro/services/m3o/api/proto"
	environments "github.com/micro/services/projects/environments/proto"
	projects "github.com/micro/services/projects/service/proto"
	users "github.com/micro/services/users/service/proto"
)

// NewProjects returns an initialised projects handler
func NewProjects(service micro.Service) *Projects {
	return &Projects{
		name:         service.Name(),
		auth:         service.Options().Auth,
		users:        users.NewUsersService("go.micro.service.users", service.Client()),
		projects:     projects.NewProjectsService("go.micro.service.projects", service.Client()),
		kubernetes:   kubernetes.NewKubernetesService("go.micro.service.kubernetes", service.Client()),
		environments: environments.NewEnvironmentsService("go.micro.service.projects.environments", service.Client()),
	}
}

// Projects implments the M3O project service proto
type Projects struct {
	name         string
	auth         auth.Auth
	users        users.UsersService
	projects     projects.ProjectsService
	kubernetes   kubernetes.KubernetesService
	environments environments.EnvironmentsService
}

// CreateProject and the underlying infra
func (p *Projects) CreateProject(ctx context.Context, req *pb.CreateProjectRequest, rsp *pb.CreateProjectResponse) error {
	// Validate the request
	if req.Project == nil {
		return errors.BadRequest(p.name, "Missing project")
	}

	// Identify the user
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return err
	}

	// Validate the user has access to the github repo
	repos, err := p.listGitHubRepos(req.GithubToken)
	if err != nil {
		return err
	}
	var isMemberOfRepo bool
	for _, r := range repos {
		if r.Name == req.Project.Repository {
			isMemberOfRepo = true
			break
		}
	}
	if !isMemberOfRepo {
		return errors.BadRequest(p.name, "Must be a member of the repository")
	}

	// create the project
	cRsp, err := p.projects.Create(ctx, &projects.CreateRequest{
		Project: &projects.Project{
			Name:        strings.ToLower(req.Project.Name),
			Description: req.Project.Description,
			Repository:  req.Project.Repository,
		},
	})
	if err != nil {
		return err
	}

	// add the user as an owner
	_, err = p.projects.AddMember(ctx, &projects.AddMemberRequest{
		Role:      projects.Role_Owner,
		ProjectId: cRsp.Project.Id,
		Member: &projects.Member{
			Type: "user",
			Id:   userID,
		},
	})
	if err != nil {
		return err
	}

	// serialize the project
	rsp.Project = serializeProject(cRsp.Project)

	// generate the auth account for the webhooks
	rsp.ClientId, rsp.ClientSecret, err = p.generateCreds(cRsp.Project.Id)
	return nil
}

// UpdateProject metadata
func (p *Projects) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest, rsp *pb.UpdateProjectResponse) error {
	// find the project
	proj, err := p.findProject(ctx, req.Id)
	if err != nil {
		return err
	}

	// assign the update attributes
	proj.Name = req.Name
	proj.Description = req.Description

	// update the project
	_, err = p.projects.Update(ctx, &projects.UpdateRequest{Project: proj})
	return err
}

// ListProjects the user has access to
func (p *Projects) ListProjects(ctx context.Context, req *pb.ListProjectsRequest, rsp *pb.ListProjectsResponse) error {
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return err
	}

	tRsp, err := p.projects.ListMemberships(ctx, &projects.ListMembershipsRequest{
		Member: &projects.Member{Type: "user", Id: userID},
	})
	if err != nil {
		return err
	}

	rsp.Projects = make([]*pb.Project, 0, len(tRsp.Projects))
	for _, pr := range tRsp.Projects {
		proj := serializeProject(pr)

		eRsp, err := p.environments.Read(ctx, &environments.ReadRequest{ProjectId: pr.Id})
		if err == nil {
			proj.Environments = make([]*pb.Environment, 0, len(eRsp.Environments))
			for _, e := range eRsp.Environments {
				proj.Environments = append(proj.Environments, serializeEnvironment(e))
			}
		}

		rsp.Projects = append(rsp.Projects, proj)
	}

	return nil
}

// ValidateProjectName validates a project name to ensure it is unique
func (p *Projects) ValidateProjectName(ctx context.Context, req *pb.ValidateProjectNameRequest, rsp *pb.ValidateProjectNameResponse) error {
	_, err := p.projects.Read(ctx, &projects.ReadRequest{Name: req.Name})
	if err == nil {
		return errors.BadRequest(p.name, "Name has already been taken")
	}
	return nil
}

// ValidateEnvironmentName validates a Environment name to ensure it is unique
func (p *Projects) ValidateEnvironmentName(ctx context.Context, req *pb.ValidateEnvironmentNameRequest, rsp *pb.ValidateEnvironmentNameResponse) error {
	eRsp, err := p.environments.Read(ctx, &environments.ReadRequest{Id: req.ProjectId})
	if err != nil {
		return err
	}

	for _, env := range eRsp.Environments {
		if env.Name == req.Name {
			return errors.BadRequest(p.name, "Name has already been taken")
		}
	}

	return nil
}

// ValidateGithubToken takes a GitHub personal token and returns the repos it has access to
func (p *Projects) ValidateGithubToken(ctx context.Context, req *pb.ValidateGithubTokenRequest, rsp *pb.ValidateGithubTokenResponse) error {
	repos, err := p.listGitHubRepos(req.Token)
	if err != nil {
		return err
	}
	rsp.Repos = repos
	return nil
}

// WebhookAPIKey generates an auth account token which can be used to authenticate against the webhook api
func (p *Projects) WebhookAPIKey(ctx context.Context, req *pb.WebhookAPIKeyRequest, rsp *pb.WebhookAPIKeyResponse) error {
	// find the project
	proj, err := p.findProject(ctx, req.ProjectId)
	if err != nil {
		return err
	}

	// generate the auth account
	rsp.ClientId, rsp.ClientSecret, err = p.generateCreds(proj.Id)
	return err
}

// CreateEnvironment for a given project
func (p *Projects) CreateEnvironment(ctx context.Context, req *pb.CreateEnvironmentRequest, rsp *pb.CreateEnvironmentResponse) error {
	// validate the request
	if req.Environment == nil {
		return errors.BadRequest(p.name, "Missing environment")
	}
	if len(req.ProjectId) == 0 {
		return errors.BadRequest(p.name, "Missing project id")
	}

	// ensure the user has access to the project
	if _, err := p.findProject(ctx, req.ProjectId); err != nil {
		return errors.Forbidden(p.name, "Unable to access project")
	}

	// create the environment
	env := &environments.Environment{
		ProjectId:   req.ProjectId,
		Name:        strings.ToLower(req.Environment.Name),
		Description: req.Environment.Description,
	}
	eRsp, err := p.environments.Create(ctx, &environments.CreateRequest{Environment: env})
	if err != nil {
		return errors.BadRequest(p.name, "Unable to create project: %v", err.Error())
	}

	// create the k8s namespace etc
	if _, err := p.kubernetes.CreateNamespace(ctx, &kubernetes.CreateNamespaceRequest{Name: eRsp.Environment.Namespace}); err != nil {
		p.environments.Delete(ctx, &environments.DeleteRequest{Id: eRsp.Environment.Id})
		return errors.BadRequest(p.name, "Unable to create k8s namespace: %v", err.Error())
	}

	// TODO: Load the projects secret (the GH token) and create an image pull secret in the above namespacce
	rsp.Environment = serializeEnvironment(eRsp.Environment)
	return nil
}

// UpdateEnvironment metadata
func (p *Projects) UpdateEnvironment(ctx context.Context, req *pb.UpdateEnvironmentRequest, rsp *pb.UpdateEnvironmentResponse) error {
	// validate the request
	if req.Environment == nil {
		return errors.BadRequest(p.name, "Missing environment")
	}
	if len(req.Environment.Id) == 0 {
		return errors.BadRequest(p.name, "Missing environment id")
	}

	// lookup the environment
	rRsp, err := p.environments.Read(ctx, &environments.ReadRequest{Id: req.Environment.Id})
	if err != nil {
		return err
	}
	env := rRsp.Environment

	// ensure the user has access to the project
	if _, err := p.findProject(ctx, env.ProjectId); err != nil {
		return err
	}

	// assign the update attributes
	env.Description = req.Environment.Description

	// update the environment
	_, err = p.environments.Update(ctx, &environments.UpdateRequest{Environment: env})
	return err
}

// DeleteEnvironment and the underlying infra
func (p *Projects) DeleteEnvironment(ctx context.Context, req *pb.DeleteEnvironmentRequest, rsp *pb.DeleteEnvironmentRequest) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest(p.name, "Missing id")
	}

	// lookup the environment
	rRsp, err := p.environments.Read(ctx, &environments.ReadRequest{Id: req.Id})
	if err != nil {
		return err
	}
	env := rRsp.Environment

	// ensure the user has access to the project
	if _, err := p.findProject(ctx, env.ProjectId); err != nil {
		return err
	}

	// delete the k8s namespace
	if _, err = p.kubernetes.DeleteNamespace(ctx, &kubernetes.DeleteNamespaceRequest{Name: env.Namespace}); err != nil {
		return err
	}

	// delete the environment
	_, err = p.environments.Delete(ctx, &environments.DeleteRequest{Id: env.Id})
	return err
}

func (p *Projects) generateCreds(projectID string) (string, string, error) {
	id := fmt.Sprintf("%v-webhook-%v", projectID, time.Now().Unix())
	md := map[string]string{"project-id": projectID}

	acc, err := p.auth.Generate(id, auth.WithRoles("webhook"), auth.WithMetadata(md))
	if err != nil {
		return "", "", err
	}

	return acc.ID, acc.Secret, nil
}

func (p *Projects) userIDFromContext(ctx context.Context) (string, error) {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return "", errors.Unauthorized(p.name, "Account Required")
	}

	uRsp, err := p.users.Read(ctx, &users.ReadRequest{Email: acc.ID})
	if err != nil {
		return "", errors.InternalServerError(p.name, "Auth error: %v", err)
	}

	return uRsp.User.Id, nil
}

func serializeProject(p *projects.Project) *pb.Project {
	return &pb.Project{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Repository:  p.Repository,
	}
}

func serializeEnvironment(e *environments.Environment) *pb.Environment {
	return &pb.Environment{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
	}
}

func (p *Projects) listGitHubRepos(token string) ([]*pb.Repository, error) {
	r, _ := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "application/vnd.github.nebula-preview+json")

	res, err := new(http.Client).Do(r)
	if err != nil {
		return nil, errors.InternalServerError(p.name, "Unable to connect to the GitHub API: %v", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.BadRequest(p.name, "Invalid GitHub token")
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.InternalServerError(p.name, "Unexpected status returned from the GitHub API: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.InternalServerError(p.name, "Invalid response returned from the GitHub API: %v", err)
	}

	var repos []struct {
		Name    string `json:"full_name"`
		Private bool   `json:"private"`
	}
	if err := json.Unmarshal(bytes, &repos); err != nil {
		return nil, errors.InternalServerError(p.name, "Invalid response returned from the GitHub API: %v", err)
	}

	repoos := make([]*pb.Repository, 0, len(repos))
	for _, r := range repos {
		repoos = append(repoos, &pb.Repository{Name: strings.ToLower(r.Name), Private: r.Private})
	}

	return repoos, nil
}

func (p *Projects) findProject(ctx context.Context, id string) (*projects.Project, error) {
	// Identify the user
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// get the projects the user belongs to
	mRsp, err := p.projects.ListMemberships(ctx, &projects.ListMembershipsRequest{
		Member: &projects.Member{Type: "user", Id: userID},
	})
	if err != nil {
		return nil, err
	}

	// check for membership
	var isMember bool
	for _, t := range mRsp.Projects {
		if t.Id == id {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, errors.Forbidden(p.name, "Not a member of this team")
	}

	// lookup the project
	rRsp, err := p.projects.Read(ctx, &projects.ReadRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return rRsp.GetProject(), nil
}
