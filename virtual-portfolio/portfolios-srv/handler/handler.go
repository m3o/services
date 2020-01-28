package handler

import (
	"context"

	"github.com/kytra-app/helpers/microgorm"
	proto "github.com/kytra-app/portfolios-srv/proto"
	"github.com/kytra-app/portfolios-srv/storage"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(storage storage.Service) *Handler {
	return &Handler{db: storage}
}

// Handler is an object can process RPC requests
type Handler struct {
	db storage.Service
}

// Create inserts a new portfolio into the database
func (h *Handler) Create(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if req.UserUuid == "" {
		return errors.BadRequest("MissingUserUUID", "A UserUUID is required to create a portfolio")
	}

	_, err := h.db.Get(storage.Portfolio{
		UserUUID: req.UserUuid,

		IndustryTargetEnergy:                1,
		IndustryTargetUtilities:             3.5,
		IndustryTargetMaterials:             3,
		IndustryTargetRealEstate:            4,
		IndustryTargetFinancials:            14.5,
		IndustryTargetHealthCare:            13.5,
		IndustryTargetIndustrials:           10,
		IndustryTargetConsumerStaples:       7.5,
		IndustryTargetInformationTechnology: 22,
		IndustryTargetConsumerDiscretionary: 10.5,
		IndustryTargetCommunicationServices: 10.5,
	})

	if err == nil {
		return errors.BadRequest("PortfolioExists", "A portfolio with this UserUUID already exists")
	} else if err != microgorm.ErrNotFound {
		return err
	}

	p, err := h.db.Create(storage.Portfolio{UserUUID: req.UserUuid})
	if err != nil {
		return err
	}
	*rsp = *h.serializePortfolio(p)

	return nil
}

// Get returns the portfolio, found using the UUID
func (h *Handler) Get(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if req.UserUuid == "" && req.Uuid == "" {
		return errors.BadRequest("MissingUUID", "A UserUUID or UUID is required to lookup a portfolio")
	}

	p, err := h.db.Get(storage.Portfolio{UUID: req.Uuid, UserUUID: req.UserUuid})
	if err != nil {
		return err
	}
	*rsp = *h.serializePortfolio(p)
	return nil
}

// Update amends the targets for the portfolio, found using the UUID
func (h *Handler) Update(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if req.Uuid == "" {
		return errors.BadRequest("MissingUUID", "A UUID is required to update a portfolio")
	}

	p, err := h.db.Update(h.reverseSerializePortfolio(req))
	if err != nil {
		return err
	}

	*rsp = *h.serializePortfolio(p)
	return nil
}

// All returns all portfolios
func (h *Handler) All(ctx context.Context, req *proto.AllRequest, rsp *proto.AllResponse) error {
	portfolios, err := h.db.All(req.Page, req.Limit)
	if err != nil {
		return err
	}

	rsp.Portfolios = make([]*proto.Portfolio, len(portfolios))
	for i, p := range portfolios {
		rsp.Portfolios[i] = h.serializePortfolio(p)
	}

	return nil
}

// List returns all portfolios which match the criteria
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	var data []storage.Portfolio
	var err error

	if len(req.UserUuids) > 0 {
		data, err = h.db.ListByUserUUIDs(req.UserUuids)
	} else {
		data, err = h.db.ListByUUIDs(req.Uuids)
	}

	if err != nil {
		return err
	}

	rsp.Portfolios = make([]*proto.Portfolio, len(data))
	for i, p := range data {
		rsp.Portfolios[i] = h.serializePortfolio(p)
	}

	return nil
}

func (h *Handler) reverseSerializePortfolio(req *proto.Portfolio) storage.Portfolio {
	return storage.Portfolio{
		UUID:                                req.Uuid,
		AssetClassTargetStocks:              req.AssetClassTargetStocks,
		AssetClassTargetCash:                req.AssetClassTargetCash,
		IndustryTargetInformationTechnology: req.IndustryTargetInformationTechnology,
		IndustryTargetFinancials:            req.IndustryTargetFinancials,
		IndustryTargetEnergy:                req.IndustryTargetEnergy,
		IndustryTargetHealthCare:            req.IndustryTargetHealthCare,
		IndustryTargetMaterials:             req.IndustryTargetMaterials,
		IndustryTargetUtilities:             req.IndustryTargetUtilities,
		IndustryTargetRealEstate:            req.IndustryTargetRealEstate,
		IndustryTargetConsumerDiscretionary: req.IndustryTargetConsumerDiscretionary,
		IndustryTargetConsumerStaples:       req.IndustryTargetConsumerStaples,
		IndustryTargetCommunicationServices: req.IndustryTargetCommunicationServices,
		IndustryTargetIndustrials:           req.IndustryTargetIndustrials,
		IndustryTargetMiscellaneous:         req.IndustryTargetMiscellaneous,
	}
}

func (h *Handler) serializePortfolio(in storage.Portfolio) *proto.Portfolio {
	return &proto.Portfolio{
		Uuid:                                in.UUID,
		UserUuid:                            in.UserUUID,
		AssetClassTargetStocks:              in.AssetClassTargetStocks,
		AssetClassTargetCash:                in.AssetClassTargetCash,
		IndustryTargetInformationTechnology: in.IndustryTargetInformationTechnology,
		IndustryTargetFinancials:            in.IndustryTargetFinancials,
		IndustryTargetEnergy:                in.IndustryTargetEnergy,
		IndustryTargetHealthCare:            in.IndustryTargetHealthCare,
		IndustryTargetMaterials:             in.IndustryTargetMaterials,
		IndustryTargetUtilities:             in.IndustryTargetUtilities,
		IndustryTargetRealEstate:            in.IndustryTargetRealEstate,
		IndustryTargetConsumerDiscretionary: in.IndustryTargetConsumerDiscretionary,
		IndustryTargetConsumerStaples:       in.IndustryTargetConsumerStaples,
		IndustryTargetCommunicationServices: in.IndustryTargetCommunicationServices,
		IndustryTargetIndustrials:           in.IndustryTargetIndustrials,
		IndustryTargetMiscellaneous:         in.IndustryTargetMiscellaneous,
	}
}
