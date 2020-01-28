package handler

import (
	"context"

	proto "github.com/kytra-app/ledger-srv/proto"
	"github.com/kytra-app/ledger-srv/storage"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	"github.com/kytra-app/helpers/microtime"
	"github.com/micro/go-micro/client"
)

// New returns an instance of Handler
func New(storage storage.Service, client client.Client) *Handler {
	return &Handler{
		db:         storage,
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	db         storage.Service
	portfolios portfolios.PortfoliosService
}

// CreateTransaction inserts a new transaction into the ledger
func (h *Handler) CreateTransaction(ctx context.Context, req *proto.Transaction, rsp *proto.Transaction) error {
	if err := h.verifyPortfolio(req.PortfolioUuid); err != nil {
		return err
	}

	t, err := h.db.CreateTransaction(storage.Transaction{
		PortfolioUUID: req.PortfolioUuid,
		Amount:        req.Amount,
		Type:          req.Type.String(),
	})

	if err != nil {
		return err
	}

	*rsp = proto.Transaction{
		PortfolioUuid: t.PortfolioUUID,
		Amount:        t.Amount,
		Uuid:          t.UUID,
	}

	switch t.Type {
	case "DEPOSIT":
		rsp.Type = proto.TransactionType_DEPOSIT
	case "WITHDRAWAL":
		rsp.Type = proto.TransactionType_WITHDRAWAL
	case "BUY_ASSET":
		rsp.Type = proto.TransactionType_BUY_ASSET
	case "SELL_ASSET":
		rsp.Type = proto.TransactionType_SELL_ASSET
	}

	return nil
}

// GetPortfolio returns the portfolio and its current balance
func (h *Handler) GetPortfolio(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if err := h.verifyPortfolio(req.Uuid); err != nil {
		return err
	}

	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	b, err := h.db.GetPortfolioBalance(time, req.Uuid)
	if err != nil {
		return err
	}
	*rsp = proto.Portfolio{Uuid: req.Uuid, CurrentBalance: b}

	return nil
}

func (h *Handler) verifyPortfolio(UUID string) error {
	_, err := h.portfolios.Get(context.Background(), &portfolios.Portfolio{Uuid: UUID})
	return err
}
