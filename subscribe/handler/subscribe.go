package handler

import (
	"context"

	s "subscribe/proto/subscribe"
)

type Subscribe struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Subscribe) Subscribe(ctx context.Context, req *s.SubscribeRequest, rsp *s.SubscribeResponse) error {
	return nil
}

func (e *Subscribe) ListSubscriptions(ctx context.Context, req *s.ListSubscriptionsRequest, rsp *s.ListSubscriptionsResponse) error {
	return nil
}
