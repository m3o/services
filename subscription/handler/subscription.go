package handler

import (
	"context"
	subscription "subscription/proto"
)

type Subscription struct{}

func (s Subscription) Create(ctx context.Context, request *subscription.CreateRequest, response *subscription.CreateResponse) error {
	panic("implement me")
}

func (s Subscription) Cancel(ctx context.Context, request *subscription.CancelRequest, response *subscription.CancelResponse) error {
	panic("implement me")
}

func (s Subscription) List(ctx context.Context, request *subscription.ListRequest, response *subscription.ListResponse) error {
	panic("implement me")
}
