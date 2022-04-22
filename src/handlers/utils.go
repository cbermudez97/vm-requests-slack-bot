package handlers

import (
	"os"

	"github.com/slack-go/slack"
)

func getSigningSecret() string {
	return os.Getenv("SLACK_SIGNING_SECRET")
}

func getBotToken() string {
	return os.Getenv("SLACK_BOT_TOKEN")
}

func getApi() *slack.Client {
	return slack.New(getBotToken())
}
