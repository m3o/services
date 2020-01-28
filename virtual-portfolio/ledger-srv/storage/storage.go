package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// PositiveTransactionTypes increment the balance on the ledger
var PositiveTransactionTypes = []string{"DEPOSIT", "SELL_ASSET"}

// NegativeTransactionTypes decrement the balance on the ledger
var NegativeTransactionTypes = []string{"WITHDRAWAL", "BUY_ASSET"}

// TransactionTypes is both positive and negative transaction types
var TransactionTypes = append(PositiveTransactionTypes, NegativeTransactionTypes...)

// Service is an wrapper of the ledger-srv db
type Service interface {
	CreateTransaction(Transaction) (Transaction, error)
	GetPortfolioBalance(time.Time, string) (int64, error)
	Close() error
}

// Transaction is a entry to the ledger
type Transaction struct {
	gorm.Model

	UUID          string `gorm:"type:uuid;primary_key;"`
	PortfolioUUID string `gorm:"type:uuid;"`
	Type          string `gorm:"type:varchar(100)"`
	Amount        int64  `gorm:"type:integer"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (t *Transaction) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (t *Transaction) BeforeSave(db *gorm.DB) {
	// validation.In requires a slice of interfaces, not strings :(
	ttInterfaces := make([]interface{}, len(TransactionTypes))
	for i, v := range TransactionTypes {
		ttInterfaces[i] = v
	}

	err := validation.ValidateStruct(t,
		validation.Field(&t.PortfolioUUID, validation.Required),
		validation.Field(&t.Amount, validation.Required, validation.Min(0)),
		validation.Field(&t.Type, validation.Required, validation.In(ttInterfaces...)))

	if err != nil {
		db.AddError(err)
	}
}
