package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/sms-verification-srv/storage"

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
	db.AutoMigrate(&storage.Verification{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Get(uuid string) (res storage.Verification, err error) {
	query := p.db.Where(&storage.Verification{UUID: uuid}).First(&res)
	return res, microgorm.TranslateErrors(query)
}
func (p postgres) Request(phoneNumber string) (storage.Verification, error) {
	ver := storage.Verification{PhoneNumber: phoneNumber}
	return ver, microgorm.TranslateErrors(p.db.Create(&ver))
}

func (p postgres) Verify(phoneNumber, code string) (storage.Verification, error) {
	// Lookup verification usng phone number
	ver := storage.Verification{PhoneNumber: phoneNumber, Verified: false}
	query := p.db.Where(&ver).Order("created_at DESC").First(&ver)
	if err := microgorm.TranslateErrors(query); err != nil {
		return ver, storage.ErrNoCodeRequested
	}

	if ver.Expired() {
		return ver, storage.ErrCodeExpired
	}

	updateQueryBase := p.db.Table("verifications").Where(&storage.Verification{UUID: ver.UUID})

	// Ensure the code matched
	if ver.Code != code {
		updateQueryBase.Update(&storage.Verification{Attempts: ver.Attempts + 1})
		return ver, storage.ErrInvalidCode
	}

	// Update the verification status to verified
	query = updateQueryBase.Update(&storage.Verification{
		Attempts: ver.Attempts + 1,
		Verified: true,
	})
	if err := microgorm.TranslateErrors(query); err != nil {
		return ver, err
	}

	return ver, nil
}
