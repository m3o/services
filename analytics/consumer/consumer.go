package consumer

import (
	"fmt"
	"os"
	"strings"
	"time"

	customers "github.com/m3o/services/customers/proto"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	groupName = "analytics"
	defaultPG = "host=localhost user=postgres dbname=postgres sslmode=disable password=password"
)

// Consumer will subscribe to events and then store them in a database
type Consumer struct {
	db      *gorm.DB
	ErrChan chan error
}

type eventHandler func(events.Event) error

// Init the consumer, connects to the database
func (c *Consumer) Init() error {
	// load the database address from config
	dsn := os.Getenv("MICRO_STORE_ADDRESS")
	if len(dsn) == 0 {
		dsn = defaultPG
	}

	// connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Error connecting to postgres: %v", err)
	}
	c.db = db

	// migrate the schema
	if err := db.AutoMigrate(&customerEvent{}); err != nil {
		return fmt.Errorf("Error migrating db schema: %v", err)
	}

	return nil
}

// Run will subscribe to the event queues and process the events async
func (c *Consumer) Run() error {
	// subscribe to the events
	if err := c.subscribeToTopic(customers.EventsTopic, c.handleCustomerEvent); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) subscribeToTopic(topic string, handler eventHandler) error {
	evChan, err := events.Consume(customers.EventsTopic,
		events.WithGroup(groupName),
		events.WithOffset(time.Unix(0, 0)),
	)
	if err != nil {
		return fmt.Errorf("Error subscribing to %v topic: %v", customers.EventsTopic, err)
	}

	go func() {
		for ev := range evChan {
			if err := handler(ev); err != nil {
				ev.Nack()
				logger.Errorf("Error processing event #%v on topic %v: %v", ev.ID, topic, err)
				continue
			}
			ev.Ack()
		}

		// exit the application
		if c.ErrChan != nil {
			c.ErrChan <- fmt.Errorf("Stopped reading from topic: %v", topic)
		}
	}()

	return nil
}

type customerEvent struct {
	EventID       string              `gorm:"primaryKey"`
	EventCallerID string              `json:"caller_id"`
	EventType     customers.EventType `json:"type"`
	Customer      customers.Customer  `gorm:"embedded"`
}

func (c *Consumer) handleCustomerEvent(msg events.Event) error {
	var ev customerEvent
	if err := msg.Unmarshal(&ev); err != nil {
		return err
	}
	ev.EventID = msg.ID

	err := c.db.Create(&ev).Error
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil
	}
	return err
}
