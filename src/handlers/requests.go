package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// Handler configs
const requestsEndpoint = "/request-vm"
const requestCommandStr = "/request-vm"

// OS block data
const osBlockId = "os"

var osOptions = []string{
	"Ubuntu",
	"Windows",
	"MacOS",
}

// Tier block data
const tierBlockId = "tier"

var tierOptions = []string{
	"Light",
	"Normal",
	"Heavy",
}

func createOptionBlockObjects(options []string) []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(options))
	for _, option := range options {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, option, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(option, optionText, nil))
	}
	return optionBlockObjects
}

func buildVMRequestModal() slack.ModalViewRequest {
	// Modal texts
	titleText := slack.NewTextBlockObject(slack.PlainTextType, "VM Request", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Submit", false, false)

	// Header section
	headerText := slack.NewTextBlockObject(slack.MarkdownType, "*Please fill out the request form below:*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// OS input
	osOptionsElems := createOptionBlockObjects(osOptions)
	osText := slack.NewTextBlockObject(slack.PlainTextType, "Select an OS", false, false)
	osOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, nil, "OS", osOptionsElems...)
	osBlock := slack.NewInputBlock(osBlockId, osText, osOption)

	// VM tier input
	tierOptionsElems := createOptionBlockObjects(tierOptions)
	tierText := slack.NewTextBlockObject(slack.PlainTextType, "Select a Tier", false, false)
	tierOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, nil, "Tier", tierOptionsElems...)
	tierBlock := slack.NewInputBlock(tierBlockId, tierText, tierOption)

	// Additional details
	// TODO: define additional details modal

	// Blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			osBlock,
			tierBlock,
		},
	}

	// Modal
	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func requestsCreation(w http.ResponseWriter, r *http.Request) {
	// Verify signing secret
	if err := verifySigningSecret(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Build command
	rCmd, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Handle command
	if rCmd.Command == requestCommandStr {
		api := getApi()

		modalRequest := buildVMRequestModal()
		_, err := api.OpenView(rCmd.TriggerID, modalRequest)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Errorf("Invalid command executed. Expected %s but got %s", requestCommandStr, rCmd.Command)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

var RequestHandler = Handler{
	Endpoint: requestsEndpoint,
	Handler:  requestsCreation,
}
