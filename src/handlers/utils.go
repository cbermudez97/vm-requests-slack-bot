package handlers

import (
	"bytes"
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

func getSigningSecret() string {
	return os.Getenv("SLACK_SIGNING_SECRET")
}

func getBotToken() string {
	return os.Getenv("SLACK_BOT_TOKEN")
}

func getApi() *slack.Client {
	return slack.New(getBotToken())
}
