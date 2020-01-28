package storage

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Service is an wrapper of the user-srv db
type Service interface {
	Count() (int32, error)
	Create(User) (User, error)
	Find(User) (User, error)
	Update(User) (User, error)
	List([]string) ([]User, error)
	ListByPhoneNumber([]string) ([]User, error)
	Query(string, int32) ([]User, error)
	All() ([]User, error)
	Close() error
}

// User is an individual who registered to Kytra
type User struct {
	gorm.Model

	UUID             string `gorm:"type:uuid;primary_key;"`
	Username         string `gorm:"type:varchar(100)"`
	FirstName        string `gorm:"type:varchar(100)"`
	LastName         string `gorm:"type:varchar(100)"`
	Email            string `gorm:"type:varchar(100)"`
	PhoneNumber      string `gorm:"type:varchar(15)"`
	Password         string `gorm:"type:varchar(250)"`
	ProfilePictureID string `gorm:"type:varchar(100)"`
	Admin            bool   `gorm:"type:boolean;DEFAULT=FALSE"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (user *User) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(user,
		validation.Field(&user.FirstName, validation.Required, validation.Length(2, 30)),
		validation.Field(&user.LastName, validation.Required, validation.Length(2, 30)),
		validation.Field(&user.Username, validation.Required, validation.Length(2, 30)),
		validation.Field(&user.Email, validation.Required, is.Email),
		// validation.Field(&user.PhoneNumber, validation.Required),
		validation.Field(&user.Password, validation.Required))

	if err != nil {
		db.AddError(err)
		return
	}

	var userWithEmail User
	db.Where(User{Email: user.Email}).First(&userWithEmail)
	if userWithEmail.UUID != user.UUID && userWithEmail.UUID != "" {
		db.AddError(errors.New("Email has already been taken"))
		return
	}

	// Conditional while the web flow is still active
	if user.PhoneNumber != "" {
		var userWithPhone User
		db.Where(User{PhoneNumber: user.PhoneNumber}).First(&userWithPhone)
		if userWithPhone.UUID != user.UUID && userWithPhone.UUID != "" {
			db.AddError(errors.New("Phone number has already been taken"))
			return
		}
	}

	var userWithUsername User
	query := db.Table("users").Where("lower(username) = ?", strings.ToLower(user.Username))
	query.First(&userWithUsername)
	if userWithUsername.UUID != user.UUID && userWithUsername.UUID != "" {
		db.AddError(errors.New("Username has already been taken"))
		return
	}

	containsSpecialChar := func(r rune) bool {
		return (r < 'A' || r > 'z') && (r < '0' || r > '9')
	}
	if strings.IndexFunc(user.Username, containsSpecialChar) != -1 {
		db.AddError(errors.New("Username contains an invalid character"))
	}
}
