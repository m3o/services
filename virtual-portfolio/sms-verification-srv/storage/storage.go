package storage

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	"github.com/kytra-app/sms-verification-srv/helpers"
	"github.com/micro/go-micro/errors"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrInvalidCode means the code was incorrect
	ErrInvalidCode = errors.BadRequest("INVALID_CODE", "The code is incorrect")

	// ErrInvalidNumber means the phone number was invalid
	ErrInvalidNumber = errors.BadRequest("INVALID_NUMBER", "The phone number is invalid")

	// ErrNoCodeRequested means no code was requested for this phone number
	ErrNoCodeRequested = errors.BadRequest("NO_CODE_REQUESTED", "No code has been requested for this phone number")

	// ErrCodeExpired means the code generated is no longer valid
	ErrCodeExpired = errors.BadRequest("CODE_EXPIRED", "The code has expired")
)

// Service is an wrapper of the trades-srv db
type Service interface {
	Get(string) (Verification, error)
	Request(string) (Verification, error)
	Verify(string, string) (Verification, error)
	Close() error
}

// Verification is an SMS Verificiaton
type Verification struct {
	gorm.Model

	UUID        string `gorm:"type:uuid;primary_key;"`
	PhoneNumber string `gorm:"type:varchar(20);"`
	Code        string `gorm:"type:varchar(6);"`
	Attempts    int    `gorm:"type:int;"`
	Verified    bool   `gorm:"type:bool;"`
}

// Expired returns a bool (has the verification expired?)
func (v *Verification) Expired() bool {
	if v.Attempts >= 3 {
		return true
	}

	maxTime := time.Duration(15 * time.Minute)
	if v.CreatedAt.Unix() < time.Now().Add(-maxTime).Unix() {
		return true
	}

	return false
}

// BeforeSave performs the validations
func (v *Verification) BeforeSave(db *gorm.DB) {
	if v.UUID == "" {
		v.UUID = uuid.NewV4().String()
	}

	if v.Code == "" {
		v.Code = helpers.GenerateCode()
	}

	if err := v.GetValidationError(db); err != nil {
		db.AddError(err)
	}
}

// GetValidationError checks for a validation error using the ozzo-validation package
func (v *Verification) GetValidationError(db *gorm.DB) error {
	fmt.Println("GetValidationError")

	err := validation.ValidateStruct(v,
		validation.Field(&v.UUID, validation.Required),
		validation.Field(&v.Code, validation.Required),
		validation.Field(&v.PhoneNumber, validation.Required))

	return err
}
