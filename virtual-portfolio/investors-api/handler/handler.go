package handler

import (
	"context"
	"strings"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/photos"
	proto "github.com/kytra-app/investors-api/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth       auth.Authenticator
	photos     photos.Service
	stocks     stocks.StocksService
	users      users.UsersService
	quotes     quotes.StockQuoteService
	portfolios portfolios.PortfoliosService
	trades     trades.TradesService
	followers  followers.FollowersService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, pics photos.Service, client client.Client) Handler {
	return Handler{
		auth:       auth,
		photos:     pics,
		stocks:     stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		trades:     trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		users:      users.NewUsersService("kytra-srv-v1-users:8080", client),
		followers:  followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		quotes:     quotes.NewStockQuoteService("kytra-srv-v1-stock-quote:8080", client),
	}
}

// Get retrieves the investor
func (h Handler) Get(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	user, err := h.users.Find(ctx, &users.User{Uuid: req.Uuid, Username: req.Username})
	if err != nil {
		return err
	}

	// Check if the current user follows this investor
	following, _ := h.getFollowingStatus(ctx, user)
	followingCount, followersCount, _ := h.getFollowingFollowersCount(ctx, user)

	rsp.User = &proto.User{
		Uuid:              user.Uuid,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Username:          user.Username,
		CreatedAt:         time.Unix(user.CreatedAt, 0).String(),
		ProfilePictureUrl: h.photos.GetURL(user.ProfilePictureId),
		Following:         following,
		FollowingCount:    followingCount,
		FollowersCount:    followersCount,
	}

	hContext, hCancel := context.WithTimeout(ctx, 250*time.Millisecond)
	defer hCancel()
	rsp.User.CurrentHoldings, err = h.holdingsForUser(hContext, user)

	return nil
}

// GetFollowingAndFollowers finds the users who are following and being followed by the user
func (h Handler) GetFollowingAndFollowers(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	user, err := h.users.Find(ctx, &users.User{Uuid: req.Uuid, Username: req.Username})
	if err != nil {
		return err
	}

	following, followers, err := h.getFollowingFollowersList(ctx, user)
	rsp.User = &proto.User{FollowingList: following, FollowersList: followers}

	return nil
}

// Discover retries a list of investors the user may wish to follow
func (h Handler) Discover(ctx context.Context, req *proto.DiscoverRequest, rsp *proto.DiscoverResponse) error {
	usersRsp, err := h.users.All(ctx, &users.AllRequest{})
	if err != nil {
		return err
	}

	rsp.Users = make([]*proto.User, len(usersRsp.Users))
	followings, err := h.getFollowingStatuses(ctx, usersRsp.Users)
	if err != nil {
		return err
	}

	for i, user := range usersRsp.Users {
		rsp.Users[i] = &proto.User{
			Uuid:              user.Uuid,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Username:          user.Username,
			CreatedAt:         time.Unix(user.CreatedAt, 0).String(),
			ProfilePictureUrl: h.photos.GetURL(user.ProfilePictureId),
			Following:         followings[user.Uuid],
		}
	}

	return nil
}

// Follow creates a follower relationship between the requested user and the authenticated uses
func (h Handler) Follow(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	if req.Uuid == "" {
		return errors.BadRequest("UUID_REQUIRED", "A UUID is required")
	}

	if u.UUID == req.Uuid {
		return errors.BadRequest("UUID_INVALID", "You cannot follow yourself")
	}

	followReq := followers.Request{
		Followee: &followers.Resource{Uuid: req.Uuid, Type: "User"},
		Follower: &followers.Resource{Uuid: u.UUID, Type: "User"},
	}

	_, err = h.followers.Follow(ctx, &followReq)
	return err
}

// Unfollow deletes the follower relationship between the requested user and the authenticated uses
func (h Handler) Unfollow(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	unFollowReq := followers.Request{
		Followee: &followers.Resource{Uuid: req.Uuid, Type: "User"},
		Follower: &followers.Resource{Uuid: u.UUID, Type: "User"},
	}

	_, err = h.followers.Unfollow(ctx, &unFollowReq)
	return err
}

func (h Handler) holdingsForUser(ctx context.Context, user *users.User) ([]*proto.Position, error) {
	// Get Portfolio
	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.Uuid})
	if err != nil {
		return []*proto.Position{}, err
	}

	// Get positions
	tradesRsp, err := h.trades.ListPositionsForPortfolio(ctx, &trades.ListRequest{PortfolioUuid: portfolio.Uuid})
	if err != nil || len(tradesRsp.Positions) == 0 {
		return []*proto.Position{}, err
	}

	// Get asset (stock) UUIDs
	stockUUIDs := make([]string, len(tradesRsp.Positions))
	for i, position := range tradesRsp.Positions {
		stockUUIDs[i] = position.Asset.Uuid
	}

	// Get and group the stocks
	stocksRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
	stocksByUUID := make(map[string]*stocks.Stock, len(stocksRsp.Stocks))
	for _, s := range stocksRsp.Stocks {
		stocksByUUID[s.Uuid] = s
	}

	// Serialize the result
	result := make([]*proto.Position, len(tradesRsp.Positions))
	for i, pos := range tradesRsp.Positions {
		stock, foundStock := stocksByUUID[pos.Asset.Uuid]
		if !foundStock {
			continue
		}

		var unitPrice int64
		if quote, err := h.quotes.GetQuote(ctx, &quotes.Stock{Symbol: stock.Symbol}); err == nil {
			unitPrice = int64(quote.Price)
		}

		result[i] = &proto.Position{
			AssetType:  pos.Asset.Type,
			AssetUuid:  pos.Asset.Uuid,
			AssetName:  stock.Name,
			Quantity:   pos.Quantity,
			BookCost:   pos.BookCost,
			UnitValue:  unitPrice,
			TotalValue: unitPrice * pos.Quantity,
		}
	}

	return result, nil
}

func (h Handler) getFollowingStatus(ctx context.Context, followee *users.User) (bool, error) {
	rsp, err := h.getFollowingStatuses(ctx, []*users.User{followee})

	if following, ok := rsp[followee.Uuid]; ok {
		return following, nil
	}

	return false, err
}

func (h Handler) getFollowingStatuses(ctx context.Context, followees []*users.User) (map[string]bool, error) {
	follower, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return map[string]bool{}, nil
	}

	fContext, fCancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer fCancel()

	uuids := make([]string, len(followees))
	for i, user := range followees {
		uuids[i] = user.Uuid
	}

	fRsp, err := h.followers.List(fContext, &followers.ListRequest{
		Follower:      &followers.Resource{Type: "User", Uuid: follower.UUID},
		FolloweeType:  "User",
		FolloweeUuids: uuids,
	})

	if err != nil {
		return map[string]bool{}, err
	}

	rsp := make(map[string]bool, len(fRsp.Resources))
	for _, r := range fRsp.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}

func (h Handler) getFollowingFollowersCount(ctx context.Context, user *users.User) (int32, int32, error) {
	context, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()

	r := &followers.Resource{Type: "User", Uuid: user.Uuid}
	rsp, err := h.followers.Count(context, r)
	if err != nil {
		return 0, 0, err
	}

	return rsp.FollowingCount, rsp.FollowerCount, nil
}

func (h Handler) getFollowingFollowersList(ctx context.Context, user *users.User) ([]*proto.Follower, []*proto.Follower, error) {
	context, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer cancel()

	r := &followers.Resource{Type: "User", Uuid: user.Uuid}
	rsp, err := h.followers.Get(context, r)
	if err != nil {
		return []*proto.Follower{}, []*proto.Follower{}, err
	}

	following, err := h.serializeFollowerResources(ctx, rsp.Following)
	if err != nil {
		return []*proto.Follower{}, []*proto.Follower{}, err

	}

	followers, err := h.serializeFollowerResources(ctx, rsp.Followers)
	if err != nil {
		return []*proto.Follower{}, []*proto.Follower{}, err
	}

	return following, followers, nil
}

func (h Handler) serializeFollowerResources(ctx context.Context, resources []*followers.Resource) ([]*proto.Follower, error) {
	var userUUIDs []string
	for _, r := range resources {
		if r.Type == "User" {
			userUUIDs = append(userUUIDs, r.Uuid)
		}
	}
	usersRsp, err := h.users.List(ctx, &users.ListRequest{Uuids: userUUIDs})
	if err != nil {
		return []*proto.Follower{}, err
	}
	usersMap := make(map[string]*users.User, len(usersRsp.Users))
	for _, user := range usersRsp.Users {
		usersMap[user.Uuid] = user
	}

	var stockUUIDs []string
	for _, r := range resources {
		if r.Type == "Stock" {
			stockUUIDs = append(stockUUIDs, r.Uuid)
		}
	}
	stocksRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return []*proto.Follower{}, err
	}
	stocksMap := make(map[string]*stocks.Stock, len(stocksRsp.Stocks))
	for _, stock := range stocksRsp.Stocks {
		stocksMap[stock.Uuid] = stock
	}

	rsp := make([]*proto.Follower, len(resources))
	for i, r := range resources {
		switch r.Type {
		case "User":
			u := usersMap[r.Uuid]
			rsp[i] = &proto.Follower{
				Type:              "User",
				Uuid:              r.Uuid,
				Following:         r.Following,
				Name:              strings.Join([]string{u.FirstName, u.LastName}, " "),
				ProfilePictureUrl: h.photos.GetURL(u.ProfilePictureId),
				Username:          u.Username,
			}
		case "Stock":
			s := stocksMap[r.Uuid]
			rsp[i] = &proto.Follower{
				Type:              "Stock",
				Uuid:              r.Uuid,
				Following:         r.Following,
				Name:              s.Name,
				ProfilePictureUrl: h.photos.GetURL(s.ProfilePictureId),
			}
		}
	}

	return rsp, nil
}
