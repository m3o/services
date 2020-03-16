package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/services/payments/provider"
	payments "github.com/micro/services/payments/provider/proto"

	pb "github.com/micro/services/apps/service/proto/apps"
	"github.com/micro/services/apps/service/store"
)

// Handler implments the
type Handler struct {
	name     string
	store    *store.Store
	payments provider.Provider
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:     srv.Name(),
		store:    store.NewStore(srv),
		payments: payments.NewProviderService("go.micro.srv.payment.stripe", srv.Client()),
	}
}

// Create a new app
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// Validate the app
	if err := h.validateApp(req.App); err != nil {
		return err
	}

	if _, err := h.store.Read(req.App.Id); err == nil {
		return errors.BadRequest(h.name, "An app with this ID already exists")
	} else if err != store.ErrNotFound {
		return err
	}

	return h.updateAppAndPricing(ctx, req.App)
}

// Update an app, looking it up by ID
func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if err := h.validateApp(req.App); err != nil {
		return err
	}

	if _, err := h.store.Read(req.App.Id); err == store.ErrNotFound {
		return errors.NotFound(h.name, err.Error())
	} else if err != nil {
		return err
	}

	return h.updateAppAndPricing(ctx, req.App)
}

// Read an app, looking it up by ID
func (h *Handler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	app, err := h.store.Read(req.Id)
	if err == store.ErrNotFound {
		return errors.NotFound(h.name, err.Error())
	} else if err != nil {
		return err
	}

	rsp.App = app
	return nil
}

// List all the apps
func (h *Handler) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	apps, err := h.store.List()
	if err != nil {
		return err
	}
	rsp.Apps = apps
	return nil
}

func (h *Handler) updateAppAndPricing(ctx context.Context, app *pb.App) error {
	// Create the product in the payments provider (retry safe)
	_, err := h.payments.CreateProduct(ctx, &payments.CreateProductRequest{
		Product: &payments.Product{
			Id:          app.Id,
			Name:        app.Name,
			Description: app.Description,
			Active:      true,
		},
	})
	if err != nil {
		return err
	}

	// Create the plans in the payment provider (retry safe)
	for _, p := range app.Plans {
		var interval payments.PlanInterval

		switch p.Interval {
		case "day":
			interval = payments.PlanInterval_DAY
		case "week":
			interval = payments.PlanInterval_WEEK
		case "month":
			interval = payments.PlanInterval_MONTH
		case "year":
			interval = payments.PlanInterval_YEAR
		default:
			return errors.BadRequest(h.name, "Invalid plan interval:  %v", p.Interval)
		}

		_, err := h.payments.CreatePlan(ctx, &payments.CreatePlanRequest{
			Plan: &payments.Plan{
				ProductId: app.Id,
				Id:        p.Id,
				Currency:  "USD",
				Name:      p.Name,
				Amount:    p.Amount,
				Interval:  interval,
			},
		})
		if err != nil {
			return err
		}
	}

	return h.store.Write(app)
}

func (h *Handler) validateApp(app *pb.App) error {
	if app == nil {
		return errors.BadRequest(h.name, "Missing app")
	}
	if len(app.Id) == 0 {
		return errors.BadRequest(h.name, "Missing app id")
	}
	if len(app.Name) == 0 {
		return errors.BadRequest(h.name, "Missing app name")
	}
	if len(app.Category) == 0 {
		return errors.BadRequest(h.name, "Missing app category")
	}
	if len(app.Plans) == 0 {
		return errors.BadRequest(h.name, "Missing app plan, at least one is requried")
	}

	for i, p := range app.Plans {
		if len(p.Name) == 0 {
			return errors.BadRequest(h.name, "Missing app plan %v name", i)
		}
		if len(p.Interval) == 0 {
			return errors.BadRequest(h.name, "Missing app plan %v interval", i)
		}
	}

	return nil
}
