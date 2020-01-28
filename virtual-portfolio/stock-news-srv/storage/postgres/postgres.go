package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/stock-news-srv/storage"

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
	db.AutoMigrate(&storage.Article{})
	db.Model(&storage.Article{}).AddUniqueIndex("idx_article_url_stock_uuid", "article_url", "stock_uuid")

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

// Create inserts a new article into the database
func (p postgres) Create(article storage.Article) (storage.Article, error) {
	req := p.db.Create(&article)
	return article, microgorm.TranslateErrors(req)
}

// ListForStock retrieves all articles which occured on the date provided
func (p postgres) ListForStock(t time.Time, stockUUIDs []string) ([]*storage.Article, error) {
	startTime := t.Truncate(time.Hour * 24)

	var articles []*storage.Article
	req := p.db.Table("articles").Where("created_at >= ? AND created_at <= ?", startTime, t)
	req = req.Where("stock_uuid IN (?)", stockUUIDs)
	req.Find(&articles)

	return articles, microgorm.TranslateErrors(req)
}

// ListForMarket retrieves all articles which occured on the date provided
func (p postgres) ListForMarket(t time.Time) ([]*storage.Article, error) {
	startTime := t.Truncate(time.Hour * 24)

	var articles []*storage.Article
	req := p.db.Table("articles").Where("created_at >= ? AND created_at <= ?", startTime, t)
	req = req.Where("stock_uuid IS NULL OR stock_uuid = ?", "")
	req.Find(&articles)

	return articles, microgorm.TranslateErrors(req)
}
