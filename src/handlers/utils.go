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

func VerifySigningSecret(r *http.Request) error {
	signingSecret := GetSigningSecret()
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

func GetMessageFrom(api *slack.Client, channelID, messageTs string) (slack.Message, error) {
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

func GetSigningSecret() string {
	return os.Getenv("SLACK_SIGNING_SECRET")
}

func GetBotToken() string {
	return os.Getenv("SLACK_BOT_TOKEN")
}

func GetRequestsChannel() string {
	return os.Getenv("SLACK_REQUESTS_CHANNEL")
}

func GetApi() *slack.Client {
	return slack.New(GetBotToken())
}
