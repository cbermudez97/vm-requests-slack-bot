package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const requestsEndpoint = "/request-vm"

const requestCommandStr = "/request-vm"

func requests(w http.ResponseWriter, r *http.Request) {
	if ok := verifySigningSecret(w, r); !ok {
		return
	}

	rCmd, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rCmd.Command == requestCommandStr {
		api := getApi()

		api.PostMessage(rCmd.ChannelID, slack.MsgOptionText("Requesting VM", false))
	} else {
		log.Errorf("Invalid command executed. Expected %s but got %s", requestCommandStr, rCmd.Command)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

var RequestHandler = Handler{
	Endpoint: requestsEndpoint,
	Handler:  requests,
}
