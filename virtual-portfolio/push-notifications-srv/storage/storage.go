package storage

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the trades-srv db
type Service interface {
	CreateToken(string, string) (Token, error)
	GetToken(string) (Token, error)
	Close() error
}

// Token is a push notification token
type Token struct {
	gorm.Model

	UserUUID string `gorm:"type:uuid;"`
	Token    string `gorm:"type:varchar(100);"`
}

// BeforeSave performs the validations
func (t *Token) BeforeSave(db *gorm.DB) {
	if err := t.GetValidationError(db); err != nil {
		db.AddError(err)
	}
}

// GetValidationError checks for a validation error using the ozzo-validation package
func (t *Token) GetValidationError(db *gorm.DB) error {
	fmt.Println("GetValidationError")

	err := validation.ValidateStruct(t,
		validation.Field(&t.UserUUID, validation.Required),
		validation.Field(&t.Token, validation.Required))

	return err
}
