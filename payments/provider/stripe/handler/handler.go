package handler

import (
	"context"

	pb "github.com/micro/services/payments/provider/proto"
)

// Handler implements the payments provider interface for stripe
type Handler struct{}

// NewHandler returns an initialised Handler, it will error if any of
// the required enviroment variables are not set
func NewHandler() (Handler, error) {
	return Handler{}, nil
}

// Test is a test function
func (h *Handler) Test(ctx context.Context, req *pb.TestRequest, rsp *pb.TestResponse) error {
	return nil
}
