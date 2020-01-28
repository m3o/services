package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/comments-srv/storage"

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

	db.AutoMigrate(&storage.Comment{})

	return postgres{db}, nil
}

func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Get(q storage.Comment) (storage.Comment, error) {
	var comment storage.Comment

	if errs := p.db.Where(q).First(&comment).GetErrors(); len(errs) > 0 {
		return comment, errs[0]
	}

	return comment, nil
}

func (p postgres) Create(comment storage.Comment) (storage.Comment, error) {
	if errs := p.db.Save(&comment).GetErrors(); len(errs) > 0 {
		return comment, errs[0]
	}

	return comment, nil
}

func (p postgres) Delete(uuid string) error {
	if errs := p.db.Delete(&storage.Comment{UUID: uuid}).GetErrors(); len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (p postgres) GetResource(r storage.Resource) (storage.Resource, error) {
	q := storage.Comment{ResourceUUID: r.UUID, ResourceType: r.Type}

	if errs := p.db.Where(&q).Find(&r.Comments).GetErrors(); len(errs) > 0 {
		return storage.Resource{}, errs[0]
	}

	return r, nil
}

func (p postgres) ListResources(resourceType string, resourceUUIDs []string) ([]*storage.Resource, error) {
	query := p.db.Table("comments").Where("resource_type = ? AND resource_uuid IN (?)", resourceType, resourceUUIDs)

	var comments []storage.Comment
	if errs := query.Find(&comments).GetErrors(); len(errs) > 0 {
		return make([]*storage.Resource, 0), errs[0]
	}

	// Group comments by resource_uuid
	cMap := make(map[string][]storage.Comment, len(resourceUUIDs))
	for _, c := range comments {
		cMap[c.ResourceUUID] = append(cMap[c.ResourceUUID], c)
	}

	// Serialize to Resource objects
	res := make([]*storage.Resource, len(resourceUUIDs))
	for i, uuid := range resourceUUIDs {
		res[i] = &storage.Resource{
			UUID:     uuid,
			Type:     resourceType,
			Comments: cMap[uuid],
		}
	}

	return res, nil
}
