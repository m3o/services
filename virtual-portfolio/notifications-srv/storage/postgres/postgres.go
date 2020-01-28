package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/notifications-srv/storage"

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
	db.AutoMigrate(&storage.Notification{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(n storage.Notification) (storage.Notification, error) {
	q := p.db.Create(&n)
	return n, microgorm.TranslateErrors(q)
}

func (p postgres) Get(query storage.Notification) (storage.Notification, error) {
	var post storage.Notification
	q := p.db.Table("notifications").Where(&query).First(&post)
	return post, microgorm.TranslateErrors(q)
}

func (p postgres) SetNotificationsSeen(UserUUID string) error {
	q := p.db.Table("notifications").Where(&storage.Notification{UserUUID: UserUUID, Seen: false})
	q.Update(&storage.Notification{Seen: true})
	return microgorm.TranslateErrors(q)
}

func (p postgres) List(query storage.Query) ([]*storage.Notification, error) {
	var notifications []*storage.Notification

	q := p.db.Table("notifications").Order("created_at DESC")

	// Add limits & page
	if query.Limit == 0 {
		query.Limit = 50
	}
	if query.Page == 0 {
		query.Page = 0
	}
	q = q.Limit(query.Limit).Offset(query.Limit * query.Page)

	// Add time constraints (if provided)
	if query.StartTime != nil {
		q = q.Where("created_at > ?", query.StartTime)
	}
	if query.EndTime != nil {
		q = q.Where("created_at < ?", query.EndTime)
	}

	// Scope to a single user, if requested
	if query.UserUUID != "" {
		q = q.Where("user_uuid = ?", query.UserUUID)
	}

	// Scope to only unseen, if requested
	if query.OnlyUnseen {
		q = q.Where("seen = ?", false)
	}

	q.Find(&notifications)
	return notifications, microgorm.TranslateErrors(q)
}
