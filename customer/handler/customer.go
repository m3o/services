package handler

import (
	"context"
	customer "customer/proto"
)

type Customer struct{}

func (c Customer) Create(ctx context.Context, in *customer.CreateRequest, opts ...interface{}) (*customer.CreateResponse, error) {
	panic("implement me")
}

func (c Customer) MarkVerified(ctx context.Context, in *customer.MarkVerifiedRequest, opts ...interface{}) (*customer.MarkVerifiedResponse, error) {
	panic("implement me")
}

func (c Customer) Read(ctx context.Context, in *customer.ReadRequest, opts ...interface{}) (*customer.ReadResponse, error) {
	panic("implement me")
}

func (c Customer) Delete(ctx context.Context, in *customer.DeleteRequest, opts ...interface{}) (*customer.DeleteResponse, error) {
	panic("implement me")
}
