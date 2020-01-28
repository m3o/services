package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the followers-srv db
type Service interface {
	Follow(Resource, Resource) error
	Unfollow(Resource, Resource) error
	GetFollowers(time.Time, Resource) ([]*Resource, error)
	GetFollowing(time.Time, Resource) ([]*Resource, error)
	CountFollowers(time.Time, Resource) (int32, error)
	CountFollowing(time.Time, Resource) (int32, error)
	ListRelationships(time.Time, Resource, string, []string) ([]*Resource, error)
	Close() error
}

// Resource is a resource which can be either a follower or a followee
type Resource struct {
	UUID      string
	Type      string
	Follower  bool
	Following bool
}

// Follow is a relationship between a follower and a followee
type Follow struct {
	gorm.Model

	FollowerUUID string `gorm:"type:varchar(100)"`
	FollowerType string `gorm:"type:varchar(100)"`
	FolloweeUUID string `gorm:"type:varchar(100)"`
	FolloweeType string `gorm:"type:varchar(100)"`
}

// BeforeSave performs the validations
func (follow *Follow) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(follow,
		validation.Field(&follow.FollowerUUID, validation.Required),
		validation.Field(&follow.FollowerType, validation.Required),
		validation.Field(&follow.FolloweeType, validation.Required),
		validation.Field(&follow.FolloweeType, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
