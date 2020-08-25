package handler

import (
	"context"

	subscription "subscription/proto"
)

type Subscription struct{}

func (s Subscription) Create(ctx context.Context, in *subscription.CreateRequest, opts ...interface{}) (*subscription.CreateResponse, error) {
	panic("implement me")
}

func (s Subscription) Cancel(ctx context.Context, in *subscription.CancelRequest, opts ...interface{}) (*subscription.CancelResponse, error) {
	panic("implement me")
}

func (s Subscription) List(ctx context.Context, in *subscription.ListRequest, opts ...interface{}) (*subscription.ListResponse, error) {
	panic("implement me")
}
