package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func interaction(w http.ResponseWriter, r *http.Request) {
	if err := verifySigningSecret(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var i slack.InteractionCallback
	if err := json.Unmarshal([]byte(r.FormValue("payload")), &i); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//TODO: handle modal
	log.Debug(i.View.State)
}

var InteractionHandler = Handler{
	Endpoint: "/interactions",
	Handler:  interaction,
}
