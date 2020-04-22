package handler

import (
	"context"

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
		users:   users.NewUsersService("go.micro.service.users", service.Client()),
		project: project.NewProjectService("go.micro.service.project", service.Client()),
	}
}

// Project implments the M3O project service proto
type Project struct {
	name    string
	users   users.UsersService
	project project.ProjectService
}

// Create a new project
func (p *Project) Create(ctx context.Context, req *pb.CreateProjectRequest, rsp *pb.CreateProjectResponse) error {
	// Identify the user
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return err
	}

	// create the project
	cRsp, err := p.project.Create(ctx, &project.CreateRequest{
		Project: &project.Project{
			Name:      req.Project.Name,
			Namespace: req.Project.Namespace,
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
	return nil
}

// Update a project
func (p *Project) Update(ctx context.Context, req *pb.UpdateProjectRequest, rsp *pb.UpdateProjectResponse) error {
	// Identify the user
	userID, err := p.userIDFromContext(ctx)
	if err != nil {
		return err
	}

	// get the projects the user belongs to
	mRsp, err := p.project.ListMemberships(ctx, &project.ListMembershipsRequest{MemberId: userID})
	if err != nil {
		return err
	}

	// check for membership
	var isMember bool
	for _, t := range mRsp.Projects {
		if t.Id == req.Id {
			isMember = true
			break
		}
	}
	if !isMember {
		return errors.Forbidden(p.name, "Not a member of this team")
	}

	// lookup the project
	rRsp, err := p.project.Read(ctx, &project.ReadRequest{Id: req.Id})
	if err != nil {
		return err
	}

	// update the project
	_, err = p.project.Update(ctx, &project.UpdateRequest{
		Project: &project.Project{
			Id:        req.Id,
			Name:      req.Name,
			Namespace: rRsp.Project.Namespace,
			WebDomain: req.WebDomain,
			ApiDomain: req.ApiDomain,
		},
	})

	return nil
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

func (p *Project) userIDFromContext(ctx context.Context) (string, error) {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return "", errors.InternalServerError(p.name, "Auth error: %v", err)
	}

	uRsp, err := p.users.Read(ctx, &users.ReadRequest{Email: acc.ID})
	if err != nil {
		return "", errors.InternalServerError(p.name, "Auth error: %v", err)
	}

	return uRsp.User.Id, nil
}

func serializeProject(p *project.Project) *pb.Project {
	return &pb.Project{
		Id:        p.Id,
		Name:      p.Name,
		Namespace: p.Namespace,
		ApiDomain: p.ApiDomain,
		WebDomain: p.WebDomain,
	}
}
