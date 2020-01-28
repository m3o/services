package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kytra-app/helpers/slack"
)

var token string

func main() {
	flag.StringVar(&token, "token", "", "slack token required:true")
	flag.Parse()

	if token == "" {
		handleErr("Missing Flag", slack.ErrAuth)
	}

	client, err := slack.NewClient(token)
	if err != nil {
		handleErr("Error initializing client", err)
	}

	err = client.SendMessage("test", "my golang client works!!!")
	if err != nil {
		handleErr("Error sending message", err)
	}
}

func handleErr(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%v: %v\n", msg, err)
	os.Exit(2)
}
