package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/insights-srv/storage"

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
	db.AutoMigrate(&storage.Insight{})
	db.AutoMigrate(&storage.UserView{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

// CreateInsight insets a new insight into the DB
func (p postgres) CreateInsight(insight storage.Insight) (storage.Insight, error) {
	req := p.db.Create(&insight)
	return insight, microgorm.TranslateErrors(req)
}

// CreateUserView insets a new user view into the DB
func (p postgres) CreateUserView(view storage.UserView) (storage.UserView, error) {
	req := p.db.Create(&view)
	return view, microgorm.TranslateErrors(req)
}

// ListAssets returns a slice of assets which have insights on the given date
func (p postgres) ListAssets(date time.Time, excludeNews bool) ([]storage.Asset, error) {
	endTime := date.Add(time.Hour * 24)

	var insights []storage.Insight
	query := p.db.Table("insights")
	query = query.Where("created_at > ? AND created_at < ?", date, endTime)
	query = query.Select("DISTINCT(asset_uuid), asset_type, created_at, type")
	if excludeNews {
		query = query.Where("type != 'NEWS'")
	}

	if err := microgorm.TranslateErrors(query.Find(&insights)); err != nil {
		return []storage.Asset{}, err
	}

	assets := make([]storage.Asset, len(insights))
	for i, insight := range insights {
		assets[i] = storage.Asset{UUID: insight.AssetUUID, Type: insight.AssetType}
	}
	return assets, nil

}

// ListInsightsForUser returns a slice of insights for the specific user on a given date
func (p postgres) ListInsightsForUser(userUUID string, date time.Time) ([]storage.Insight, error) {
	endTime := date.Add(time.Hour * 24)

	var insights []storage.Insight
	query := p.db.Table("insights")
	query = query.Where("created_at > ? AND created_at < ?", date, endTime)
	query = query.Where("user_uuid = ? OR user_uuid IS NULL or user_uuid = ''", userUUID)

	err := microgorm.TranslateErrors(query.Find(&insights))
	return insights, err
}

// ListInsightsForAssets returns a slice of insights for any of the given assets on a given date
// which are not specific to any one user
// func (p postgres) ListInsightsForAssets(assetType string, assetUUIDs []string, date time.Time) ([]storage.Insight, error) {
// 	endTime := date.Add(time.Hour * 24)

// 	var insights []storage.Insight
// 	query := p.db.Table("insights")
// 	query = query.Where("user_uuid = '' OR user_uuid IS NULL")
// 	query = query.Where("created_at > ? AND created_at < ?", date, endTime)
// 	query = query.Where("asset_type = ? AND asset_uuid IN (?)", assetType, assetUUIDs)

// 	err := microgorm.TranslateErrors(query.Find(&insights))
// 	return insights, err
// }

// GetUserView returns the last UserView for the given date
func (p postgres) GetUserView(userUUID string, asset storage.Asset, date time.Time) (storage.UserView, error) {
	endTime := date.Add(time.Hour * 24)

	var view storage.UserView
	query := p.db.Table("user_views")
	query = query.Order("created_at DESC")
	query = query.Where("user_uuid = ?", userUUID)
	query = query.Where("created_at > ? AND created_at < ?", date, endTime)
	query = query.Where("asset_uuid = ? AND asset_type = ?", asset.UUID, asset.Type)

	err := microgorm.TranslateErrors(query.First(&view))
	return view, err
}
