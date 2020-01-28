package handler

import (
	"context"

	proto "github.com/kytra-app/investors-api/proto"
	allocation "github.com/kytra-app/portfolio-allocation-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
)

// Discover retries a list of investors the user may wish to follow
func (h Handler) Discover(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {

	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return err
	}

	uRsp, err := h.users.All(ctx, &users.AllRequest{})
	if err != nil {
		return err
	}

	uuids := []string{}
	for _, u := range uRsp.GetUsers() {
		if u.Uuid == user.UUID {
			continue
		}

		uuids = append(uuids, u.Uuid)
	}
	aRsp, err := h.allocation.List(ctx, &allocation.ListRequest{
		UserUuids: uuids,
	})

	uuids = []string{}
	for _, p := range aRsp.GetPortfolios() {
		var holdingsCount int
		var stocksPct float32

		for _, c := range p.GetAssetClasses() {
			if c.Name != "Stocks" {
				continue
			}

			holdingsCount = len(c.GetHoldings())
			stocksPct = c.PercentOfPortfolio
		}

		if stocksPct < 70 || holdingsCount < 10 {
			continue
		}

		uuids = append(uuids, p.UserUuid)
	}

	return h.serializeUsers(ctx, rsp, uuids, aRsp.GetPortfolios()...)
}
