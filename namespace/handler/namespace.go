package handler

import (
	"context"
	namespace "namespace/proto"
)

type Namespace struct{}

func (n Namespace) Create(ctx context.Context, request *namespace.CreateRequest, response *namespace.CreateResponse) error {
	panic("implement me")
}

func (n Namespace) Read(ctx context.Context, request *namespace.ReadRequest, response *namespace.ReadResponse) error {
	panic("implement me")
}

func (n Namespace) Delete(ctx context.Context, request *namespace.DeleteRequest, response *namespace.DeleteResponse) error {
	panic("implement me")
}

func (n Namespace) List(ctx context.Context, request *namespace.ListRequest, response *namespace.ListResponse) error {
	panic("implement me")
}

func (n Namespace) AddUser(ctx context.Context, request *namespace.AddUserRequest, response *namespace.AddUserResponse) error {
	panic("implement me")
}

func (n Namespace) RemoveUser(ctx context.Context, request *namespace.RemoveUserRequest, response *namespace.RemoveUserResponse) error {
	panic("implement me")
}
