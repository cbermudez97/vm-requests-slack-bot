package interactions

import (
	"encoding/json"
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// Attachment buttons data
const acceptOrDenyBlockID = "accept_or_deny"

const acceptActionID = "accept"
const acceptActionText = "Accept"
const acceptActionValue = "accept"

const denyActionID = "deny"
const denyActionText = "Deny"
const denyActionValue = "deny"

type VMRequestData struct {
	Requester    string
	Distribution string
	Type         string
}

func interactions(w http.ResponseWriter, r *http.Request) {
	if err := handlers.VerifySigningSecret(r); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var i slack.InteractionCallback
	if err := json.Unmarshal([]byte(r.FormValue("payload")), &i); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if i.Type == slack.InteractionTypeViewSubmission {
		switch i.View.CallbackID {
		case handlers.RequestModalCallbackId: // Handle request modal
			handleRequestModal(w, r, i)
		}
	} else if i.Type == slack.InteractionTypeBlockActions {
		for _, action := range i.ActionCallback.BlockActions {
			switch action.ActionID { // Allow to handle more block actions
			case acceptActionID:
				handleAcceptCallback(w, r, i)
				return
			case denyActionID:
				handleDenyCallback(w, r, i)
				return
			}
		}
	}
}

var InteractionHandler = handlers.Handler{
	Endpoint: "/interactions",
	Handler:  interactions,
}
