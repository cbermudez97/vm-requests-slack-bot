package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		"Hi %s. Your request for a VM with:\n\nDistribution: %s\nType: %s\n\nhave been correctly created.",
		i.User.RealName,
		modalValues.Dist,
		modalValues.Type,
	)

	_, _, err := api.PostMessage(
		i.User.ID,
		slack.MsgOptionText(msgText, false),
	)
	return err
}

func sendChannelNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMModalValues) error {
	requestsChannelId := getRequestsChannel()

	msgText := fmt.Sprintf(
		"VM Request from <@%s>.\nData:\nDistribution: %s\nType: %s\n",
		i.User.ID,
		modalValues.Dist,
		modalValues.Type,
	)

	// Build message block
	messageBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			msgText,
			true,
			false,
		),
		nil,
		nil,
	)

	// Build actions block
	acceptOrDenyRequestBlock := slack.NewActionBlock(
		acceptOrDenyBlockID,
		slack.NewButtonBlockElement(
			denyActionID,
			denyActionValue,
			slack.NewTextBlockObject(
				slack.PlainTextType,
				denyActionText,
				false,
				false,
			),
		),
		slack.NewButtonBlockElement(
			acceptActionID,
			acceptActionValue,
			slack.NewTextBlockObject(
				slack.PlainTextType,
				acceptActionText,
				false,
				false,
			),
		),
	)

	_, _, err := api.PostMessage(
		requestsChannelId,
		slack.MsgOptionBlocks(messageBlock, acceptOrDenyRequestBlock),
	)
	return err
}

func interactions(w http.ResponseWriter, r *http.Request) {
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
		switch i.View.CallbackID {
		case requestModalCallbackId: // Handle request modal
			handleRequestModal(w, r, i)
		}
		return

	} else if i.Type == slack.InteractionTypeMessageAction {
		// TODO: message interaction
		log.Info(i)
	}
}

func handleRequestModal(w http.ResponseWriter, r *http.Request, i slack.InteractionCallback) {
	modalValues := buildVMModalValues(i)
	log.Infof("VM Modal config: %s", modalValues)

	api := getApi()

	//Send request to requests channel
	if err := sendChannelNotification(api, i, modalValues); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Notify user that the request is created
	if err := sendUserNotification(api, i, modalValues); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var InteractionHandler = Handler{
	Endpoint: "/interactions",
	Handler:  interactions,
}
