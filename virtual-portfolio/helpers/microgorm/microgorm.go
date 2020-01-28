package microgorm

import (
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/errors"
)

var (
	// ErrNotFound is a 404 Error
	ErrNotFound = errors.NotFound("NOT_FOUND", "Resource not found")
	// ErrDatabase is a 500 error
	ErrDatabase = errors.InternalServerError("DATABASE_ERROR", "A database error occured")
)

// TranslateErrors takes a pointer to a gorm database, gets the errors and transforms
// them into a single micro error which can be returned safely to a handler.
func TranslateErrors(db *gorm.DB) error {
	errs := db.GetErrors()

	if len(errs) == 0 {
		return nil
	}

	switch errs[0] {
	case gorm.ErrInvalidSQL:
	case gorm.ErrUnaddressable:
	case gorm.ErrCantStartTransaction:
		return ErrDatabase
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	}

	return errors.BadRequest("VALIDATION_FAILED", errs[0].Error())
}
