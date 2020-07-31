package handler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service"
	mstore "github.com/micro/micro/v3/service/store"

	pb "github.com/m3o/services/projects/environments/proto"
	projects "github.com/m3o/services/projects/service/proto"
)

// NewEnvironments returns an initialised Environments handler
func NewEnvironments(srv *service.Service) *Environments {
	return &Environments{
		name:     srv.Name(),
		projects: projects.NewProjectsService("go.micro.service.projects"),
	}
}

// Environments implements the proto service interface
type Environments struct {
	name     string
	projects projects.ProjectsService
}

// Create an Environment
func (e *Environments) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Environment == nil {
		return errors.BadRequest(e.name, "Missing Environment")
	}
	if len(req.Environment.Name) == 0 {
		return errors.BadRequest(e.name, "Missing Environment name")
	}
	if len(req.Environment.ProjectId) == 0 {
		return errors.BadRequest(e.name, "Missing Environment project id")
	}

	// lookup the project
	pRsp, err := e.projects.Read(ctx, &projects.ReadRequest{Id: req.Environment.ProjectId})
	if err != nil {
		return errors.BadRequest(e.name, "Error finding project: %v", err)
	}

	// generate the namespace (projectName-EnvironmentName)
	namespace := strings.ToLower(pRsp.Project.Name + "-" + req.Environment.Name)

	// validiate the namespace is unique
	if _, err := e.findEnvironmentByNamespace(namespace); err == nil {
		return errors.BadRequest(e.name, "%v already taken in the %v project", req.Environment.Name, pRsp.Project.Name)
	} else if err != store.ErrNotFound {
		return errors.InternalServerError(e.name, "Error reading from store: %v", err)
	}

	// create the record
	req.Environment.Id = uuid.New().String()
	req.Environment.Namespace = namespace
	bytes, err := json.Marshal(req.Environment)
	if err != nil {
		return errors.InternalServerError(e.name, "Error marshaling record: %v", err)
	}
	key := req.Environment.ProjectId + "/" + req.Environment.Id
	if err := mstore.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(e.name, "Error writing to store: %v", err)
	}

	// serialize the result
	rsp.Environment = req.Environment
	return nil
}

// Read a singular Environment using ID / Namespace or multiple Environments using Project ID
func (e *Environments) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Id) > 0 {
		env, err := e.findEnvironmentByID(req.Id)
		rsp.Environment = env
		return err
	}

	if len(req.Namespace) > 0 {
		env, err := e.findEnvironmentByNamespace(req.Namespace)
		if err == store.ErrNotFound {
			return errors.BadRequest(e.name, "Environment with %v namespace not found", req.Namespace)
		}
		rsp.Environment = env
		return err
	}

	if len(req.ProjectId) > 0 {
		envs, err := e.findEnvironmentsForProject(req.ProjectId)
		if err == store.ErrNotFound {
			return errors.BadRequest(e.name, "Environment with %v id not found", req.Id)
		}
		rsp.Environments = envs
		return err
	}

	return errors.BadRequest(e.name, "Missing ID / Namespace / ProjectID")
}

// Update an Environment
func (e *Environments) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if req.Environment == nil {
		return errors.BadRequest(e.name, "Missing Environment")
	}
	if len(req.Environment.Id) == 0 {
		return errors.BadRequest(e.name, "Missing Environment id")
	}

	// lookup the Environment
	env, err := e.findEnvironmentByID(req.Environment.Id)
	if err == store.ErrNotFound {
		return errors.BadRequest(e.name, "Environment not found")
	} else if err != nil {
		return err
	}

	// assign the update attributees
	env.Description = req.Environment.Description

	// update in the store
	bytes, err := json.Marshal(env)
	if err != nil {
		return errors.InternalServerError(e.name, "Error marshaling record: %v", err)
	}
	key := env.ProjectId + "/" + env.Id
	if err := mstore.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(e.name, "Error writing to store: %v", err)
	}
	return nil
}

// Delete an Environment
func (e *Environments) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// lookup the Environment
	env, err := e.findEnvironmentByID(req.Id)
	if err == store.ErrNotFound {
		return errors.BadRequest(e.name, "Environment not found")
	} else if err != nil {
		return err
	}

	// delete from the store
	key := env.ProjectId + "/" + env.Id
	if err := mstore.Delete(key); err != nil {
		return errors.InternalServerError(e.name, "Error deleting from store: %v", err)
	}
	return nil
}

func (e *Environments) findEnvironmentsForProject(id string) ([]*pb.Environment, error) {
	recs, err := mstore.Read(id+"/", store.ReadPrefix())
	if err != nil {
		return nil, err
	}

	envs := make([]*pb.Environment, 0, len(recs))
	for _, r := range recs {
		var env *pb.Environment
		if err := json.Unmarshal(r.Value, &env); err != nil {
			return nil, errors.InternalServerError(e.name, "Error unmarshaling record: %v", err)
		}
		envs = append(envs, env)
	}

	return envs, nil
}

func (e *Environments) findEnvironmentByID(id string) (*pb.Environment, error) {
	keys, err := mstore.List()
	if err != nil {
		return nil, err
	}

	var envKey string
	for _, k := range keys {
		if strings.HasSuffix(k, "/"+id) {
			envKey = k
			break
		}
	}
	if len(envKey) == 0 {
		return nil, store.ErrNotFound
	}

	recs, err := mstore.Read(envKey)
	if err != nil {
		return nil, err
	}

	var env *pb.Environment
	if err := json.Unmarshal(recs[0].Value, &env); err != nil {
		return nil, errors.InternalServerError(e.name, "Error unmarshaling record: %v", err)
	}
	return env, nil
}

func (e *Environments) findEnvironmentByNamespace(ns string) (*pb.Environment, error) {
	recs, err := mstore.Read("", store.ReadPrefix())
	if err != nil {
		return nil, err
	}

	for _, r := range recs {
		var env *pb.Environment
		if err := json.Unmarshal(r.Value, &env); err != nil {
			return nil, errors.InternalServerError(e.name, "Error unmarshaling record: %v", err)
		}
		if env.Namespace == ns {
			return env, nil
		}
	}

	return nil, store.ErrNotFound
}
