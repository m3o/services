package handler

import (
	"context"
	customer "customer/proto"
)

type Customer struct{}

func (c Customer) Create(ctx context.Context, request *customer.CreateRequest, response *customer.CreateResponse) error {
	panic("implement me")
}

func (c Customer) MarkVerified(ctx context.Context, request *customer.MarkVerifiedRequest, response *customer.MarkVerifiedResponse) error {
	panic("implement me")
}

func (c Customer) Read(ctx context.Context, request *customer.ReadRequest, response *customer.ReadResponse) error {
	panic("implement me")
}

func (c Customer) Delete(ctx context.Context, request *customer.DeleteRequest, response *customer.DeleteResponse) error {
	panic("implement me")
}
