package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	pb "github.com/micro/services/payments/provider/proto"
)

// Handler implements the payments provider interface for stripe
type Handler struct {
	name   string      // name of the service
	store  store.Store // go-micro store (key/value)
	client *client.API // stripe api client
}

// NewHandler returns an initialised Handler, it will error if any of
// the required enviroment variables are not set
func NewHandler(srv micro.Service) *Handler {
	apiKey := os.Getenv("STRIPE_API_KEY")
	if len(apiKey) == 0 {
		log.Fatalf("Missing required env: STRIPE_API_KEY")
	}

	return &Handler{
		store:  store.DefaultStore,
		client: client.New(apiKey, nil),
		name:   srv.Name(),
	}
}

// CreateProduct via the Stripe API, e.g. "Notes"
func (h *Handler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest, rsp *pb.CreateProductResponse) error {
	if req.Product == nil {
		return errors.BadRequest(h.name, "Product required")
	}

	// Construct the stripe product params
	params := &stripe.ProductParams{
		ID:          stripe.String(req.Product.Id),
		Name:        stripe.String(req.Product.Name),
		Description: stripe.String(req.Product.Description),
		Active:      stripe.Bool(req.Product.Active),
	}

	// Create the product
	_, err := h.client.Products.New(params)
	if err == nil {
		return nil
	}

	// Handle the error
	switch err.(*stripe.Error).Code {
	case stripe.ErrorCodeResourceAlreadyExists:
		// the product already exists, update it
		params.ID = nil // don't pass ID again in req body
		_, updateErr := h.client.Products.Update(req.Product.Id, params)
		return updateErr
	default:
		// the error was not expected
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}
}

// CreatePlan via the Stripe API, e.g. "Gold"
func (h *Handler) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest, rsp *pb.CreatePlanResponse) error {
	if req.Plan == nil {
		return errors.BadRequest(h.name, "Plan required")
	}

	// Format the interval
	interval := strings.ToLower(req.Plan.Interval.String())

	// Construct the stripe product plan params
	params := &stripe.PlanParams{
		ID:        stripe.String(req.Plan.Id),
		Nickname:  stripe.String(req.Plan.Name),
		Currency:  stripe.String(req.Plan.Currency),
		ProductID: stripe.String(req.Plan.ProductId),
		Interval:  stripe.String(interval),
		Amount:    stripe.Int64(req.Plan.Amount),
	}

	// Create the product plan
	_, err := h.client.Plans.New(params)
	if err == nil {
		return nil
	}

	// Handle the error
	switch err.(*stripe.Error).Code {
	case stripe.ErrorCodeResourceAlreadyExists:
		// the product plan already exists and it cannot be updated
		return nil
	default:
		// the error was not expected
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}
}

// CreateUser via the Stripe API, e.g. "John Doe"
func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest, rsp *pb.CreateUserResponse) error {
	if req.User == nil {
		return errors.BadRequest(h.name, "User required")
	}
	if req.User.Metadata == nil {
		req.User.Metadata = make(map[string]string, 0)
	}

	// Check to see if the user has already been created
	stripeID, err := h.getStripeIDForUser(req.User.Id)
	if err != nil {
		return err
	}

	// Construct the params
	var params stripe.CustomerParams
	if email := req.User.Metadata["email"]; len(email) > 0 {
		params.Email = stripe.String(email)
	}
	if name := req.User.Metadata["name"]; len(name) > 0 {
		params.Name = stripe.String(name)
	}
	if phone := req.User.Metadata["phone"]; len(phone) > 0 {
		params.Phone = stripe.String(phone)
	}

	// If the user already exists, update using the existing attrbutes
	if len(stripeID) > 0 {
		if _, err := h.client.Customers.Update(stripeID, &params); err != nil {
			return errors.InternalServerError(h.name, "Unexepcted stripe update error: %v", err)
		}
		return nil
	}

	// Create the user in stripe
	c, err := h.client.Customers.New(&params)
	if err != nil {
		return errors.InternalServerError(h.name, "Unexepcted stripe create error: %v", err)
	}

	// Write the ID to the database
	return h.setStripeIDForUser(c.ID, req.User.Id)
}

// CreateSubscription via the Stripe API, e.g. "Subscribe John Doe to Notes Gold"
func (h *Handler) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest, rsp *pb.CreateSubscriptionResponse) error {
	return nil
}

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

// User is the datatype stored in the store
type User struct {
	StripeID string `json:"stripe_id"`
}

// getStripeIDForUser returns the stripe ID from the store for the given user
func (h *Handler) getStripeIDForUser(userID string) (string, error) {
	recs, err := h.store.Read(userID)
	if err == store.ErrNotFound || len(recs) == 0 {
		return "", nil
	} else if err != nil {
		return "", errors.InternalServerError(h.name, "Could not read from store: %v", err)
	}

	var user *User
	if err := json.Unmarshal(recs[0].Value, &user); err != nil {
		return "", errors.InternalServerError(h.name, "Could not unmarshal json: %v", err)
	}

	return user.StripeID, nil
}

// setStripeIDForUser writes the stripe ID to the store for the given user
func (h *Handler) setStripeIDForUser(stripeID, userID string) error {
	bytes, err := json.Marshal(&User{StripeID: stripeID})
	if err != nil {
		return errors.InternalServerError(h.name, "Could not marshal json: %v", err)
	}

	if err := h.store.Write(&store.Record{Key: userID, Value: bytes}); err != nil {
		return errors.InternalServerError(h.name, "Could not write to store: %v", err)
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
