package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/followers-srv/storage"
	"github.com/kytra-app/helpers/microgorm"

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
	db.AutoMigrate(&storage.Follow{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Follow(follower storage.Resource, followee storage.Resource) error {
	f := storage.Follow{
		FolloweeUUID: followee.UUID,
		FolloweeType: followee.Type,
		FollowerUUID: follower.UUID,
		FollowerType: follower.Type,
	}

	var existing int
	q := p.db.Model(&storage.Follow{}).Where(&f).Count(&existing)
	if err := microgorm.TranslateErrors(q); err != nil {
		return err
	} else if existing > 0 {
		return nil // Relationship already exists
	}

	return microgorm.TranslateErrors(p.db.Create(&f))
}

func (p postgres) Unfollow(follower storage.Resource, followee storage.Resource) error {
	query := storage.Follow{
		FolloweeUUID: followee.UUID,
		FolloweeType: followee.Type,
		FollowerUUID: follower.UUID,
		FollowerType: follower.Type,
	}

	var f storage.Follow

	q := p.db.Model(&storage.Follow{}).Where(&query).First(&f)
	if err := microgorm.TranslateErrors(q); err != nil {
		return err
	}

	q = p.db.Unscoped().Delete(&f)
	if err := microgorm.TranslateErrors(q); err != nil {
		return err
	}

	return nil
}

func (p postgres) GetFollowers(time time.Time, r storage.Resource) ([]*storage.Resource, error) {
	var res []*storage.Follow
	q := p.db.Table("follows")
	q = q.Where("followee_uuid = ? AND followee_type = ?", r.UUID, r.Type)
	q = q.Where("created_at < ?", time).Find(&res)

	if err := microgorm.TranslateErrors(q); err != nil {
		return make([]*storage.Resource, 0), err
	}

	followers := make([]*storage.Resource, len(res))
	for i, f := range res {
		followers[i] = &storage.Resource{UUID: f.FollowerUUID, Type: f.FollowerType}
	}

	return followers, nil
}

func (p postgres) GetFollowing(time time.Time, r storage.Resource) ([]*storage.Resource, error) {
	var res []*storage.Follow
	q := p.db.Table("follows")
	q = q.Where("follower_uuid = ? AND follower_type = ?", r.UUID, r.Type)
	q = q.Where("created_at < ?", time).Find(&res)

	if err := microgorm.TranslateErrors(q); err != nil {
		return make([]*storage.Resource, 0), err
	}

	following := make([]*storage.Resource, len(res))
	for i, f := range res {
		following[i] = &storage.Resource{UUID: f.FolloweeUUID, Type: f.FolloweeType}
	}

	return following, nil
}

func (p postgres) CountFollowers(time time.Time, r storage.Resource) (int32, error) {
	var count int32
	q := p.db.Table("follows")
	q = q.Where("followee_uuid = ? AND followee_type = ?", r.UUID, r.Type)
	q = q.Where("created_at < ?", time).Count(&count)
	return count, microgorm.TranslateErrors(q)
}

func (p postgres) CountFollowing(time time.Time, r storage.Resource) (int32, error) {
	var count int32
	q := p.db.Table("follows")
	q = q.Where("follower_uuid = ? AND follower_type = ?", r.UUID, r.Type)
	q = q.Where("created_at < ?", time).Count(&count)
	return count, microgorm.TranslateErrors(q)
}

func (p postgres) ListRelationships(time time.Time, follower storage.Resource, followeeType string, followeeUUIDs []string) ([]*storage.Resource, error) {
	baseQuery := storage.Follow{
		FollowerType: follower.Type,
		FollowerUUID: follower.UUID,
		FolloweeType: followeeType,
	}

	var res []*storage.Follow
	q := p.db.Model(&storage.Follow{}).Where(&baseQuery)
	q = q.Where("followee_uuid IN (?)", followeeUUIDs)
	q = q.Where("created_at < ?", time).Find(&res)
	if err := microgorm.TranslateErrors(q); err != nil {
		return make([]*storage.Resource, 0), nil
	}

	// Group the found follows relationships by followee UUID
	followsMap := make(map[string]*storage.Follow, len(res))
	for _, f := range res {
		followsMap[f.FolloweeUUID] = f
	}

	// Serialize the data
	rsp := make([]*storage.Resource, len(followeeUUIDs))
	for i, uuid := range followeeUUIDs {
		_, following := followsMap[uuid] // If no relationship is found, they're not following

		rsp[i] = &storage.Resource{
			UUID:      uuid,
			Type:      followeeType,
			Following: following,
		}
	}

	return rsp, nil
}
