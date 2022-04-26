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
	dist := i.View.State.Values[distBlockId][distActionId].SelectedOption.Value
	tier := i.View.State.Values[vmTypeBlockId][vmTypeActionId].SelectedOption.Value

	return VMModalValues{
		Dist: dist,
		Type: tier,
	}
}

func sendUserNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMModalValues) error {
	msgText := fmt.Sprintf(
		"Hi %s. Your request for a VM with:\nDistribution: %s\nType: %s\n\nhave been correctly created.",
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

func sendChannelNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMModalValues) error {
	requestsChannelId := getRequestsChannel()

	msgText := fmt.Sprintf(
		"VM Request from %s.\nData:\nDistribution: %s\nType: %s\n",
		i.User.Name,
		modalValues.Dist,
		modalValues.Type,
	)

	_, _, err := api.PostMessage(
		requestsChannelId,
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

	if i.Type == slack.InteractionTypeViewSubmission {
		modalValues := buildVMModalValues(i)
		log.Info(modalValues)

		api := getApi()

		//Notify user that the request is created
		if err := sendUserNotification(api, i, modalValues); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Send request to requests channel
		if err := sendChannelNotification(api, i, modalValues); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

var InteractionHandler = Handler{
	Endpoint: "/interactions",
	Handler:  interaction,
}
