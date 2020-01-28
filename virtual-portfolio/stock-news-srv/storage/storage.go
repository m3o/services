package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the stock-news-srv db
type Service interface {
	Create(Article) (Article, error)
	ListForStock(time.Time, []string) ([]*Article, error)
	ListForMarket(time.Time) ([]*Article, error)
	Close() error
}

// Article represents a piece of news
type Article struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	ArticleURL  string `gorm:"unique_index:idx_article_url_stock_uuid" json:"article_url"`
	Title       string `json:"title"`
	Source      string `json:"source"`
	Description string `gorm:"type:text" json:"Description"`
	ImageURL    string `json:"image_url"`
	StockUUID   string `json:"stock_uuid"`
}

// BeforeSave performs the validations
func (m *Article) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(m,
		validation.Field(&m.ArticleURL, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Source, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
