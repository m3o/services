package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	pb "github.com/m3o/services/mixpanel/proto"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
)

type Mixpanel struct {
	client *MixpanelClient
}

type conf struct {
	User    string
	Secret  string
	Project string
}

func NewHandler() *Mixpanel {
	c, err := config.Get("micro.mixpanel")
	if err != nil {
		logger.Fatalf("Failed to load config %s", err)
	}
	var cObj conf
	if err := c.Scan(&cObj); err != nil {
		logger.Fatalf("Failed to load config %s", err)
	}
	m := &Mixpanel{
		client: &MixpanelClient{
			User:    cObj.User,
			Pass:    cObj.Secret,
			Project: cObj.Project,
		},
	}
	go m.consumeEvents()
	return m
}

type EventData struct {
	Event string `json:"event"`
	// distinct_id property tracks a user
	// token property maps your mixpanel project
	// time property is unix timestamp in secs
	// $insert_id is the UUID of the event
	Properties map[string]interface{} `json:"properties"`
}

type Event struct {
	Data EventData `json:"data"`
}

type MixpanelClient struct {
	User    string
	Pass    string
	Project string
}

func (m *MixpanelClient) newMixpanelEvent(topic, typeStr, customerID, evtID string, ts int64, data interface{}) Event {
	mev := Event{
		Data: EventData{
			Event: fmt.Sprintf("%s_%s", topic, typeStr),
			Properties: map[string]interface{}{
				"token":       m.Project,
				"time":        ts,
				"distinct_id": customerID,
				"$insert_id":  evtID,
				"data":        data,
			},
		},
	}
	return mev
}

func (m *MixpanelClient) Track(ev Event) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	logger.Infof("Tracking %s", string(b))
	req, err := http.NewRequest(http.MethodPost, "https://api.mixpanel.com/track#live-event-deduplicate", bytes.NewReader(b))
	if err != nil {
		logger.Errorf("Error creating http req %s", err)
		return err
	}
	req.SetBasicAuth(m.User, m.Pass)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Error creating http req %s", err)
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode > 299 || rsp.StatusCode < 200 {
		logger.Errorf("Error creating http req %s", err)
		return err
	}
	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		logger.Errorf("Error reading http rsp %s %s", err, string(b))
		return nil // ignore
	}
	logger.Infof("Response %s", string(b))
	return nil
}

func (b *Mixpanel) Ping(ctx context.Context, request *pb.PingRequest, response *pb.PingResponse) error {
	return nil
}
