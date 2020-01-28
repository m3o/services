package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/feeditems-srv/storage"

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

	db.AutoMigrate(&storage.FeedItem{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(item storage.FeedItem) (storage.FeedItem, error) {
	if errs := p.db.Create(&item).GetErrors(); len(errs) > 0 {
		return item, errs[0]
	}

	return item, nil
}

func (p postgres) Get(query storage.FeedItem) (storage.FeedItem, error) {
	var item storage.FeedItem
	errs := p.db.Table("feed_items").Where(&query).First(&item).GetErrors()

	if len(errs) > 0 {
		return item, errs[0]
	}

	return item, nil
}

func (p postgres) GetFeed(feedType, feedUUID string, page, limit int32) ([]*storage.FeedItem, error) {
	var items []*storage.FeedItem

	q := storage.FeedItem{FeedUUID: feedUUID, FeedType: feedType}
	query := p.db.Table("feed_items").Order("created_at DESC").Where(&q)
	query.Limit(limit).Offset(limit * page).Find(&items)

	if errs := query.GetErrors(); len(errs) > 0 {
		return items, errs[0]
	}

	return items, nil
}

func (p postgres) Delete(query storage.FeedItem) error {
	if errs := p.db.Delete(&query).GetErrors(); len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (p postgres) BulkDelete(feedType, feedUUID string, postUUIDs []string) error {
	query := p.db.Where(&storage.FeedItem{FeedType: feedType, FeedUUID: feedUUID})
	query = query.Where("post_uuid IN (?)", postUUIDs)

	if errs := query.Delete(storage.FeedItem{}).GetErrors(); len(errs) > 0 {
		return errs[0]
	}

	return nil
}
