package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/errors"
	uuid "github.com/satori/go.uuid"
)

// TradeTypes is the valid values for the type attribute of trade
var TradeTypes = []string{"BUY", "SELL"}

// Service is an wrapper of the trades-srv db
type Service interface {
	PrevalidateTrade(Trade) error
	CreateTrade(Trade) (Trade, error)
	GetTrade(Trade) (Trade, error)
	SetTradeMetadata(Trade) (Trade, error)
	ListTrades(time.Time, time.Time, []string) ([]Trade, error)
	ListTradesForPosition(Position, time.Time) ([]Trade, error)
	ListTradesForPortfolio(string, time.Time) ([]Trade, error)
	ListPositionsForPortfolio(string, time.Time) ([]Position, error)
	ListPositions([]string, string, []string, time.Time) ([]Position, error)
	AllAssets() ([]Asset, error)
	Close() error
}

// Trade is a transaction, buying or selling an asset, e.g. BUY 1000 units of Apple Inc at $22.24 a share
type Trade struct {
	gorm.Model

	UUID          string `gorm:"type:uuid;primary_key;"`
	ClientUUID    string `gorm:"type:uuid;"`
	Type          string `gorm:"type:varchar(25);"`
	AssetUUID     string `gorm:"type:uuid;"`
	AssetType     string `gorm:"type:varchar(25);"`
	PortfolioUUID string `gorm:"type:uuid;"`
	Quantity      int64  `gorm:"type:integer"`
	UnitPrice     int64  `gorm:"type:integer"`
	TargetPrice   int64  `gorm:"type:integer"`
	Notes         string `gorm:"type:text"`
}

// Asset is a resource which can be trades
type Asset struct {
	UUID string
	Type string
}

// Position is a summation of trades, e.g. if you buy 10 shares then sell 6, your position will be 4 shares.
type Position struct {
	AssetUUID, AssetType string
	PortfolioUUID        string
	Quantity             int64
	BookCost             int64
}

// BeforeCreate will set a UUID rather than numeric ID.
func (t *Trade) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (t *Trade) BeforeSave(db *gorm.DB) {
	if err := t.GetValidationError(db); err != nil {
		db.AddError(err)
	}
}

// GetValidationError checks for a validation error using the ozzo-validation package
func (t *Trade) GetValidationError(db *gorm.DB) error {
	// validation.In requires a slice of interfaces, not strings :(
	ttInterfaces := make([]interface{}, len(TradeTypes))
	for i, v := range TradeTypes {
		ttInterfaces[i] = v
	}

	err := validation.ValidateStruct(t,
		validation.Field(&t.ClientUUID, validation.Required),
		validation.Field(&t.AssetUUID, validation.Required),
		validation.Field(&t.AssetType, validation.Required),
		validation.Field(&t.PortfolioUUID, validation.Required),
		validation.Field(&t.Quantity, validation.Required, validation.Min(0)),
		validation.Field(&t.UnitPrice, validation.Required, validation.Min(0)),
		validation.Field(&t.Type, validation.Required, validation.In(ttInterfaces...)))

	if err != nil {
		return err
	}

	// Ensure the Client UUID has not been taken. This prevents the front-end submitting the same request twice.
	var count int
	q := db.Table("trades").Where(Trade{PortfolioUUID: t.PortfolioUUID, ClientUUID: t.ClientUUID})

	// exclude the current record, allows validation to pass on update
	if t.UUID != "" {
		q.Where("uuid != ?", t.UUID).Count(&count)
	} else {
		q.Count(&count)
	}

	if count > 0 {
		return errors.BadRequest("INVALID_CLIENT_UUID", "This client_uuid has already been taken")
	}

	return nil
}
