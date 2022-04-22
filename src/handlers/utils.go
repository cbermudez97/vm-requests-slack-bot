package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func verifySigningSecret(w http.ResponseWriter, r *http.Request) bool {
	signingSecret := getSigningSecret()
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &sv))

	if err := sv.Ensure(); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
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
