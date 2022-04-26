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

type VMRequestData struct {
	Requester string
	Dist      string
	Type      string
}

func buildVMModalValues(i slack.InteractionCallback) VMRequestData {
	dist := i.View.State.Values[distBlockId][distActionId].SelectedOption.Value
	tier := i.View.State.Values[vmTypeBlockId][vmTypeActionId].SelectedOption.Value

	return VMRequestData{
		Requester: i.User.ID,
		Dist:      dist,
		Type:      tier,
	}
}

func sendUserNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMRequestData) error {
	msgText := fmt.Sprintf(
		"Hi %s. Your request for a VM with:\n\nDistribution: %s\nType: %s\n\nhave been correctly created.",
		i.User.Name,
		modalValues.Dist,
		modalValues.Type,
	)

	_, _, err := api.PostMessage(
		i.User.ID,
		slack.MsgOptionText(msgText, false),
	)
	return err
}

func sendChannelNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMRequestData) error {
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
			false,
			false,
		),
		nil,
		nil,
	)

	// Build actions block
	acceptOrDenyRequestBlock := slack.NewActionBlock(
		acceptOrDenyBlockID,
		slack.ButtonBlockElement{
			Type:     slack.METButton,
			ActionID: acceptActionID,
			Value:    acceptActionValue,
			Text: slack.NewTextBlockObject(
				slack.PlainTextType,
				acceptActionText,
				false,
				false,
			),
			Style: slack.StylePrimary,
		},
		slack.ButtonBlockElement{
			Type:     slack.METButton,
			ActionID: denyActionID,
			Value:    denyActionValue,
			Text: slack.NewTextBlockObject(
				slack.PlainTextType,
				denyActionText,
				false,
				false,
			),
			Style: slack.StyleDanger,
		},
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

func parseRequestDataFrom(message slack.Message) (VMRequestData, error) {
	data := VMRequestData{}

	// FIXME: fix parsing
	// if len(message.Blocks.BlockSet) < 1 {
	// 	return data, fmt.Errorf("Invalid message structure")
	// }

	// msgTextSection := message.Blocks.BlockSet[0]
	// switch msgTextSection.(type) {
	// case slack.SectionBlock:
	// 	// TODO: parse data from text
	// default:
	// 	return data, fmt.Errorf("Invalid message structure")
	// }

	return data, nil
}

func handleAcceptCallback(w http.ResponseWriter, r *http.Request, i slack.InteractionCallback) {
	api := getApi()
	msgTs := i.Container.MessageTs
	channelID := i.Container.ChannelID

	// Get request data from message
	requestMsg, err := getMessageFrom(api, channelID, msgTs)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = parseRequestDataFrom(requestMsg)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: run creation workflow

	// Clear message buttons
	l := len(requestMsg.Blocks.BlockSet)
	blocks := requestMsg.Blocks.BlockSet[:l-1]
	blocks = append(blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			fmt.Sprintf("Accepted by <@%s>", i.User.ID),
			false,
			false,
		),
		nil,
		nil,
	))

	_, _, _, err = api.UpdateMessage(
		channelID,
		msgTs,
		slack.MsgOptionBlocks(blocks...),
	)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: notify user

	w.WriteHeader(http.StatusOK)
}

func handleDenyCallback(w http.ResponseWriter, r *http.Request, i slack.InteractionCallback) {
	log.Infof("VM Request Block: Denied")

	// TODO: run denial workflow

	// TODO: clear message buttons

	// TODO: notify user

	w.WriteHeader(http.StatusOK)
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

	w.WriteHeader(http.StatusOK)
}

var InteractionHandler = Handler{
	Endpoint: "/interactions",
	Handler:  interactions,
}
