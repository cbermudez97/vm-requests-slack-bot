package handlers

import (
	"net/http"

	"github.com/slack-go/slack"
)

const requestsEndpoint = "/request-vm"

const requestCommandStr = "/request-vm"

func requests(w http.ResponseWriter, r *http.Request) {
	body, ok := handleRequestBody(w, r)
	if !ok {
		return
	}

	if ok := handleSigningSecret(w, r.Header, body); !ok {
		return
	}

	rCmd, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rCmd.Command == requestCommandStr {
		api := getApi()

		api.PostMessage(rCmd.ChannelID, slack.MsgOptionText("Requesting VM", false))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

var RequestHandler = Handler{
	Endpoint: requestsEndpoint,
	Handler:  requests,
}
