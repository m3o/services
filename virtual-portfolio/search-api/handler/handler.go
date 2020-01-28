package handler

import (
	"context"
	"fmt"
	"strings"

	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	photos "github.com/kytra-app/helpers/photos"
	"github.com/kytra-app/search-api/handler/scorer"
	proto "github.com/kytra-app/search-api/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
)

// Handler is an object can process RPC requests
type Handler struct {
	pics      photos.Service
	auth      auth.Authenticator
	users     users.UsersService
	stocks    stocks.StocksService
	followers followers.FollowersService
}

// New returns an instance of Handler
func New(client client.Client, auth auth.Authenticator, pics photos.Service) Handler {
	return Handler{
		auth:      auth,
		pics:      pics,
		users:     users.NewUsersService("kytra-srv-v1-users:8080", client),
		stocks:    stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		followers: followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
	}
}

// Search returns all the stocks and users who match the query
func (h Handler) Search(ctx context.Context, req *proto.Query, rsp *proto.Response) error {
	sc := scorer.New(req.Query)

	if strings.ToLower(req.Type) != "user" {
		stocksQuery := &stocks.SearchRequest{Query: req.Query, Limit: 25}
		stocksRsp, err := h.stocks.Search(ctx, stocksQuery)
		if err != nil {
			return err
		}

		stockFollows, _ := h.fetchFollowingsForStocks(ctx, stocksRsp.Stocks)
		for _, stock := range sc.SortStocks(stocksRsp.Stocks) {
			rsp.Results = append(rsp.Results, &proto.Result{
				Type:              "Stock",
				Uuid:              stock.Uuid,
				Title:             stock.Name,
				Subtitle:          fmt.Sprintf("%v:%v", stock.Exchange, stock.Symbol),
				ProfilePictureUrl: h.pics.GetURL(stock.ProfilePictureId),
				Color:             stock.Color,
				Description:       stock.Description,
				Following:         stockFollows[stock.Uuid],
			})
		}
	}

	if strings.ToLower(req.Type) != "stock" {
		usersQuery := &users.SearchRequest{Query: req.Query, Limit: 25}
		usersRsp, err := h.users.Search(ctx, usersQuery)
		if err != nil {
			return err
		}

		userFollows, _ := h.fetchFollowingsForUsers(ctx, usersRsp.Users)
		for _, user := range sc.SortUsers(usersRsp.Users) {
			rsp.Results = append(rsp.Results, &proto.Result{
				Type:              "Investor",
				Uuid:              user.Uuid,
				Title:             strings.Join([]string{user.FirstName, user.LastName}, " "),
				Subtitle:          fmt.Sprintf("@%v", user.Username),
				ProfilePictureUrl: h.pics.GetURL(user.ProfilePictureId),
				Following:         userFollows[user.Uuid],
			})
		}
	}

	return nil
}

func (h Handler) fetchFollowingsForStocks(ctx context.Context, stocks []*stocks.Stock) (map[string]bool, error) {
	// Setup response, default to not following users
	rsp := make(map[string]bool, len(stocks))
	for _, u := range stocks {
		rsp[u.Uuid] = false
	}

	// Try and get the user from the context
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return rsp, nil
	}

	// Get the user UUIDS
	uuids := make([]string, len(stocks))
	for i, user := range stocks {
		uuids[i] = user.Uuid
	}

	// Construct the query
	req := &followers.ListRequest{
		Follower:      &followers.Resource{Uuid: u.UUID, Type: "User"},
		FolloweeType:  "Stock",
		FolloweeUuids: uuids,
	}

	// Request the data
	data, err := h.followers.List(ctx, req)
	if err != nil {
		return rsp, err
	}

	// Update the response
	for _, r := range data.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}

func (h Handler) fetchFollowingsForUsers(ctx context.Context, users []*users.User) (map[string]bool, error) {
	// Setup response, default to not following users
	rsp := make(map[string]bool, len(users))
	for _, u := range users {
		rsp[u.Uuid] = false
	}

	// Try and get the user from the context
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return rsp, nil
	}

	// Get the user UUIDS
	uuids := make([]string, len(users))
	for i, user := range users {
		uuids[i] = user.Uuid
	}

	// Construct the query
	req := &followers.ListRequest{
		Follower:      &followers.Resource{Uuid: u.UUID, Type: "User"},
		FolloweeType:  "User",
		FolloweeUuids: uuids,
	}

	// Request the data
	data, err := h.followers.List(ctx, req)
	if err != nil {
		return rsp, err
	}

	// Update the response
	for _, r := range data.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}
