package postgres

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/users-srv/storage"

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

	db.AutoMigrate(&storage.User{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(user storage.User) (storage.User, error) {
	req := p.db.Create(&user)
	return user, microgorm.TranslateErrors(req)
}

func (p postgres) Count() (int32, error) {
	var count int32
	req := p.db.Table("users").Count(&count)
	return count, microgorm.TranslateErrors(req)
}

func (p postgres) Find(query storage.User) (storage.User, error) {
	var user storage.User
	req := p.db.Where(query).First(&user)
	return user, microgorm.TranslateErrors(req)
}

func (p postgres) Update(params storage.User) (storage.User, error) {
	user, err := p.Find(storage.User{UUID: params.UUID})
	if err != nil {
		return params, err
	}

	req := p.db.Model(&user).Updates(mapWithoutBlank(params))
	return user, microgorm.TranslateErrors(req)
}

func (p postgres) List(uuids []string) ([]storage.User, error) {
	var users []storage.User
	req := p.db.Table("users").Where("uuid IN (?)", uuids).Find(&users)
	return users, microgorm.TranslateErrors(req)
}

func (p postgres) ListByPhoneNumber(numbers []string) ([]storage.User, error) {
	var users []storage.User
	req := p.db.Table("users").Where("phone_number IN (?)", numbers).Limit(50).Find(&users)
	return users, microgorm.TranslateErrors(req)
}

func (p postgres) All() ([]storage.User, error) {
	var users []storage.User
	req := p.db.Table("users").Find(&users)
	return users, microgorm.TranslateErrors(req)
}

func (p postgres) Query(query string, limit int32) ([]storage.User, error) {
	q := fmt.Sprintf("lower(username) LIKE lower('%%%v%%') OR lower(concat(first_name,' ',last_name)) LIKE lower('%%%v%%')",
		query, query,
	)

	var users []storage.User
	req := p.db.Table("users").Where(q).Limit(limit).Find(&users)
	return users, microgorm.TranslateErrors(req)
}

func mapWithoutBlank(data interface{}) map[string]interface{} {
	in := structs.Map(data)
	out := make(map[string]interface{}, len(in))

	for key, val := range in {
		if val != "" && val != nil {
			out[key] = val
		}
	}

	return out
}
