package interactions

import (
	"fmt"
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func buildVMModalValues(i slack.InteractionCallback) VMRequestData {
	dist := i.View.State.Values[handlers.DistBlockId][handlers.DistActionId].SelectedOption.Value
	tier := i.View.State.Values[handlers.VMTypeBlockId][handlers.VMTypeActionId].SelectedOption.Value

	return VMRequestData{
		Requester:    i.User.ID,
		Distribution: dist,
		Type:         tier,
	}
}

func sendChannelNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMRequestData) error {
	requestsChannelId := handlers.GetRequestsChannel()

	msgText := fmt.Sprintf(
		"VM Request\n\nFrom: <@%s>\nDistribution: %s\nType: %s\n",
		i.User.ID,
		modalValues.Distribution,
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

func sendUserNotification(api *slack.Client, i slack.InteractionCallback, modalValues VMRequestData) error {
	msgText := fmt.Sprintf(
		"Hi %s. Your request for a VM with:\n\nDistribution: %s\nType: %s\n\nhave been correctly created.",
		i.User.Name,
		modalValues.Distribution,
		modalValues.Type,
	)

	_, _, err := api.PostMessage(
		i.User.ID,
		slack.MsgOptionText(msgText, false),
	)
	return err
}

func handleRequestModal(w http.ResponseWriter, r *http.Request, i slack.InteractionCallback) {
	modalValues := buildVMModalValues(i)

	api := handlers.GetApi()

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
