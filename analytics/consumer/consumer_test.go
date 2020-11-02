package consumer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	customers "github.com/m3o/services/customers/proto"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/config/env"
	"github.com/micro/micro/v3/service/events"
)

func TestHandleCustomerEvent(t *testing.T) {
	// load a noop test to prevent erorrs
	config.DefaultConfig, _ = env.NewConfig()

	// setup the consumer and connect to the db
	c := &Consumer{}
	assert.NoErrorf(t, c.Init(), "Init should not error")
	c.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&customerEvent{})

	// construct the event payload
	payload := &customers.Event{
		Type: customers.EventType_EventTypeCreated,
		Customer: &customers.Customer{
			Id:      uuid.New().String(),
			Status:  "active",
			Created: time.Now().Unix(),
			Email:   "foo@bar.com",
		},
		CallerId: uuid.New().String(),
	}
	bytes, err := json.Marshal(payload)
	assert.NoErrorf(t, err, "Error marshaling object")

	// handle the event
	ev := events.Event{
		ID:      uuid.New().String(),
		Payload: bytes,
	}
	assert.NoError(t, c.handleCustomerEvent(ev))

	// load the object from the database
	var res customerEvent
	err = c.db.First(&res).Error
	assert.NoErrorf(t, err, "There was no event written to the database")

	// compare the object
	assert.Equal(t, ev.ID, res.EventID)
	assert.Equal(t, payload.Type, res.EventType)
	assert.Equal(t, payload.Customer.Id, res.Customer.Id)
	assert.Equal(t, payload.CallerId, res.EventCallerID)

	// process the event again, it should not create a second record
	assert.NoError(t, c.handleCustomerEvent(ev))
	var count int64
	assert.NoError(t, c.db.Model(&customerEvent{}).Count(&count).Error)
	assert.EqualValues(t, 1, count, "Only one record should exist")
}
