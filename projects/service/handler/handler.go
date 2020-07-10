package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/projects/service/proto"
)

// Project implements the project service interface
type Project struct {
	name  string
	store store.Store
}

// New returns an initialized project handler
func New(service micro.Service) *Project {
	return &Project{
		name:  service.Name(),
		store: service.Options().Store,
	}
}

const (
	// projectsPrefix is the store prefix for projects. projects are stored with
	// keys in the following format "project/{id}".
	projectsPrefix = "project/"
	// membersPrefix is the stroe prefix for memberships. Memberships are
	// stored with key in the following format "membership/{projectID}/{userID}".
	// The value is the user ID (string, stored as bytes).
	membersPrefix = "member/"
)

// Read looks up a project using id
func (p *Project) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// lookup the project
	var err error
	if len(req.Id) > 0 {
		rsp.Project, err = p.findProjectByID(req.Id)
	}
	if len(req.Name) > 0 {
		rsp.Project, err = p.findProjectByName(req.Name)
	}
	if err != nil {
		return err
	}

	// lookup the project members
	recs, err := p.store.Read(membersPrefix+rsp.Project.Id+"/", store.ReadPrefix())
	if err != nil {
		return nil
	}
	rsp.Project.Members = make([]*pb.Member, 0, len(recs))
	for _, r := range recs {
		var m *membership
		if err := json.Unmarshal(r.Value, &m); err != nil {
			return errors.BadRequest(p.name, "Error unmarshaling json: %v", err)
		}

		rsp.Project.Members = append(rsp.Project.Members, &pb.Member{
			Id: m.MemberID, Type: m.MemberType, Role: m.Role,
		})
	}

	return nil
}

// Create a new projects
func (p *Project) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Project == nil {
		return errors.BadRequest(p.name, "Missing project")
	}
	if len(req.Project.Name) == 0 {
		return errors.BadRequest(p.name, "Missing project name")
	}
	if len(req.Project.Repository) == 0 {
		return errors.BadRequest(p.name, "Missing project repository")
	}

	// add the default fields
	req.Project.Id = uuid.New().String()

	// write to the store
	if err := p.writeProjectToStore(req.Project); err != nil {
		return err
	}

	// return the project in the response
	rsp.Project = req.Project
	return nil
}

// Update a project
func (p *Project) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if req.Project == nil {
		return errors.BadRequest(p.name, "Missing project")
	}
	if len(req.Project.Id) == 0 {
		return errors.BadRequest(p.name, "Missing project id")
	}
	if len(req.Project.Name) == 0 {
		return errors.BadRequest(p.name, "Missing project name")
	}

	// lookup the project
	project, err := p.findProjectByID(req.Project.Id)
	if err != nil {
		return errors.BadRequest(p.name, "Error finding project: %v", err)
	}

	// assign the update params
	project.Description = req.Project.Description

	// write to the store
	if err := p.writeProjectToStore(req.Project); err != nil {
		return errors.InternalServerError(p.name, "Error writing project: %v", err)
	}

	return nil
}

// List all the projects (does not return membership)
func (p *Project) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// get the records with the project prefix
	recs, err := p.store.Read(projectsPrefix, store.ReadPrefix())
	if err != nil {
		return err
	}

	// unmarshal and return in the response
	rsp.Projects = make([]*pb.Project, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &rsp.Projects[i]); err != nil {
			return errors.InternalServerError(p.name, "Error unmarsaling json: %v", err)
		}
	}

	return nil
}

type membership struct {
	ProjectID  string
	MemberType string
	MemberID   string
	Role       pb.Role
}

func (m *membership) Key() string {
	return fmt.Sprintf("%v%v/%v/%v", membersPrefix, m.ProjectID, m.MemberType, m.MemberID)
}

func (m *membership) Bytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

// AddMember to a project
func (p *Project) AddMember(ctx context.Context, req *pb.AddMemberRequest, rsp *pb.AddMemberResponse) error {
	// validate the request
	if _, err := p.findProjectByID(req.ProjectId); err != nil {
		return err
	}
	if req.Member == nil {
		return errors.BadRequest(p.name, "Missing member")
	}
	if req.Role == pb.Role_Unknown {
		return errors.BadRequest(p.name, "Missing role")
	}

	// construct the membership
	m := &membership{
		ProjectID:  req.ProjectId,
		Role:       req.Role,
		MemberID:   req.Member.Id,
		MemberType: req.Member.Type,
	}

	// write the membership to the store
	return p.store.Write(&store.Record{Key: m.Key(), Value: m.Bytes()})
}

// RemoveMember from a project
func (p *Project) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest, rsp *pb.RemoveMemberResponse) error {
	// validate the request
	if req.Member == nil {
		return errors.BadRequest(p.name, "Missing member")
	}

	// construct the membership
	m := &membership{
		ProjectID:  req.ProjectId,
		Role:       req.Role,
		MemberID:   req.Member.Id,
		MemberType: req.Member.Type,
	}

	return p.store.Delete(m.Key())
}

// ListMemberships returns all the projects a member belongs to
func (p *Project) ListMemberships(ctx context.Context, req *pb.ListMembershipsRequest, rsp *pb.ListMembershipsResponse) error {
	// validate the request
	if req.Member == nil {
		return errors.BadRequest(p.name, "Missing member")
	}

	// member id is the last component of the key, so list all
	// the keys in the store which relate to memberships
	keys, err := p.store.List(store.ListPrefix(membersPrefix))
	if err != nil {
		return err
	}

	// filter to get the project ids which the member belongs to
	var projectIDs []string
	for _, k := range keys {
		if strings.HasSuffix(k, "/"+req.Member.Type+"/"+req.Member.Id) {
			projectIDs = append(projectIDs, strings.Split(k, "/")[1])
		}
	}

	// get each of the projects
	rsp.Projects = make([]*pb.Project, 0, len(projectIDs))
	for _, id := range projectIDs {
		project, err := p.findProjectByID(id)
		if err != nil {
			return err
		}
		rsp.Projects = append(rsp.Projects, project)
	}

	return nil
}

func (p *Project) findProjectByID(id string) (*pb.Project, error) {
	recs, err := p.store.Read(projectsPrefix + id)
	if err != nil {
		return nil, err
	}

	var project *pb.Project
	err = json.Unmarshal(recs[0].Value, &project)
	return project, err
}

func (p *Project) findProjectByName(name string) (*pb.Project, error) {
	recs, err := p.store.Read(projectsPrefix, store.ReadPrefix())
	if err != nil {
		return nil, err
	}

	for _, r := range recs {
		var project *pb.Project
		if err = json.Unmarshal(r.Value, &project); err != nil {
			return nil, errors.InternalServerError(p.name, "Error unmarsaling json: %v", err)
		}
		if project.Name == name {
			return project, nil
		}
	}

	return nil, store.ErrNotFound
}

// writeProjectToStore marshals a project and writes it to the store under
// the corresponding key (prefix + id)
func (p *Project) writeProjectToStore(project *pb.Project) error {
	bytes, err := json.Marshal(project)
	if err != nil {
		return errors.InternalServerError(p.name, "Error marsaling json: %v", err)
	}

	key := projectsPrefix + project.Id
	if err := p.store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(p.name, "Error writing to the store: %v", err)
	}

	return nil
}
