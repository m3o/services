package handler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/kytra-app/helpers/mailer"
	"github.com/kytra-app/helpers/unique"
	notifications "github.com/kytra-app/notifications-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
)

// Handler is responsible for sending emails
type Handler struct {
	template      *template.Template
	mailer        mailer.Service
	users         users.UsersService
	notifications notifications.NotificationsService
}

// New returns an instance of Handler
func New(client client.Client, mailer mailer.Service) (Handler, error) {
	template, err := template.ParseFiles("./bin/template.html")
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		mailer:        mailer,
		template:      template,
		users:         users.NewUsersService("kytra-srv-v1-users:8080", client),
		notifications: notifications.NewNotificationsService("kytra-srv-v1-notifications:8080", client),
	}, nil
}

// SendDailyEmails sends an email to every user who has missed notifications since this time yesterday
func (h Handler) SendDailyEmails() {
	// Fetch the notifications
	startTime := time.Now().AddDate(0, 0, -1)
	nRsp, err := h.notifications.ListNotifications(context.Background(), &notifications.Query{
		StartTime:  startTime.Unix(),
		EndTime:    time.Now().Unix(),
		OnlyUnseen: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the users
	userUUIDs := make([]string, len(nRsp.Notifications))
	for i, n := range nRsp.Notifications {
		userUUIDs[i] = n.UserUuid
	}
	uRsp, err := h.users.List(context.Background(), &users.ListRequest{Uuids: unique.Strings(userUUIDs)})
	if err != nil {
		log.Fatal(err)
	}

	// Group the notifications by UserUUID
	userNotifications := make(map[string][]*notifications.Notification, len(userUUIDs))
	for _, n := range nRsp.Notifications {
		if _, ok := userNotifications[n.UserUuid]; !ok {
			userNotifications[n.UserUuid] = []*notifications.Notification{n}
		} else {
			userNotifications[n.UserUuid] = append(userNotifications[n.UserUuid], n)
		}
	}

	// Send the notifications for each user
	for _, user := range uRsp.Users {
		h.sendEmailToUser(user, userNotifications[user.Uuid], time.Now())
	}
}

func (h Handler) sendEmailToUser(user *users.User, notifications []*notifications.Notification, date time.Time) {
	fmt.Printf("Sending an email to %v, containing %v notification(s)\n", user.Username, len(notifications))

	type Notification struct {
		Title       string
		Description string
		Link        string
		LinkTitle   string
	}

	type EmailData struct {
		Subject          string
		Notifications    []Notification
		NumNotifications int
	}

	ns := make([]Notification, len(notifications))
	for i, n := range notifications {
		ns[i] = Notification{
			Title:       n.Title,
			Description: n.Description,
			LinkTitle:   "View Post",
			Link:        fmt.Sprintf("https://kytra.app/posts/%v", n.ResourceUuid),
		}
	}

	dateSuffix := "th"
	switch date.Day() {
	case 1, 21, 31:
		dateSuffix = "st"
	case 2, 22:
		dateSuffix = "nd"
	case 3, 23:
		dateSuffix = "rd"
	}

	templateData := EmailData{
		Subject:          "Kytra Notifications - " + date.Format("Monday 1"+dateSuffix+" January"),
		Notifications:    ns,
		NumNotifications: len(ns),
	}

	buf := new(bytes.Buffer)
	if err := h.template.Execute(buf, templateData); err != nil {
		fmt.Println(err)
		return
	}

	email := mailer.Email{
		ToName:    strings.Join([]string{user.FirstName, user.LastName}, " "),
		ToAddress: user.Email,
		Subject:   templateData.Subject,
		Body:      buf.String(),
	}

	if err := h.mailer.Send(email); err != nil {
		fmt.Println(err)
		return
	}
}
