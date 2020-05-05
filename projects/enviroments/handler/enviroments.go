package handler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/projects/enviroments/proto"
	projects "github.com/micro/services/projects/service/proto"
)

// NewEnviroments returns an initialised enviroments handler
func NewEnviroments(srv micro.Service) *Enviroments {
	return &Enviroments{
		name:     srv.Name(),
		store:    srv.Options().Store,
		projects: projects.NewProjectsService("go.micro.service.projects", srv.Client()),
	}
}

// Enviroments implements the proto service interface
type Enviroments struct {
	name     string
	store    store.Store
	projects projects.ProjectsService
}

// Create an enviroment
func (e *Enviroments) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Enviroment == nil {
		return errors.BadRequest(e.name, "Missing enviroment")
	}
	if len(req.Enviroment.Name) == 0 {
		return errors.BadRequest(e.name, "Missing enviroment name")
	}
	if len(req.Enviroment.ProjectId) == 0 {
		return errors.BadRequest(e.name, "Missing enviroment project id")
	}

	// lookup the project
	pRsp, err := e.projects.Read(ctx, &projects.ReadRequest{Id: req.Enviroment.ProjectId})
	if err != nil {
		return err
	}

	// generate the namespace (projectName.enviromentName)
	namespace := pRsp.Project.Name + "." + req.Enviroment.Name

	// validiate the namespace is unique
	if _, err := e.findEnviromentByNamespace(namespace); err == nil {
		return errors.BadRequest(e.name, "%v already taken in the %v project", req.Enviroment.Name, pRsp.Project.Name)
	} else if err != store.ErrNotFound {
		return err
	}

	// create the record
	req.Enviroment.Id = uuid.New().String()
	req.Enviroment.Namespace = namespace
	bytes, err := json.Marshal(req.Enviroment)
	if err != nil {
		return errors.InternalServerError(e.name, "Error marshaling record: %v", err)
	}
	key := req.Enviroment.ProjectId + "/" + req.Enviroment.Id
	if err := e.store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(e.name, "Error writing to store: %v", err)
	}
	return nil
}

// Read a singular enviroment using ID / Namespace or multiple enviroments using Project ID
func (e *Enviroments) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Id) > 0 {
		env, err := e.findEnviromentByID(req.Id)
		rsp.Enviroment = env
		return err
	}

	if len(req.Namespace) > 0 {
		env, err := e.findEnviromentByNamespace(req.Namespace)
		rsp.Enviroment = env
		return err
	}

	if len(req.ProjectId) > 0 {
		envs, err := e.findEnviromentsForProject(req.ProjectId)
		rsp.Enviroments = envs
		return err
	}

	return errors.BadRequest(e.name, "Missing ID / Namespace / ProjectID")
}

// Update an enviroment
func (e *Enviroments) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if req.Enviroment == nil {
		return errors.BadRequest(e.name, "Missing enviroment")
	}
	if len(req.Enviroment.Id) == 0 {
		return errors.BadRequest(e.name, "Missing enviroment id")
	}

	// lookup the enviroment
	env, err := e.findEnviromentByID(req.Enviroment.Id)
	if err == store.ErrNotFound {
		return errors.BadRequest(e.name, "Enviroment not found")
	} else if err != nil {
		return err
	}

	// assign the update attributees
	env.Description = req.Enviroment.Description

	// update in the store
	bytes, err := json.Marshal(env)
	if err != nil {
		return errors.InternalServerError(e.name, "Error marshaling record: %v", err)
	}
	key := env.ProjectId + "/" + env.Id
	if err := e.store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(e.name, "Error writing to store: %v", err)
	}
	return nil
}

// Delete an enviroment
func (e *Enviroments) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// lookup the enviroment
	env, err := e.findEnviromentByID(req.Id)
	if err == store.ErrNotFound {
		return errors.BadRequest(e.name, "Enviroment not found")
	} else if err != nil {
		return err
	}

	// delete from the store
	key := env.ProjectId + "/" + env.Id
	if err := e.store.Delete(key); err != nil {
		return errors.InternalServerError(e.name, "Error deleting from store: %v", err)
	}
	return nil
}

func (e *Enviroments) findEnviromentsForProject(id string) ([]*pb.Enviroment, error) {
	recs, err := e.store.Read(id+"/", store.ReadPrefix())
	if err != nil {
		return nil, err
	}

	envs := make([]*pb.Enviroment, 0, len(recs))
	for _, r := range recs {
		var env *pb.Enviroment
		if err := json.Unmarshal(r.Value, &env); err != nil {
			return nil, errors.InternalServerError(e.name, "Error unmarshaling record: %v", err)
		}
		envs = append(envs, env)
	}

	return envs, nil
}

func (e *Enviroments) findEnviromentByID(id string) (*pb.Enviroment, error) {
	keys, err := e.store.List()
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

	recs, err := e.store.Read(envKey)
	if err != nil {
		return nil, err
	}

	var env *pb.Enviroment
	if err := json.Unmarshal(recs[0].Value, &env); err != nil {
		return nil, errors.InternalServerError(e.name, "Error unmarshaling record: %v", err)
	}
	return env, nil
}

func (e *Enviroments) findEnviromentByNamespace(ns string) (*pb.Enviroment, error) {
	recs, err := e.store.Read("", store.ReadPrefix())
	if err != nil {
		return nil, err
	}

	for _, r := range recs {
		var env *pb.Enviroment
		if err := json.Unmarshal(r.Value, &env); err != nil {
			return nil, errors.InternalServerError(e.name, "Error unmarshaling record: %v", err)
		}
		if env.Namespace == ns {
			return env, nil
		}
	}

	return nil, store.ErrNotFound
}
