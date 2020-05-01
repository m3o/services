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

	pb "github.com/micro/services/m3o/api/proto"
	project "github.com/micro/services/project/service/proto"
	users "github.com/micro/services/users/service/proto"
)

// NewProject returns an initialised project handler
func NewProject(service micro.Service) *Project {
	return &Project{
		name:    service.Name(),
		auth:    service.Options().Auth,
		users:   users.NewUsersService("go.micro.service.users", service.Client()),
		project: project.NewProjectService("go.micro.service.project", service.Client()),
	}
}

// Project implments the M3O project service proto
type Project struct {
	name    string
	auth    auth.Auth
	users   users.UsersService
	project project.ProjectService
}

// Create a new project
func (p *Project) Create(ctx context.Context, req *pb.CreateProjectRequest, rsp *pb.CreateProjectResponse) error {
	// Validate the request
	if req.Project == nil {
		return errors.BadRequest(p.name, "Missing project")
	}

	// Identify the user
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return err
	}

	// verify the user has access to the github repo
	repos, err := p.listGitHubRepos(req.GithubToken)
	if err != nil {
		return err
	}
	var isMemberOfRepo bool
	for _, r := range repos {
		if r == req.Project.Repository {
			isMemberOfRepo = true
			break
		}
	}
	if !isMemberOfRepo {
		return errors.BadRequest(p.name, "Must be a member of the repository")
	}

	// create the project
	cRsp, err := p.project.Create(ctx, &project.CreateRequest{
		Project: &project.Project{
			Name:        req.Project.Name,
			Description: req.Project.Description,
			Repository:  req.Project.Repository,
		},
	})
	if err != nil {
		return err
	}

	// add the user as a member
	_, err = p.project.AddMember(ctx, &project.AddMemberRequest{
		MemberId: userID, ProjectId: cRsp.Project.Id,
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

// Update a project
func (p *Project) Update(ctx context.Context, req *pb.UpdateProjectRequest, rsp *pb.UpdateProjectResponse) error {
	// find the project
	proj, err := p.findProject(ctx, req.Id)
	if err != nil {
		return err
	}

	// assign the update attributes
	proj.Name = req.Name
	proj.Description = req.Description

	// update the project
	_, err = p.project.Update(ctx, &project.UpdateRequest{Project: proj})
	return err
}

// List all the projects the user has access to
func (p *Project) List(ctx context.Context, req *pb.ListProjectsRequest, rsp *pb.ListProjectsResponse) error {
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return err
	}

	tRsp, err := p.project.ListMemberships(ctx, &project.ListMembershipsRequest{MemberId: userID})
	if err != nil {
		return err
	}

	rsp.Projects = make([]*pb.Project, 0, len(tRsp.Projects))
	for _, t := range tRsp.Projects {
		rsp.Projects = append(rsp.Projects, serializeProject(t))
	}

	return nil
}

// VerifyGithubToken takes a GitHub personal token and returns the repos it has access to
func (p *Project) VerifyGithubToken(ctx context.Context, req *pb.VerifyGithubTokenRequest, rsp *pb.VerifyGithubTokenResponse) error {
	repos, err := p.listGitHubRepos(req.Token)
	if err != nil {
		return err
	}
	rsp.Repos = repos
	return nil
}

// WebhookAPIKey generates an auth account token which can be used to authenticate against the webhook api
func (p *Project) WebhookAPIKey(ctx context.Context, req *pb.WebhookAPIKeyRequest, rsp *pb.WebhookAPIKeyResponse) error {
	// find the project
	proj, err := p.findProject(ctx, req.ProjectId)
	if err != nil {
		return err
	}

	// generate the auth account
	rsp.ClientId, rsp.ClientSecret, err = p.generateCreds(proj.Id)
	return err
}

func (p *Project) generateCreds(projectID string) (string, string, error) {
	id := fmt.Sprintf("%v-webhook-%v", projectID, time.Now().Unix())
	md := map[string]string{"project-id": projectID}

	acc, err := p.auth.Generate(id, auth.WithRoles("webhook"), auth.WithMetadata(md))
	if err != nil {
		return "", "", err
	}

	return acc.ID, acc.Secret, nil
}

func (p *Project) userIDFromContext(ctx context.Context) (string, error) {
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

func serializeProject(p *project.Project) *pb.Project {
	return &pb.Project{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Repository:  p.Repository,
	}
}

func (p *Project) listGitHubRepos(token string) ([]string, error) {
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
		Name string `json:"full_name"`
	}
	if err := json.Unmarshal(bytes, &repos); err != nil {
		return nil, errors.InternalServerError(p.name, "Invalid response returned from the GitHub API: %v", err)
	}

	repoos := make([]string, 0, len(repos))
	for _, r := range repos {
		repoos = append(repoos, strings.ToLower(r.Name))
	}

	return repoos, nil
}

func (p *Project) findProject(ctx context.Context, id string) (*project.Project, error) {
	// Identify the user
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// get the projects the user belongs to
	mRsp, err := p.project.ListMemberships(ctx, &project.ListMembershipsRequest{MemberId: userID})
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
	rRsp, err := p.project.Read(ctx, &project.ReadRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return rRsp.GetProject(), nil
}
