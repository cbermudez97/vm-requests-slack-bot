package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type VMModalValues struct {
	Dist string
	Type string
}

func buildVMModalValues(i slack.InteractionCallback) VMModalValues {
	dist := i.View.State.Values[distBlockId][distActionId].Value
	tier := i.View.State.Values[vmTypeBlockId][vmTypeActionId].Value

	return VMModalValues{
		Dist: dist,
		Type: tier,
	}
}

func sendUserNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMModalValues) error {
	msgText := fmt.Sprintf(
		"Hi %s. Your request for a VM with:\nDistribution: %s\nType: %s\n\nhave been correctly sended.",
		i.User.Name,
		modalValues.Dist,
		modalValues.Type,
	)

	_, _, err := api.PostMessage(
		i.User.ID,
		slack.MsgOptionText(msgText, false),
		slack.MsgOptionAttachments(),
	)
	return err
}

func interaction(w http.ResponseWriter, r *http.Request) {
	if err := verifySigningSecret(r); err != nil {
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

	modalValues := buildVMModalValues(i)
	log.Info(modalValues)

	botToken := getBotToken()
	api := slack.New(botToken)

	//Notify user that the request is created
	sendUserNotification(api, i, modalValues)

	//TODO: send request to requests channel
}

var InteractionHandler = Handler{
	Endpoint: "/interactions",
	Handler:  interaction,
}
