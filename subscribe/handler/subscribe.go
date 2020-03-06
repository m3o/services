package handler

import (
	"context"
	"fmt"

	s "subscribe/proto/subscribe"

	store "github.com/micro/go-micro/v2/store"
	serviceStore "github.com/micro/go-micro/v2/store/service"
)

const prefix = "subscribe"

type Subscribe struct {
	store store.Store
}

func NewSubscribe() *Subscribe {
	return &Subscribe{
		store: serviceStore.NewStore(),
	}
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Subscribe) Subscribe(ctx context.Context, req *s.SubscribeRequest, rsp *s.SubscribeResponse) error {
	return e.store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v", prefix, req.GetNamespace()),
		Value: []byte(req.GetEmail()),
	})
}

func (e *Subscribe) ListSubscriptions(ctx context.Context, req *s.ListSubscriptionsRequest, rsp *s.ListSubscriptionsResponse) error {
	records, err := e.store.Read(fmt.Sprintf("%v:%v", prefix, req.GetNamespace()))
	if err != nil {
		return err
	}
	subs := []*s.Subscription{}
	for _, record := range records {
		subs = append(subs, &s.Subscription{
			Email: string(record.Value),
		})
	}
	rsp.Subscriptions = subs
	return nil
}
