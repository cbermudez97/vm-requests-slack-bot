package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func verifySigningSecret(r *http.Request) error {
	signingSecret := getSigningSecret()
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Error(err)
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func getMessageFrom(api *slack.Client, channelID, messageTs string) (slack.Message, error) {
	conversation, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Latest:    messageTs,
		Limit:     1,
		Inclusive: true,
	})
	if err != nil {
		return slack.Message{}, err
	}

	if len(conversation.Messages) < 1 {
		return slack.Message{}, fmt.Errorf("No message found on channel %s with Ts %s", channelID, messageTs)
	}

	msg := conversation.Messages[0]
	return msg, nil
}

func getSigningSecret() string {
	return os.Getenv("SLACK_SIGNING_SECRET")
}

func getBotToken() string {
	return os.Getenv("SLACK_BOT_TOKEN")
}

func getRequestsChannel() string {
	return os.Getenv("SLACK_REQUESTS_CHANNEL")
}

func getApi() *slack.Client {
	return slack.New(getBotToken())
}
