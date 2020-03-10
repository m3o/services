package handler

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/services/payments/provider/proto"
	stripe "github.com/stripe/stripe-go"
)

// CreatePaymentMethod via the Stripe API, e.g. "Add payment method pm_s93483932 to John Doe"
func (h *Handler) CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest, rsp *pb.CreatePaymentMethodResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest(h.name, "ID required")
	}
	if len(req.UserId) == 0 {
		return errors.BadRequest(h.name, "User ID required")
	}

	// Check to see if the user has exists
	stripeID, err := h.getStripeIDForUser(req.UserId)
	if err != nil {
		return err
	}
	if stripeID == "" {
		return errors.BadRequest(h.name, "User ID doesn't exist")
	}

	// Create the payment method
	pm, err := h.client.PaymentMethods.Attach(req.Id, &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(stripeID),
	})
	if err != nil {
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}

	// Serialize the response
	rsp.PaymentMethod = serializePaymentMethod(pm, req.UserId)
	return nil
}

// DeletePaymentMethod via the Stripe API, e.g. "Remove payment method pm_s93483932"
func (h *Handler) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest, rsp *pb.DeletePaymentMethodResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest(h.name, "ID required")
	}

	// Delete the payment method
	_, err := h.client.PaymentMethods.Detach(req.Id, &stripe.PaymentMethodDetachParams{})
	if err != nil {
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}
	return nil
}

// ListPaymentMethods via the Stripe API, e.g. "List payment methods for John Doe"
func (h *Handler) ListPaymentMethods(ctx context.Context, req *pb.ListPaymentMethodsRequest, rsp *pb.ListPaymentMethodsResponse) error {
	if len(req.UserId) == 0 {
		return errors.BadRequest(h.name, "User ID required")
	}

	// Check to see if the user has exists
	stripeID, err := h.getStripeIDForUser(req.UserId)
	if err != nil {
		return err
	}
	if stripeID == "" {
		return errors.BadRequest(h.name, "User ID doesn't exist")
	}

	// List the payment methods
	iter := h.client.PaymentMethods.List(&stripe.PaymentMethodListParams{
		Customer: stripe.String(stripeID),
	})
	if iter.Err() != nil {
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}

	// Loop through and serialize
	rsp.PaymentMethods = make([]*pb.PaymentMethod, 0)
	for {
		pm := serializePaymentMethod(iter.PaymentMethod(), req.UserId)
		rsp.PaymentMethods = append(rsp.PaymentMethods, pm)

		if !iter.Next() {
			break
		}
	}

	return nil
}

func serializePaymentMethod(pm *stripe.PaymentMethod, userID string) *pb.PaymentMethod {
	rsp := &pb.PaymentMethod{
		Id:      pm.ID,
		Created: pm.Created,
		UserId:  userID,
		Type:    fmt.Sprint(pm.Type),
	}

	if pm.Type == stripe.PaymentMethodTypeCard && pm.Card != nil {
		rsp.CardBrand = fmt.Sprint(pm.Card.Brand)
		rsp.CardExpMonth = fmt.Sprint(pm.Card.ExpMonth)
		rsp.CardExpYear = fmt.Sprint(pm.Card.ExpYear)
		rsp.CardLast_4 = fmt.Sprint(pm.Card.Last4)
	}

	return rsp
}
