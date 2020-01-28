package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/ledger-srv/storage"
	"github.com/micro/go-micro/errors"

	// The PG driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type postgres struct{ db *gorm.DB }

// New returns an instance of storage.Service backed by a PG database
func New(host, dbname, user, password string) (storage.Service, error) {
	// Construct the DB Source
	src := fmt.Sprintf("host=%v dbname=%v user=%v password=%v port=5432 sslmode=disable",
		host, dbname, user, password)

	// Connect to the DB
	db, err := gorm.Open("postgres", src)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)
	db.AutoMigrate(&storage.Transaction{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) CreateTransaction(t storage.Transaction) (storage.Transaction, error) {
	if err := p.ensureBalanceForTransaction(t); err != nil {
		return t, err
	}

	req := p.db.Create(&t)
	return t, microgorm.TranslateErrors(req)
}

func (p postgres) GetPortfolioBalance(time time.Time, portfolioUUID string) (int64, error) {
	positive, err := p.sumTransactions(time, portfolioUUID, storage.PositiveTransactionTypes)
	if err != nil {
		return 0, err
	}

	negative, err := p.sumTransactions(time, portfolioUUID, storage.NegativeTransactionTypes)
	if err != nil {
		return 0, err
	}

	return positive - negative, nil
}

// sumTransactions finds the sum of the amounts, given a portfolioUUID and a slice of transaction types
func (p postgres) sumTransactions(time time.Time, portfolioUUID string, types []string) (int64, error) {
	var result struct{ Value int64 }

	query := p.db.Table("transactions").Select("SUM(amount) as Value")
	query = query.Where("type IN (?) AND portfolio_uuid=?", types, portfolioUUID)
	query = query.Where("created_at < ?", time)

	query.Scan(&result)
	return result.Value, microgorm.TranslateErrors(query)
}

func (p postgres) ensureBalanceForTransaction(transaction storage.Transaction) error {
	// Don't check the balance if the transaction will be incrementing it
	for _, t := range storage.PositiveTransactionTypes {
		if transaction.Type == t {
			return nil
		}
	}

	balance, err := p.GetPortfolioBalance(time.Now(), transaction.PortfolioUUID)
	if err != nil {
		return err
	}
	if balance < transaction.Amount {
		return errors.BadRequest("INSUFFICIENT_FUNDS", "You do not have enough funds to make this transaction")
	}

	return nil
}
