package handler

import (
	"context"

	proto "github.com/kytra-app/stocks-srv/proto"
	"github.com/kytra-app/stocks-srv/storage"
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

// Create inserts a new stock into the database
func (h *Handler) Create(ctx context.Context, req *proto.Stock, rsp *proto.Response) error {
	s, err := h.db.Create(h.permittedParams(req))

	if err != nil {
		rsp.Error = &proto.Error{Code: 400, Message: err.Error()}
		return nil
	}

	rsp.Stock = h.serializeStock(&s)
	return nil
}

// Get returns the stock, found using the UUID
func (h *Handler) Get(ctx context.Context, req *proto.Stock, rsp *proto.Response) error {
	if req.Uuid == "" && req.Symbol == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing UUID or Symbol"}
		return nil
	}

	s, err := h.db.Get(storage.Stock{UUID: req.Uuid, Symbol: req.Symbol})
	if err != nil {
		rsp.Error = &proto.Error{Code: 404}
		return err
	}

	rsp.Stock = h.serializeStock(&s)
	return nil
}

// Delete deletes the stock found using the UUID
func (h *Handler) Delete(ctx context.Context, req *proto.Stock, rsp *proto.Response) error {
	if req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing UUID"}
		return nil
	}

	if err := h.db.Delete(storage.Stock{UUID: req.Uuid}); err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching stock"}
		return err
	}

	return nil
}

// Update finds the Stock with the UUID provided and updates that object
func (h *Handler) Update(ctx context.Context, req *proto.Stock, rsp *proto.Response) error {
	params := h.permittedParams(req)
	params.UUID = req.Uuid // UUID is needed for update but not create

	stock, err := h.db.Update(params)
	if err != nil {
		return err
	}

	rsp.Stock = h.serializeStock(&stock)
	return nil
}

// List returns all the stocks matching the UUIDs provided
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	var stocks []*storage.Stock
	var err error

	// Lookup by UUIDs
	if len(req.Uuids) > 0 {
		stocks, err = h.db.List(req.Uuids)
	}

	// Lookup by Symbols
	if len(req.Symbols) > 0 {
		stocks, err = h.db.ListBySymbol(req.Symbols)
	}

	// Lookup by Industries
	if req.Sector != "" {
		var industries []string

		for i, sector := range industryMap {
			if sector == req.Sector {
				industries = append(industries, i)
			}
		}

		stocks, err = h.db.ListByIndustries(industries)
	}

	// Check for an error
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.Stocks = make([]*proto.Stock, len(stocks))
	for i, p := range stocks {
		rsp.Stocks[i] = h.serializeStock(p)
	}

	return nil
}

// Search returns all the stocks matching the query provided
func (h *Handler) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.ListResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("INVALID_QUERY", "Query can't be blank")
	}

	limit := req.Limit
	if limit == 0 {
		limit = 30
	}

	stocks, err := h.db.Query(req.Query, limit)
	if err != nil {
		return err
	}

	rsp.Stocks = make([]*proto.Stock, len(stocks))
	for i, p := range stocks {
		rsp.Stocks[i] = h.serializeStock(p)
	}

	return nil
}

// All returns all the stocks
func (h *Handler) All(ctx context.Context, req *proto.AllRequest, rsp *proto.ListResponse) error {
	stocks, err := h.db.All()
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching stocks"}
		return err
	}

	rsp.Stocks = make([]*proto.Stock, len(stocks))
	for i, p := range stocks {
		rsp.Stocks[i] = h.serializeStock(p)
	}

	return nil
}

func (h *Handler) serializeStock(p *storage.Stock) *proto.Stock {
	return &proto.Stock{
		Uuid:             p.UUID,
		Name:             p.Name,
		Symbol:           p.Symbol,
		Exchange:         p.Exchange,
		Type:             p.Type,
		Region:           p.Region,
		Currency:         p.Currency,
		ProfilePictureId: p.ProfilePictureID,
		Sector:           industryMap[p.Industry],
		Industry:         p.Industry,
		Website:          p.Website,
		Description:      p.Description,
		Color:            p.Color,
	}
}

func (h *Handler) permittedParams(p *proto.Stock) storage.Stock {
	return storage.Stock{
		Name:             p.Name,
		Symbol:           p.Symbol,
		Exchange:         p.Exchange,
		Type:             p.Type,
		Region:           p.Region,
		Currency:         p.Currency,
		ProfilePictureID: p.ProfilePictureId,
		Industry:         p.Industry,
		Website:          p.Website,
		Description:      p.Description,
		Color:            p.Color,
	}
}

var industryMap = map[string]string{
	"Industrial Conglomerates":           "Industrials",
	"Oil & Gas Production":               "Energy",
	"Electric Utilities":                 "Utilities",
	"Electronics Distributors":           "Information Technology",
	"Personnel Services":                 "Industrials",
	"Motor Vehicles":                     "Consumer Discretionary",
	"Computer Communications":            "Communication Services",
	"Tools & Hardware":                   "Consumer Discretionary",
	"Discount Stores":                    "Consumer Discretionary",
	"Chemicals: Major Diversified":       "Materials",
	"Electronic Components":              "Information Technology",
	"Containers/Packaging":               "Materials",
	"Investment Trusts/Mutual Funds":     "Financials",
	"Specialty Telecommunications":       "Communication Services",
	"Internet Retail":                    "Consumer Discretionary",
	"Drugstore Chains":                   "Consumer Staples",
	"Coal":                               "Energy",
	"Data Processing Services":           "Information Technology",
	"Electronic Production Equipment":    "Information Technology",
	"Apparel/Footwear":                   "Consumer Discretionary",
	"Pharmaceuticals: Generic":           "Health Care",
	"Agricultural Commodities/Milling":   "Industrials",
	"Financial Conglomerates":            "Financials",
	"Wholesale Distributors":             "Consumer Discretionary",
	"Casinos/Gaming":                     "Consumer Discretionary",
	"Broadcasting":                       "Communication Services",
	"Apparel/Footwear Retail":            "Consumer Discretionary",
	"Investment Banks/Brokers":           "Financials",
	"Food: Specialty/Candy":              "Consumer Staples",
	"Forest Products":                    "Materials",
	"Office Equipment/Supplies":          "Industrials",
	"Advertising/Marketing Services":     "Communication Services",
	"Multi-Line Insurance":               "Financials",
	"Property/Casualty Insurance":        "Financials",
	"Construction Materials":             "Materials",
	"Real Estate Development":            "Real Estate",
	"Movies/Entertainment":               "Communication Services",
	"Electrical Products":                "Information Technology",
	"Hospital/Nursing Management":        "Health Care",
	"Electronic Equipment/Instruments":   "Information Technology",
	"Railroads":                          "Industrials",
	"Precious Metals":                    "Materials",
	"Chemicals: Agricultural":            "Materials",
	"Managed Health Care":                "Health Care",
	"Food: Meat/Fish/Dairy":              "Consumer Staples",
	"Other Metals/Minerals":              "Materials",
	"Household/Personal Care":            "Consumer Staples",
	"Beverages: Non-Alcoholic":           "Consumer Staples",
	"Major Banks":                        "Financials",
	"Chemicals: Specialty":               "Materials",
	"Cable/Satellite TV":                 "Communication Services",
	"Regional Banks":                     "Financials",
	"Industrial Specialties":             "Materials",
	"Medical/Nursing Services":           "Health Care",
	"Auto Parts: OEM":                    "Consumer Discretionary",
	"Catalog/Specialty Distribution":     "Industrials",
	"Major Telecommunications":           "Communication Services",
	"Commercial Printing/Forms":          "Industrials",
	"Publishing: Newspapers":             "Communication Services",
	"Industrial Machinery":               "Industrials",
	"Department Stores":                  "Consumer Discretionary",
	"Life/Health Insurance":              "Financials",
	"Wireless Telecommunications":        "Communication Services",
	"Trucks/Construction/Farm Machinery": "Industrials",
	"Textiles":                           "Consumer Discretionary",
	"Engineering & Construction":         "Industrials",
	"Automotive Aftermarket":             "Consumer Discretionary",
	"Electronics/Appliance Stores":       "Consumer Discretionary",
	"Telecommunications Equipment":       "Communication Services",
	"Steel":                              "Materials",
	"Pulp & Paper":                       "Materials",
	"Finance/Rental/Leasing":             "Financials",
	"Metal Fabrication":                  "Materials",
	"Miscellaneous":                      "",
	"Restaurants":                        "Consumer Discretionary",
	"Food Distributors":                  "Consumer Staples",
	"Home Furnishings":                   "Consumer Discretionary",
	"Real Estate Investment Trusts":      "Financials",
	"Electronics/Appliances":             "Consumer Discretionary",
	"Computer Peripherals":               "Information Technology",
	"Airlines":                           "Industrials",
	"Medical Distributors":               "Health Care",
	"Oilfield Services/Equipment":        "Energy",
	"Trucking":                           "Industrials",
	"Homebuilding":                       "Consumer Discretionary",
	"Specialty Insurance":                "Financials",
	"Recreational Products":              "Consumer Discretionary",
	"Biotechnology":                      "Health Care",
	"Food: Major Diversified":            "Consumer Staples",
	"Hotels/Resorts/Cruiselines":         "Consumer Discretionary",
	"Miscellaneous Commercial Services":  "Industrials",
	"Computer Processing Hardware":       "Information Technology",
	"Integrated Oil":                     "Energy",
	"Building Products":                  "Industrials",
	"Pharmaceuticals: Major":             "Health Care",
	"Other Consumer Specialties":         "Consumer Discretionary",
	"Oil Refining/Marketing":             "Energy",
	"Services to the Health Industry":    "Health Care",
	"Home Improvement Chains":            "Consumer Discretionary",
	"Gas Distributors":                   "Energy",
	"Alternative Power Generation":       "Energy",
	"Tobacco":                            "Consumer Staples",
	"Investment Managers":                "Financials",
	"Contract Drilling":                  "Energy",
	"Publishing: Books/Magazines":        "Communication Services",
	"Specialty Stores":                   "Consumer Discretionary",
	"Beverages: Alcoholic":               "Consumer Staples",
	"Food Retail":                        "Consumer Staples",
	"Aerospace & Defense":                "Industrials",
	"Media Conglomerates":                "Communication Services",
	"Internet Software/Services":         "Information Technology",
	"Pharmaceuticals: Other":             "Health Care",
	"Medical Specialties":                "Health Care",
	"Air Freight/Couriers":               "Industrials",
	"Aluminum":                           "Materials",
	"Other Transportation":               "Industrials",
	"Information Technology Services":    "Information Technology",
	"Packaged Software":                  "Information Technology",
	"Environmental Services":             "Industrials",
	"Insurance Brokers/Services":         "Financials",
	"Marine Shipping":                    "Industrials",
	"Other Consumer Services":            "Consumer Discretionary",
	"Consumer Sundries":                  "Consumer Discretionary",
	"Savings Banks":                      "Financials",
	"Water Utilities":                    "Utilities",
	"Financial Publishing/Services":      "Information Technology",
	"Oil & Gas Pipelines":                "Energy",
	"Semiconductors":                     "Information Technology",
	"Miscellaneous Manufacturing":        "Materials",
}
