package interactions

import (
	"fmt"
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers"
	"github.com/cbermudez97/vm-requests-slack-bot/src/vms"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func sendChannelNotification(api *slack.Client, i slack.InteractionCallback, modalValues vms.VMRequest) error {
	requestsChannelId := handlers.GetRequestsChannel()

	msgText := vms.BuildRequestNotificationMessage(modalValues)

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

func sendUserNotification(api *slack.Client, i slack.InteractionCallback, modalValues vms.VMRequest) error {
	msgText := fmt.Sprintf(
		`Hi %s. Your request for the VM named "%s" have been correctly created. I will notify you upon it's acceptance or denial.`,
		i.User.Name,
		modalValues.Name,
	)

	_, _, err := api.PostMessage(
		i.User.ID,
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject(
					slack.MarkdownType,
					msgText,
					false,
					false,
				),
				nil,
				nil,
			),
		),
	)
	return err
}

func handleRequestModal(w http.ResponseWriter, r *http.Request, i slack.InteractionCallback) {
	modalValues := vms.NewVMRequestFromModal(i)

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

func handleModalProviderChangedCallback(w http.ResponseWriter, r *http.Request, i slack.InteractionCallback) {
	modalValues := vms.NewVMRequestFromModal(i)

	api := handlers.GetApi()
	// viewValues := i.View.State.Values
	// if _, ok := viewValues[vms.VMTypeBlockId]; ok { // Modal have a previous selected provider
	// 	// TODO: handle this case
	// } else { // No previous selected provider
	// Add types for new provider
	provider, err := vms.FindProviderByValue(modalValues.Provider)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	modal := vms.BuildVMRequestModalWithTypes(*provider)

	_, err = api.UpdateView(
		modal,
		i.View.ExternalID,
		i.View.Hash,
		i.View.ID,
	)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// }
}
