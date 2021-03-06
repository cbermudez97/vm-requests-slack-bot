package interactions

import (
	"fmt"
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers"
	"github.com/cbermudez97/vm-requests-slack-bot/src/vms"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func handleDenyCallback(w http.ResponseWriter, r *http.Request, i slack.InteractionCallback) {
	api := handlers.GetApi()
	msgTs := i.Container.MessageTs
	channelID := i.Container.ChannelID

	// Get request data from message
	requestMsg, err := handlers.GetMessageFrom(api, channelID, msgTs)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestData, err := vms.ParseRequestDataFrom(requestMsg)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: run denied workflow

	// Clear message buttons
	l := len(requestMsg.Blocks.BlockSet)
	blocks := requestMsg.Blocks.BlockSet[:l-1]
	blocks = append(blocks, slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			fmt.Sprintf("*UPDATE*: Denied by <@%s>", i.User.ID),
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

	// Notify user
	requesterUser, err := api.GetUserInfo(requestData.Requester)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _, err = api.PostMessage(
		requestData.Requester,
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject(
					slack.MarkdownType,
					fmt.Sprintf(`Hi %s. Your request for the VM named "%s" have been denied by <@%s>. Please contact him/her for more information.`, requesterUser.Name, requestData.Name, i.User.ID),
					false,
					false,
				),
				nil,
				nil,
			),
		),
	)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
