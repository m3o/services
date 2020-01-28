package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/push-notifications-srv/storage"

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
	db.AutoMigrate(&storage.Token{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) CreateToken(token, userUUID string) (storage.Token, error) {
	t := storage.Token{Token: token, UserUUID: userUUID}
	q := p.db.Create(&t)
	return t, microgorm.TranslateErrors(q)
}

func (p postgres) GetToken(userUUID string) (storage.Token, error) {
	token := storage.Token{UserUUID: userUUID}
	query := p.db.Where(&token).Order("created_at DESC").First(&token)
	return token, microgorm.TranslateErrors(query)
}
