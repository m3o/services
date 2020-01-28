package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kytra-app/helpers/microtime"
	valuation "github.com/kytra-app/portfolio-valuation-srv/proto"
	proto "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	"github.com/kytra-app/portfolio-value-tracking-srv/storage"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	db         storage.Service
	portfolios portfolios.PortfoliosService
	valuation  valuation.PortfolioValuationService
}

// New returns an instance of Handler
func New(storage storage.Service, client client.Client) *Handler {
	return &Handler{
		db:         storage,
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		valuation:  valuation.NewPortfolioValuationService("kytra-srv-v1-portfolio-valuation:8080", client),
	}
}

// ListValuations returns the latest values of the portfolios
func (h *Handler) ListValuations(ctx context.Context, req *proto.ListValuationsRequest, rsp *proto.ListValuationsResponse) error {
	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	history, err := h.db.ListValuations(req.PortfolioUuids, time)
	if err != nil {
		return err
	}

	rsp.Valuations = make([]*proto.Valuation, len(history))
	for i, h := range history {
		rsp.Valuations[i] = &proto.Valuation{
			Amount:        h.Value,
			Date:          h.CreatedAt.Unix(),
			PortfolioUuid: h.PortfolioUUID,
		}
	}

	return nil
}

// GetDailyHistory returns the historical values of the portfolio, with a single value per date
func (h *Handler) GetDailyHistory(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if req.Uuid == "" {
		return errors.BadRequest("MissingUUID", "A UUID is required to lookup a portfolios history")
	}

	history, err := h.db.GetDailyHistory(req.Uuid)
	if err != nil {
		return err
	}

	rsp.Uuid = req.Uuid
	rsp.History = make([]*proto.Valuation, len(history))

	for i, h := range history {
		rsp.History[i] = &proto.Valuation{
			Amount: h.Value, Date: h.Date.Truncate(time.Hour * 24).Unix(),
		}
	}

	return nil
}

// GetIntradayHistory returns the historical values of the portfolio, with a single value per date
func (h *Handler) GetIntradayHistory(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if req.Uuid == "" {
		return errors.BadRequest("MissingUUID", "A UUID is required to lookup a portfolios history")
	}

	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	history, err := h.db.GetIntradayHistory(req.Uuid, time)
	if err != nil {
		return err
	}

	rsp.Uuid = req.Uuid
	rsp.History = make([]*proto.Valuation, len(history))

	for i, h := range history {
		rsp.History[i] = &proto.Valuation{
			Amount: h.Value, Date: h.CreatedAt.Unix(),
		}
	}

	return nil
}

// GetPriceMovement gives the price movements for the requested portfolio
func (h *Handler) GetPriceMovement(ctx context.Context, req *proto.GetPriceMovementsRequest, rsp *proto.GetPriceMovementsResponse) error {
	if req.StartDate == 0 {
		return errors.BadRequest("START_DATE_REQUIRED", "A start date is required")
	}
	if req.EndDate == 0 {
		return errors.BadRequest("END_DATE_REQUIRED", "An end date is required")
	}

	startDate := time.Unix(req.StartDate, 0)
	endDate := time.Unix(req.EndDate, 0)

	result, err := h.db.GetPriceMovements([]string{req.PortfolioUuid}, startDate, endDate)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return nil
	}

	p := result[0]
	rsp.PriceMovement = &proto.PriceMovement{
		PortfolioUuid:    p.PortfolioUUID,
		PercentageChange: p.PercentageChange(),
		EarliestValue:    p.EarliestValue,
		LatestValue:      p.LatestValue,
	}

	return nil
}

// ListPriceMovements gives the price movements for the requested portfolios over the range requested
func (h *Handler) ListPriceMovements(ctx context.Context, req *proto.ListPriceMovementsRequest, rsp *proto.ListPriceMovementsResponse) error {
	if req.StartDate == 0 {
		return errors.BadRequest("START_DATE_REQUIRED", "A start date is required")
	}
	if req.EndDate == 0 {
		return errors.BadRequest("END_DATE_REQUIRED", "An end date is required")
	}

	startDate := time.Unix(req.StartDate, 0)
	endDate := time.Unix(req.EndDate, 0)

	result, err := h.db.GetPriceMovements(req.GetPortfolioUuids(), startDate, endDate)
	if err != nil {
		return err
	}

	rsp.PriceMovements = make([]*proto.PriceMovement, len(result))
	for i, p := range result {
		rsp.PriceMovements[i] = &proto.PriceMovement{
			PortfolioUuid:    p.PortfolioUUID,
			PercentageChange: p.PercentageChange(),
			EarliestValue:    p.EarliestValue,
			LatestValue:      p.LatestValue,
		}
	}

	return nil
}

// RecordValuations is called by the daily CRON job. It created a valuation record for each
// active portfolio.
func (h *Handler) RecordValuations() {
	const ResultsPerBatch = 100
	fmt.Println("RecordValuations")

	today := time.Now().Truncate(time.Hour * 24)
	marketsOpen := today.Add(time.Hour * 14).Add(time.Minute - 30).Unix()
	marketsClose := today.Add(time.Hour * 21).Add(time.Minute + 30).Unix()
	now := time.Now().Unix()

	if now < marketsOpen || now > marketsClose {
		fmt.Println("Skipping as outside of market hrs")
		return // Outside trading hours
	}

	for page := 0; ; page++ {
		pRsp, err := h.portfolios.All(context.Background(), &portfolios.AllRequest{
			Page: int64(page), Limit: ResultsPerBatch,
		})

		if err != nil {
			log.Fatal(err)
		}

		for _, p := range pRsp.Portfolios {
			h.createValuation(p)
		}

		// No more results to fetch
		if len(pRsp.Portfolios) != ResultsPerBatch {
			break
		}
	}
}

func (h *Handler) createValuation(p *portfolios.Portfolio) {
	fmt.Println(p.Uuid)
	val, err := h.valuation.GetPortfolio(context.Background(), &valuation.Portfolio{Uuid: p.Uuid})

	if err != nil {
		log.Printf("Error getting valuation for Portfolio %v, %v", p.Uuid, err)
		return
	}

	_, err = h.db.Create(storage.Valuation{
		PortfolioUUID: p.Uuid,
		Value:         val.TotalValue,
		Date:          time.Now(),
	})

	if err != nil {
		log.Printf("Error saving valuation for Portfolio %v, %v", p.Uuid, err)
		return
	}
}
