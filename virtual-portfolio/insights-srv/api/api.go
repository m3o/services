package api

import (
	"context"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	"github.com/kytra-app/helpers/microtime"
	proto "github.com/kytra-app/insights-srv/proto"
	"github.com/kytra-app/insights-srv/storage"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(storage storage.Service, client client.Client) *Handler {
	followersSrv := followers.NewFollowersService("kytra-srv-v1-followers:8080", client)
	return &Handler{storage, followersSrv}
}

// Handler is an object can process RPC requests
type Handler struct {
	db        storage.Service
	followers followers.FollowersService
}

// ListAssets returns all the assets which have insights for the time provided in the context
func (h *Handler) ListAssets(ctx context.Context, req *proto.ListAssetsRequest, rsp *proto.ListAssetsResponse) error {
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}
	date = date.Truncate(time.Hour * 24)

	assets, err := h.db.ListAssets(date, req.ExcludeNews)
	if err != nil {
		return err
	}

	rsp.Assets = make([]*proto.Asset, len(assets))
	for i, a := range assets {
		rsp.Assets[i] = &proto.Asset{Uuid: a.UUID, Type: a.Type}
	}

	return nil
}

// ListInsights returns all the insights for the requested user on the date provided in the context
func (h *Handler) ListInsights(ctx context.Context, req *proto.ListInsightsRequest, rsp *proto.ListInsightsResponse) error {
	if req.UserUuid == "" {
		return errors.BadRequest("MISSING_USER_UUID", "A user uuid is required")
	}

	// Step 1. Parse Date
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}
	date = date.Truncate(time.Hour * 24)

	// // Step 2. Get user specific insights
	// specificInsights, err := h.db.ListInsightsForUser(req.UserUuid, date)
	// if err != nil {
	// 	return err
	// }

	// // Step 3. Determine which stocks they're following
	// fRsp, err := h.followers.Get(ctx, &followers.Resource{Uuid: req.UserUuid, Type: "User"})
	// if err != nil {
	// 	return err
	// }
	// stockUUIDs := []string{}
	// for _, r := range fRsp.GetFollowing() {
	// 	if r.Type == "Stock" {
	// 		stockUUIDs = append(stockUUIDs, r.Uuid)
	// 	}
	// }

	// // Step 4. Get the generic insights for those stocks
	// genericInsights, err := h.db.ListInsightsForAssets("Stock", stockUUIDs, date)
	// if err != nil {
	// 	return err
	// }

	// // Step 5. Combine the insights and serialize
	// insights := append(specificInsights, genericInsights...)
	insights, _ := h.db.ListInsightsForUser(req.UserUuid, date)
	rsp.Insights = make([]*proto.Insight, len(insights))

	for i, insight := range insights {
		// Step 5.1. Determine when the user last saw the these insights
		asset := storage.Asset{Type: "Stock", UUID: insight.AssetUUID}
		view, _ := h.db.GetUserView(req.UserUuid, asset, date)

		// Step 5.2 Serialize the asset
		rsp.Insights[i] = &proto.Insight{
			Asset: &proto.Asset{
				Uuid: insight.AssetUUID,
				Type: insight.AssetType,
			},
			UserUuid:  req.UserUuid,
			LinkUrl:   insight.LinkURL,
			PostUuid:  insight.PostUUID,
			Title:     insight.Title,
			Subtitle:  insight.Subtitle,
			Type:      insight.Type,
			Seen:      view.CreatedAt.After(insight.CreatedAt),
			CreatedAt: insight.CreatedAt.Unix(),
		}
	}

	return nil
}

// CreateUserView records the user viewed the insights for an asset
func (h *Handler) CreateUserView(ctx context.Context, req *proto.UserView, rsp *proto.UserView) error {
	if req.UserUuid == "" {
		return errors.BadRequest("MISSING_USER_UUID", "A user uuid is required")
	}
	if req.Asset == nil {
		return errors.BadRequest("MISSING_ASSET", "An asset is required")
	}

	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	_, err = h.db.CreateUserView(storage.UserView{
		CreatedAt: date,
		UserUUID:  req.UserUuid,
		AssetUUID: req.Asset.Uuid,
		AssetType: req.Asset.Type,
	})

	return err
}
