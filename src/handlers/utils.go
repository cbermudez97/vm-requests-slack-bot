package handlers

import (
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func handleRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return body, false
	}
	return body, true
}

func handleSigningSecret(w http.ResponseWriter, headers http.Header, body []byte) bool {
	signingSecret := getSigningSecret()

	sv, err := slack.NewSecretsVerifier(headers, signingSecret)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if _, err := sv.Write(body); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
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
