package handler

import (
	"context"
	"fmt"

	pb "github.com/m3o/services/payments/provider/proto"
	"github.com/micro/go-micro/v3/errors"
	stripe "github.com/stripe/stripe-go/v71"
)

// CreatePaymentMethod via the Stripe API, e.g. "Add payment method pm_s93483932 to John Doe"
func (h *Provider) CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest, rsp *pb.CreatePaymentMethodResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest(h.name, "ID required")
	}
	if len(req.CustomerId) == 0 {
		return errors.BadRequest(h.name, "Customer ID required")
	}
	if len(req.CustomerType) == 0 {
		return errors.BadRequest(h.name, "Customer Type required")
	}

	// Check to see if the user has exists
	stripeID, err := h.getStripeIDForCustomer(req.CustomerType, req.CustomerId)
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
	rsp.PaymentMethod = serializePaymentMethod(pm, req.CustomerType, req.CustomerId)
	return nil
}

// DeletePaymentMethod via the Stripe API, e.g. "Remove payment method pm_s93483932"
func (h *Provider) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest, rsp *pb.DeletePaymentMethodResponse) error {
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
func (h *Provider) ListPaymentMethods(ctx context.Context, req *pb.ListPaymentMethodsRequest, rsp *pb.ListPaymentMethodsResponse) error {
	if len(req.CustomerType) == 0 {
		return errors.BadRequest(h.name, "Customer Type required")
	}
	if len(req.CustomerId) == 0 {
		return errors.BadRequest(h.name, "Customer ID required")
	}

	// Check to see if the user has exists
	stripeID, err := h.getStripeIDForCustomer(req.CustomerType, req.CustomerId)
	if err != nil {
		return err
	}
	if stripeID == "" {
		return errors.BadRequest(h.name, "User ID doesn't exist")
	}

	// Get the customers (need the default payment method ID)
	c, err := h.client.Customers.Get(stripeID, &stripe.CustomerParams{})
	if err != nil {
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}
	var defaultPaymentID string
	if c.InvoiceSettings != nil && c.InvoiceSettings.DefaultPaymentMethod != nil {
		defaultPaymentID = c.InvoiceSettings.DefaultPaymentMethod.ID
	}

	// List the payment methods
	iter := h.client.PaymentMethods.List(&stripe.PaymentMethodListParams{
		Customer: stripe.String(stripeID),
		Type:     stripe.String("card"),
	})
	if iter.Err() != nil {
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", iter.Err())
	}

	// Loop through and serialize
	rsp.PaymentMethods = make([]*pb.PaymentMethod, 0)
	for {
		if !iter.Next() {
			break
		}

		pm := serializePaymentMethod(iter.PaymentMethod(), req.CustomerType, req.CustomerId)
		if pm.Id == defaultPaymentID {
			pm.Default = true
		}
		rsp.PaymentMethods = append(rsp.PaymentMethods, pm)
	}

	return nil
}

// SetDefaultPaymentMethod sets the users default payment method
func (h *Provider) SetDefaultPaymentMethod(ctx context.Context, req *pb.SetDefaultPaymentMethodRequest, rsp *pb.SetDefaultPaymentMethodResponse) error {
	// Check to see if the user has already been created
	stripeID, err := h.getStripeIDForCustomer(req.CustomerType, req.CustomerId)
	if err != nil {
		return err
	}

	// Construct the params
	var params stripe.CustomerParams
	params.InvoiceSettings = &stripe.CustomerInvoiceSettingsParams{
		DefaultPaymentMethod: stripe.String(req.PaymentMethodId),
	}

	// Update the payment method
	if _, err := h.client.Customers.Update(stripeID, &params); err != nil {
		return errors.InternalServerError(h.name, "Unexepcted stripe update error: %v", err)
	}

	return nil
}

func serializePaymentMethod(pm *stripe.PaymentMethod, CustomerType, CustomerID string) *pb.PaymentMethod {
	rsp := &pb.PaymentMethod{
		Id:           pm.ID,
		Created:      pm.Created,
		CustomerId:   CustomerID,
		CustomerType: CustomerType,
		Type:         fmt.Sprint(pm.Type),
	}

	if pm.Type == stripe.PaymentMethodTypeCard && pm.Card != nil {
		rsp.CardBrand = fmt.Sprint(pm.Card.Brand)
		rsp.CardExpMonth = fmt.Sprint(pm.Card.ExpMonth)
		rsp.CardExpYear = fmt.Sprint(pm.Card.ExpYear)
		rsp.CardLast_4 = fmt.Sprint(pm.Card.Last4)
	}

	return rsp
}
