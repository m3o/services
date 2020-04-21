package handler

import (
	"context"
	"fmt"

	pb "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// HandleUserEvent handles the events published by the uses service
func (h *Handler) HandleUserEvent(ctx context.Context, event *users.Event) error {
	switch event.Type {
	case users.EventType_UserCreated, users.EventType_UserUpdated:
		req := pb.CreateCustomerRequest{
			Customer: &pb.Customer{
				Id:   event.User.Id,
				Type: "user",
				Metadata: map[string]string{
					"email": event.User.Email,
					"name":  fmt.Sprintf("%v %v", event.User.FirstName, event.User.LastName),
				},
			},
		}

		return h.CreateCustomer(ctx, &req, &pb.CreateCustomerResponse{})
	default:
		return nil
	}
}
