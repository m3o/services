package handler

import (
	"context"
	namespace "namespace/proto"
)

type Namespace struct{}

func (n Namespace) Create(ctx context.Context, in *namespace.CreateRequest, opts ...interface{}) (*namespace.CreateResponse, error) {
	panic("implement me")
}

func (n Namespace) Read(ctx context.Context, in *namespace.ReadRequest, opts ...interface{}) (*namespace.ReadResponse, error) {
	panic("implement me")
}

func (n Namespace) Delete(ctx context.Context, in *namespace.DeleteRequest, opts ...interface{}) (*namespace.DeleteResponse, error) {
	panic("implement me")
}

func (n Namespace) List(ctx context.Context, in *namespace.ListRequest, opts ...interface{}) (*namespace.ListResponse, error) {
	panic("implement me")
}

func (n Namespace) AddUser(ctx context.Context, in *namespace.AddUserRequest, opts ...interface{}) (*namespace.AddUserResponse, error) {
	panic("implement me")
}

func (n Namespace) RemoveUser(ctx context.Context, in *namespace.RemoveUserRequest, opts ...interface{}) (*namespace.RemoveUserResponse, error) {
	panic("implement me")
}
