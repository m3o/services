package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/kytra-app/helpers/mailer"
)

func main() {
	handler := mailer.New(os.Getenv("MAILER_USERNAME"), os.Getenv("MAILER_PASSWORD"))

	type Notification struct {
		Title       string
		Description string
		Link        string
		LinkTitle   string
	}

	type EmailData struct {
		Name          string
		URL           string
		Subject       string
		Notifications []Notification
	}

	templateData := EmailData{
		Name:    "JohnDoe",
		URL:     "https://kytra.app",
		Subject: "Kytra Notifications - Monday 17th June",
		Notifications: []Notification{
			Notification{Title: "Title", Description: "Description", Link: "https://kytra.app/demo", LinkTitle: "ClickMe"},
			Notification{Title: "Title2", Description: "Description2", Link: "https://kytra.app/demo", LinkTitle: "ClickMe"},
		},
	}

	t, err := template.ParseFiles("demo.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		fmt.Println(err)
		return
	}
	body := buf.String()

	email := mailer.Email{
		ToName:    "Ben Toogood",
		ToAddress: "bentoogood@gmail.com",
		Subject:   templateData.Subject,
		Body:      body,
	}

	if err := handler.Send(email); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Email Sent!")
}
