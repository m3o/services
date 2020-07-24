package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/micro/go-micro/v2/config"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"

	"github.com/google/uuid"
	alert "github.com/m3o/services/alert/proto/alert"
	"github.com/slack-go/slack"
)

const (
	storePrefixEvents = "events/"
)

type Alert struct {
	store       store.Store
	slackClient *slack.Client
}

type event struct {
	ID       string            `json:"id"`
	Category string            `json:"category"`
	Action   string            `json:"action"`
	Label    string            `json:"label"`
	Value    uint64            `json:"value"`
	Metadata map[string]string `json:"metadata"`
}

func NewAlert(store store.Store) *Alert {
	slackToken := config.Get("micro", "alert", "slack_token").String("")
	if len(slackToken) == 0 {
		log.Errorf("Slack token missing")
	}

	return &Alert{
		store:       store,
		slackClient: slack.New(slackToken),
	}
}

// ReportEvent ingests events and sends alerts if needed
func (e *Alert) ReportEvent(ctx context.Context, req *alert.ReportEventRequest, rsp *alert.ReportEventResponse) error {
	if req.Event == nil {
		return errors.New("event can't be empty")
	}
	// ignoring the error intentionally here so we still sends alerts
	// even if persistence is failing
	e.saveEvent(&event{
		ID:       uuid.New().String(),
		Category: req.Event.Category,
		Action:   req.Event.Action,
		Label:    req.Event.Label,
		Value:    req.Event.Value,
	})
	if req.Event.Action == "error" {
		jsond, err := json.MarshalIndent(req.Event, "", "   ")
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Error event received:\n```\n%v\n```", string(jsond))
		_, _, _, err = e.slackClient.SendMessage("errors", slack.MsgOptionUsername("Alert Service"), slack.MsgOptionText(msg, false))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Alert) saveEvent(ev *event) error {
	bytes, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	return e.store.Write(&store.Record{
		Key:   storePrefixEvents,
		Value: bytes})
}
