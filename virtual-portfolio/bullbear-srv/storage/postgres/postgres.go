package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/bullbear-srv/storage"

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

	db.AutoMigrate(&storage.Opinion{})

	return postgres{db}, nil
}

func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Get(r storage.Resource) (storage.Resource, error) {
	q := storage.Opinion{ResourceUUID: r.UUID, ResourceType: r.Type}

	var opinions []storage.Opinion
	if errs := p.db.Where(&q).Find(&opinions).GetErrors(); len(errs) > 0 {
		return storage.Resource{}, errs[0]
	}

	for _, o := range opinions {
		switch o.Opinion {
		case "BULLISH":
			r.BullsCount++
			r.Bulls = append(r.Bulls, o.UserUUID)
			break
		case "BEARISH":
			r.BearsCount++
			r.Bears = append(r.Bears, o.UserUUID)
			break
		}
	}

	return r, nil
}

func (p postgres) List(resourceType string, resourceUUIDs []string, userUUID string) ([]*storage.Resource, error) {
	query := p.db.Table("opinions").Where("resource_type = ? AND resource_uuid IN (?)", resourceType, resourceUUIDs)

	var options []storage.Opinion
	if errs := query.Find(&options).GetErrors(); len(errs) > 0 {
		return make([]*storage.Resource, 0), errs[0]
	}

	// Count Bulls & Bears for each resource
	bulls := make(map[string]int, len(resourceUUIDs))
	bears := make(map[string]int, len(resourceUUIDs))
	for _, o := range options {
		switch o.Opinion {
		case "BULLISH":
			bulls[o.ResourceUUID]++
			break
		case "BEARISH":
			bears[o.ResourceUUID]++
			break
		}
	}

	// Find the requesting users opinions
	userOpinions := make(map[string]string, len(resourceUUIDs))
	for _, o := range options {
		if o.UserUUID != userUUID {
			continue
		}
		userOpinions[o.ResourceUUID] = o.Opinion
	}

	// Serialize to Resource objects
	res := make([]*storage.Resource, len(resourceUUIDs))
	for i, uuid := range resourceUUIDs {
		res[i] = &storage.Resource{
			UUID:       uuid,
			Type:       resourceType,
			BullsCount: bulls[uuid],
			BearsCount: bears[uuid],
			Opinion:    userOpinions[uuid],
		}
	}

	return res, nil
}

func (p postgres) Create(params storage.Opinion) error {
	lookup := storage.Opinion{
		UserUUID:     params.UserUUID,
		ResourceType: params.ResourceType,
		ResourceUUID: params.ResourceUUID,
	}

	var opinion storage.Opinion
	var errors []error

	if p.db.Find(&opinion, &lookup).RecordNotFound() {
		errors = p.db.Create(&params).GetErrors()
	} else {
		opinion.Opinion = params.Opinion
		errors = p.db.Save(&opinion).GetErrors()
	}

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
