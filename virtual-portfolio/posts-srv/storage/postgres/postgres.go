package postgres

import (
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/posts-srv/storage"

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
	db.AutoMigrate(&storage.Post{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(post storage.Post) (storage.Post, error) {
	if errs := p.db.Create(&post).GetErrors(); len(errs) > 0 {
		return post, errs[0]
	}

	return post, nil
}

func (p postgres) Get(query storage.Post) (storage.Post, error) {
	var post storage.Post
	errs := p.db.Table("posts").Where(&query).First(&post).GetErrors()

	if len(errs) > 0 {
		return post, errs[0]
	}

	return post, nil
}

func (p postgres) Recent(limit, page int32) ([]*storage.Post, error) {
	var posts []*storage.Post

	query := p.db.Table("posts").Order("created_at DESC")
	query.Limit(limit).Offset(limit * page).Find(&posts)

	if errs := query.GetErrors(); len(errs) > 0 {
		return posts, errs[0]
	}

	return posts, nil
}

func (p postgres) Update(params storage.Post) (storage.Post, error) {
	post, err := p.Get(storage.Post{UUID: params.UUID})
	if err != nil {
		return params, err
	}

	req := p.db.Model(&post).Updates(mapWithoutBlank(params))
	if errs := req.GetErrors(); len(errs) > 0 {
		return post, errs[0]
	}

	return post, nil
}

func (p postgres) CountByUser(uuids []string, startTime, endTime time.Time) (map[string]int32, error) {
	var res []struct {
		UserUUID string
		Count    int32
	}

	q := p.db.Table("posts")
	q = q.Where("created_at >= ? AND created_at <= ?", startTime, endTime)
	q = q.Where("user_uuid IN (?)", uuids)
	q = q.Select("user_uuid, COUNT(*)")
	q = q.Group("user_uuid").Find(&res)

	if err := microgorm.TranslateErrors(q); err != nil {
		return map[string]int32{}, err
	}

	data := map[string]int32{}
	for _, r := range res {
		data[r.UserUUID] = r.Count
	}

	return data, nil
}

func (p postgres) Count(query storage.Post) (int32, error) {
	var count int32
	errs := p.db.Table("posts").Where(&query).Count(&count).GetErrors()

	if len(errs) > 0 {
		return count, errs[0]
	}

	return count, nil
}

func (p postgres) List(uuids []string) ([]*storage.Post, error) {
	var posts []*storage.Post

	query := p.db.Table("posts").Where("uuid IN (?)", uuids)
	query.Order("created_at DESC").Find(&posts)

	if errs := query.GetErrors(); len(errs) > 0 {
		return posts, errs[0]
	}

	return posts, nil
}

func (p postgres) ListFeed(feedType, feedUUID string) ([]*storage.Post, error) {
	feed := &storage.Post{FeedType: feedType, FeedUUID: feedUUID}

	var posts []*storage.Post
	query := p.db.Table("posts").Where(feed).Find(&posts)

	if errs := query.GetErrors(); len(errs) > 0 {
		return posts, errs[0]
	}

	return posts, nil
}

func (p postgres) ListUser(UUID string, limit, page int32) ([]*storage.Post, error) {
	var posts []*storage.Post

	query := p.db.Table("posts").Where(&storage.Post{UserUUID: UUID})
	query.Order("created_at DESC").Limit(limit).Offset(limit * page).Find(&posts)

	if errs := query.GetErrors(); len(errs) > 0 {
		return posts, errs[0]
	}

	return posts, nil
}

func (p postgres) Delete(query storage.Post) error {
	if errs := p.db.Where(&query).Delete(storage.Post{}).GetErrors(); len(errs) > 0 {
		return errs[0]
	}

	return nil
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
